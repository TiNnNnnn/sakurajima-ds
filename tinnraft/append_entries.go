package tinnraft

import (
	"context"
	"sakurajima-ds/tinnraftpb"

	_ "google.golang.org/grpc/peer"
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
	for _, peer := range rf.peers {
		if int(peer.id) == rf.me {
			//防止Leader节点无意义的选举发送
			rf.resetElectionTimer()
			continue
		}

		prevLogIndex := uint64(rf.nextIndex[peer.id] - 1)
		/*
		   在复制快照的时候先判断到peer的prevLogIndex,如果比当前日志的第一条
		   索引号还小，就说明Leader已经把这条日志打到快照中了，此时构造
		   InstallSnapshotArgs调用Snapshot RPC将快照数据发送给Followr节点
		*/
		if prevLogIndex < uint64(rf.log.GetPersistFirstEntry().Index) {
			firstLog := rf.log.GetPersistFirstEntry()
			snapShotArgs := &tinnraftpb.InstallSnapshotArgs{
				Term:              int64(rf.currentTerm),
				LeaderId:          int64(rf.me),
				LastIncludedIndex: firstLog.Index,
				LastIncludeTerm:   int64(firstLog.Term),
				Data:              rf.ReadSnapshot(),
			}
			go rf.leaderSendSnapshots(int(peer.id), snapShotArgs)

		} else {
			lastLog := rf.log.GetPersistLastEntry()
			if int(lastLog.Index) >= rf.nextIndex[peer.id] || isHeartbeat {
				// nextsend_index := rf.nextIndex[peer.id]
				// if nextsend_index <= 0 {
				// 	nextsend_index = 1
				// }
				// if int(lastLog.Index)+1 < nextsend_index {
				// 	nextsend_index = int(lastLog.Index)
				// }
				firstLogIndex := rf.log.GetPersistFirstEntry().Index
				entires := make([]*tinnraftpb.Entry, len(rf.log.TruncatePersistLog(int64(prevLogIndex)+1-firstLogIndex)))
				copy(entires, rf.log.TruncatePersistLog(int64(prevLogIndex)+1-firstLogIndex))
				args := tinnraftpb.AppendEntriesArgs{
					Term:         int64(rf.currentTerm),
					LeaderId:     int64(rf.me),
					PrevLogIndex: int64(prevLogIndex),
					PrevLogTerm:  int64(rf.log.GetEntryWithoutLock(int64(prevLogIndex) - firstLogIndex).Term),
					Entries:      entires,
					LeaderCommit: int64(rf.commitIndex),
				}
				//copy(args.Entries, rf.log.slice2(nextsend_index))
				go rf.leaderSendEntries(int(peer.id), &args)

			}
		}
	}
}

