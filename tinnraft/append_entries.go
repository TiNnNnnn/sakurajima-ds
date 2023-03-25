package tinnraft

import (
	"context"
	"fmt"
	"sakurajima-ds/tinnraftpb"
)

// // 日志添加请求
// type AppendEntriesArgs struct {
// 	Term         int                //当前任期
// 	LeaderId     int                //领导者id
// 	PrevLogIndex int                //紧邻待追加日志条目之前的那个日志条目索引
// 	PrevLogTerm  int                //紧邻待追加日志条目之前的那个日志条目任期
// 	Entries      []tinnraftpb.Entry //需要追加到follwer的日志
// 	LeaderCommit int                //领导者已知的已经提交的最高的日志条目
// }

// // 日志添加响应
// type AppendEntriesReply struct {
// 	Term     int  //当前任期，对于leader,会更新自己的term
// 	Success  bool //true:follower所含的条目与prevLogIndex以及prevLogTerm匹配成功
// 	Conflict bool //日志追加是否发生冲突

// 	//当日志追加冲突时的记录字段
// 	XTerm  int
// 	XIndex int
// 	XLen   int
// }

func (rf *Raft) appendEntries(isHeartbeat bool) {
	lastLog := rf.log.GetPersistLastEntry()
	for peer := range rf.peers {
		if peer == rf.me {
			rf.resetElectionTimer()
			continue
		}
		if int(lastLog.Index) >= rf.nextIndex[peer] || isHeartbeat {
			nextsend_index := rf.nextIndex[peer]
			if nextsend_index <= 0 {
				nextsend_index = 1
			}
			if int(lastLog.Index)+1 < nextsend_index {
				nextsend_index = int(lastLog.Index)
			}
			nextsend_index_prev := rf.log.at(nextsend_index - 1)

			args := tinnraftpb.AppendEntriesArgs{
				Term:         int64(rf.currentTerm),
				LeaderId:     int64(rf.me),
				PrevLogIndex: int64(nextsend_index_prev.Index),
				PrevLogTerm:  int64(nextsend_index_prev.Term),
				Entries:      make([]*tinnraftpb.Entry, int(lastLog.Index)-nextsend_index+1),
				LeaderCommit: int64(rf.commitIndex),
			}
			copy(args.Entries, rf.log.slice2(nextsend_index))
			go rf.leaderSendEntries(peer, &args)

		}
	}
}

func (rf *Raft) leaderSendEntries(serverId int, args *tinnraftpb.AppendEntriesArgs) {
	reply, err := rf.sendAppendEntries(serverId, args)
	if err != nil {
		return
	}

	rf.mu.Lock()
	defer rf.mu.Unlock()

	if int(reply.Term) > rf.currentTerm {
		rf.setNewTerm(int(reply.Term))
		return
	}
	if int(args.Term) == rf.currentTerm {
		if reply.Success { //追加日志成功
			match := int(args.PrevLogIndex) + len(args.Entries)
			next := match + 1
			rf.nextIndex[serverId] = max(rf.nextIndex[serverId], next)
			rf.matchIndex[serverId] = max(rf.matchIndex[serverId], match)
			fmt.Printf("[%v]: %v 追加成功 next %v match %v\n", rf.me, serverId, rf.nextIndex[serverId], rf.matchIndex[serverId])
		} else if reply.Conflict { //追加失败，产生冲突
			fmt.Printf("[%v]: 冲突来自 %v %#v\n", rf.me, serverId, reply)
			if reply.XTerm == -1 { //preLogIndex发生冲突
				rf.nextIndex[serverId] = int(reply.XLen)
			} else { //preLogTerm发生冲突
				lastLogInXTerm := rf.findLastLogInTerm(int(reply.XTerm))
				fmt.Printf("[%v]: lastLogInXTerm %v\n", rf.me, lastLogInXTerm)
				if lastLogInXTerm > 0 {
					rf.nextIndex[serverId] = lastLogInXTerm
				} else {
					rf.nextIndex[serverId] = int(reply.XIndex)
				}
			}
			fmt.Printf("[%v]: leader nextIndex[%v] %v\n", rf.me, serverId, rf.nextIndex[serverId])
		} else if rf.nextIndex[serverId] > 1 {
			//减少nextIndex进行重发
			rf.nextIndex[serverId]--
		}
		rf.leaderCommitRule()
	}

}

