package tinnraft

import (
	"context"
	"sakurajima-ds/tinnraftpb"
	"sync"
)

// // 请求投票request
// type RequestVoteArgs struct {
// 	// Your data here (2A, 2B).
// 	Term         int //候选人任期
// 	CandidateId  int //候选人id
// 	LastLogIndex int //候选人最后日志条目的索引值
// 	LastLogTerm  int //候选人最后日志条目的任期号
// }

// // 请求投票response
// type RequestVoteReply struct {
// 	// Your data here (2A).
// 	Term        int  //当前任期,方便候选人更新自己的任期
// 	VoteGranted bool //候选人是否赢得了该张选票
// }

// 候选人发起投票请求
func (rf *Raft) candidateRequestVote(serverId int, args *tinnraftpb.RequestVoteArgs, voteCounter *int, becomeLeader *sync.Once) {
	reply := tinnraftpb.RequestVoteReply{}
	//发起rpc投票并接受结果
	ok := rf.sendRequestVote(serverId, args, &reply)
	if !ok {
		return
	}
	rf.mu.Lock()
	defer rf.mu.Unlock()

	//发现fllower的term比自己大
	if reply.Term > args.Term {
		rf.setNewTerm(int(reply.Term))
		return
	}

	//如果fllower的term比自己小，投票作废
	if reply.Term < args.Term {
		return
	}
	//follwer未投票,直接返回
	if !reply.VoteGranted {
		return
	}

	//选票增加
	*voteCounter++

	//已经获得超过半数的选票
	if *voteCounter > len(rf.peers)/2 && rf.currentTerm == int(args.Term) && rf.state == Candidate {
		becomeLeader.Do(func() {
			rf.state = Leader
			LastLogIndex := rf.log.lastLog().Index
			for i := range rf.peers {
				rf.nextIndex[i] = int(LastLogIndex) + 1
				rf.matchIndex[i] = 0
			}
			//发送心跳给其他server
			rf.appendEntries(true)
		})
	}
}

func (rf *Raft) RequestVote(args *tinnraftpb.RequestVoteArgs, reply *tinnraftpb.RequestVoteReply) {
	// Your code here (2A, 2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	//发现候选人的term比自己大
	if int(args.Term) > rf.currentTerm {
		rf.setNewTerm(int(args.Term))
	}
	//发现候选人的term比自己小
	if int(args.Term) < rf.currentTerm {
		reply.Term = int64(rf.currentTerm)
		reply.VoteGranted = false
		return
	}

	follow_lastLog := rf.log.lastLog()
	upToDate := uint64(args.LastLogTerm) > follow_lastLog.Term ||
		(uint64(args.LastLogTerm) == follow_lastLog.Term &&
			args.LastLogIndex >= follow_lastLog.Index)

	if (rf.votedFor == -1 || rf.votedFor == int(args.CandidateId)) && upToDate {
		reply.VoteGranted = true
		rf.votedFor = int(args.CandidateId)
		//持久化
		rf.persist()
		rf.resetElectionTimer()
	} else {
		reply.VoteGranted = false
	}
	reply.Term = int64(rf.currentTerm)
}

func (rf *Raft) sendRequestVote(server int, args *tinnraftpb.RequestVoteArgs, reply *tinnraftpb.RequestVoteReply) bool {
	//ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	reply, err := (*rf.peers[server].raftServiceCli).RequestVote(context.Background(), args)
	if err != nil {
		return false
	}
	return true
}