func (rf *Raft) leaderSendSnapshots(serverId int, args *tinnraftpb.InstallSnapshotArgs) {
	DLog("[%v]: term %v | send snapshot to %v with args: %s", rf.me, rf.currentTerm, serverId, args.String())
	reply, err := rf.sendAppendSnapshots(serverId, args)
	if err != nil {
		DLog("[%v]: term %v | send snapshot to %v failed! %v", rf.me, rf.currentTerm, serverId, err.Error())
		return
	}

	DLog("[%v]: term %v | send snapshot to %v with reply: %s", rf.me, rf.currentTerm, serverId, reply.String())

	rf.mu.Lock()
	defer rf.mu.Unlock()

	if reply != nil {
		if rf.state == Leader && rf.currentTerm == int(reply.Term) {
			if reply.Term > int64(rf.currentTerm) {
				rf.setNewTerm(int(reply.Term))
			} else {
				DLog("set [%v] matchIDx: %d , nextIdx: %d", serverId, args.LastIncludedIndex, args.LastIncludedIndex+1)
				rf.matchIndex[serverId] = int(args.LastIncludedIndex)
				rf.nextIndex[serverId] = int(args.LastIncludedIndex) + 1
			}
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
	if rf.state == Leader && int(args.Term) == rf.currentTerm {
		if reply.Success {
			//追加日志成功
			match := int(args.PrevLogIndex) + len(args.Entries)
			next := match + 1
			rf.nextIndex[serverId] = max(rf.nextIndex[serverId], next)
			rf.matchIndex[serverId] = max(rf.matchIndex[serverId], match)
			if len(args.Entries) > 0 {
				DLog("[%v]: append entries to %v success, next %v match %v\n", rf.me, serverId, rf.nextIndex[serverId], rf.matchIndex[serverId])
			} else {
				DLog("[%v]: term: %v | send heartbeats to %v success", rf.me, rf.currentTerm, serverId)
			}

		} else if reply.Conflict {
			//追加失败，产生冲突
			DLog("[%v]: confict from %v: %#v\n", rf.me, serverId, reply)
			if reply.XTerm == -1 {
				//preLogIndex发生冲突
				rf.nextIndex[serverId] = int(reply.XLen)
			} else {
				//preLogTerm发生冲突
				lastLogInXTerm := rf.findLastLogInTerm(int(reply.XTerm))
				DLog("[%v]: lastLogInXTerm %v\n", rf.me, lastLogInXTerm)
				if lastLogInXTerm > 0 {
					rf.nextIndex[serverId] = lastLogInXTerm
				} else {
					rf.nextIndex[serverId] = int(reply.XIndex)
				}
			}
			DLog("[%v]: leader nextIndex[%v] %v\n", rf.me, serverId, rf.nextIndex[serverId])
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
		if int(rf.log.GetPersistEntryByidx(int64(k)).Term) != rf.currentTerm {
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
		term := int(rf.log.GetPersistEntryByidx(int64(i)).Term)
		if term == x {
			return i
		} else if term < x {
			break
		}
	}
	return -1
}

// 处理leader发来的AppendEntries请求
func (rf *Raft) HandleAppendEntries(args *tinnraftpb.AppendEntriesArgs, reply *tinnraftpb.AppendEntriesReply) {
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

	rf.leaderId = int(args.LeaderId)
	rf.resetElectionTimer()
	//Candidater在选举中收到了来自其他Leader的心跳，且任期更大
	//说明选举失败，变回Follower
	if rf.state == Candidate {
		rf.state = Follower
	}

	if args.PrevLogIndex < rf.log.GetPersistFirstEntry().Index {
		return
	}

	//PrevLogIndex发生冲突
	if rf.log.GetPersistLastEntry().Index < args.PrevLogIndex {
		reply.Conflict = true
		reply.XTerm = -1
		reply.XIndex = -1
		reply.XLen = int64(rf.log.LogPersistLen())
		DLog("[%v]: Conflict XTerm %v,XIndex %v,XLen %v", rf.me, reply.XTerm, reply.XIndex, reply.XLen)
		return
	}

	//PrevLogTerm发生冲突
	if int64(rf.log.GetPersistEntryByidx(args.PrevLogIndex).Term) != args.PrevLogTerm {
		reply.Conflict = true
		xTerm := rf.log.GetPersistEntryByidx(args.PrevLogIndex).Term
		firstIndex := rf.log.GetPersistFirstEntry().Index
		//找到上个任期的最后一条日志条目的索引，并记录为Xindex
		for xIndex := args.PrevLogIndex; xIndex >= firstIndex; xIndex-- {
			if rf.log.GetPersistEntryByidx(xIndex-1).Term != xTerm {
				reply.XIndex = xIndex
				break
			}
		}
		reply.XTerm = int64(xTerm)
		reply.XLen = int64(rf.log.LogPersistLen())
		DLog("[%v]: Conflict XTerm %v,XIndex %v,XLen %v", rf.me, reply.XTerm, reply.XIndex, reply.XLen)
		return
	}

	for idx, entry := range args.Entries {
		/*
			如果follower中已经存在的日志条目和追加日志条目发生
			冲突(索引相同,但是任期不同),那么就删除这个已经存在
			的条目以及之后的所有条目
		*/
		if entry.Index <= rf.log.GetPersistLastEntry().Index &&
			rf.log.GetPersistEntryByidx(entry.Index).Term != entry.Term {
			rf.log.TruncatePersistLog(int64(entry.Index))
			rf.persist()
		}
		//追加日志中尚未保存的任何新日志条目
		if entry.Index > rf.log.GetPersistLastEntry().Index {
			for _, newEntry := range args.Entries[idx:] {
				rf.log.PersistAppend(newEntry)
			}
			DLog("[%v]: append entries from leader: [%v]", rf.me, args.Entries[idx:])
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

func (rf *Raft) sendAppendSnapshots(server int, args *tinnraftpb.InstallSnapshotArgs) (*tinnraftpb.InstallSnapshotReply, error) {
	return (*rf.peers[server].raftServiceCli).Snapshot(context.Background(), args)
}