// 如果存在一个N满足N>commitIndex,并且大多数的matchIndex[i]>N成立
// 且log[N].term == currentTerm,那么commitIndex等于这个N
func (rf *Raft) leaderCommitRule() {
	if rf.state != Leader {
		return
	}

	for k := rf.commitIndex + 1; k <= int(rf.log.GetPersistLastEntry().Index); k++ {
		if int(rf.log.at(k).Term) != rf.currentTerm {
			continue
		}
		counter := 1
		for serverId := 0; serverId < len(rf.peers); serverId++ {
			if serverId != rf.me && rf.matchIndex[serverId] >= k {
				counter++
			}
			if counter > len(rf.peers)/2 {
				rf.commitIndex = k
				rf.apply()
				break
			}
		}
	}
}

// 从后往前查找任期为x的日志条目的索引
func (rf *Raft) findLastLogInTerm(x int) int {
	for i := int(rf.log.GetPersistLastEntry().Index); i > 0; i-- {
		term := int(rf.log.at(i).Term)
		if term == x {
			return i
		} else if term < x {
			break
		}
	}
	return -1
}

func (rf *Raft) AppendEntries(args *tinnraftpb.AppendEntriesArgs, reply *tinnraftpb.AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	reply.Success = false
	reply.Term = int64(rf.currentTerm)
	if args.Term > int64(rf.currentTerm) {
		rf.setNewTerm(int(args.Term))
		return
	}
	//Leader的任期小于follower的任期，直接结束
	if args.Term < int64(rf.currentTerm) {
		return
	}
	rf.resetElectionTimer()

	//Candidater在选举中收到了来自其他Leader的心跳，且任期更大
	//说明选举失败，变回Follower
	if rf.state == Candidate {
		rf.state = Follower
	}

	//PrevLogIndex发生冲突
	if rf.log.GetPersistLastEntry().Index < args.PrevLogIndex {
		reply.Conflict = true
		reply.XTerm = -1
		reply.XIndex = -1
		reply.XLen = int64(rf.log.LenLog())
		return
	}
	//PrevLogTerm发生冲突
	if int64(rf.log.at(int(args.PrevLogIndex)).Term) != args.PrevLogTerm {
		reply.Conflict = true
		xTerm := rf.log.at(int(args.PrevLogIndex)).Term
		//找到上个任期的最后一条日志条目的索引，并记录为Xindex
		for xIndex := args.PrevLogIndex; xIndex > 0; xIndex-- {
			if rf.log.at(int(xIndex-1)).Term != xTerm {
				reply.XIndex = xIndex
				break
			}
		}
		reply.XTerm = int64(xTerm)
		reply.XLen = int64(rf.log.LenLog())
		return
	}

	for idx, entry := range args.Entries {
		//如果follower中已经存在的日志条目和追加日志条目发生
		//冲突(索引相同,但是任期不同),那么就删除这个已经存在
		//的条目以及之后的所有条目
		if entry.Index <= rf.log.GetPersistLastEntry().Index &&
			rf.log.at(int(entry.Index)).Term != entry.Term {
			rf.log.TruncatePersistLog(int64(entry.Index))
			rf.persist()
		}
		//追加日志中尚未保存的任何新日志条目
		if entry.Index > rf.log.GetPersistLastEntry().Index {

			rf.log.append2(args.Entries[idx:])
			rf.persist()
			break
		}
	}
	//如果Leader提交的最高的日志的索引LeaderCommit大于fllower的
	//已经提交的最高的日志条目的索引,commitIndex重置为
	//min(LeaderCommit,最新追加日志的最高索引)
	if args.LeaderCommit > int64(rf.commitIndex) {
		rf.commitIndex = min(int(args.LeaderCommit), int(rf.log.GetPersistLastEntry().Index))
		rf.apply()
	}
	reply.Success = true

}

// rpc 调用对端的AppendEntries方法
func (rf *Raft) sendAppendEntries(server int, args *tinnraftpb.AppendEntriesArgs) (*tinnraftpb.AppendEntriesReply, error) {
	//ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	return (*rf.peers[server].raftServiceCli).AppendEntries(context.Background(), args)
}
