package tinnraft

import (
	"context"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"sync"
	"syscall"
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
	//发起rpc投票并接受结果

	//LOG
	raftlog := &tinnraftpb.LogArgs{
		Op:       tinnraftpb.LogOp_RequestVote,
		Contents: "send request vote",
		FromId:   strconv.Itoa(rf.me),
		ToId:     strconv.Itoa(serverId),
		CurState: "candidate",
		Pid:      int64(syscall.Getpid()),
		Term:     int64(rf.currentTerm),
		Layer:    tinnraftpb.LogLayer_RAFT,
	}
	rf.apiGateClient.SendLogToGate(raftlog)

	DLog("[%v]: term %v | send request vote to %v ", rf.me, rf.currentTerm, serverId)
	reply, err := (*rf.peers[serverId].raftServiceCli).RequestVote(context.Background(), args)
	if err != nil {
		DLog("[%v]: term %v | send request vote to %d failed! %v", rf.me, rf.currentTerm, serverId, err.Error())
		return
	}

	if reply == nil {
		return
	}

	rf.mu.Lock()
	defer rf.mu.Unlock()

	//发现fllower的term比自己大
	if reply.Term > args.Term {
		rf.ChangeRaftState(Follower)
		rf.currentTerm = int(reply.Term)
		rf.votedFor = -1
		rf.persist()
		return
	}

	//发现foller的term比自己小，说明foller任期失效
	if reply.Term < args.Term {
		return
	}

	//follwer未投票给自己,直接返回
	if !reply.VoteGranted {
		return
	}

	//选票增加
	DLog("[%v]: recieve one vote from [%v]", rf.me, serverId)
	*voteCounter++

	//已经获得超过半数的选票
	if *voteCounter > len(rf.peers)/2 && rf.currentTerm == int(args.Term) && rf.state == Candidate {
		//LOG
		raftlog := &tinnraftpb.LogArgs{
			Op:       tinnraftpb.LogOp_ToLeader,
			Contents: "win vote and become leader",
			FromId:   strconv.Itoa(rf.me),
			PreState: "follower",
			CurState: "leader",
			Pid:      int64(syscall.Getpid()),
			Term:     int64(rf.currentTerm),
			Layer:    tinnraftpb.LogLayer_RAFT,
		}
		rf.apiGateClient.SendLogToGate(raftlog)

		DLog("[%d]: term %v | recieve the most votes, win the election,become the leader\n", rf.me, rf.currentTerm)
		becomeLeader.Do(func() {
			rf.ChangeRaftState(Leader)
			//发送心跳给其他server
			rf.appendEntries(true)
			*voteCounter = 0
		})
	}
}

// 处理其他server发来的requestvote请求,并回应投票结果
func (rf *Raft) HandleRequestVote(args *tinnraftpb.RequestVoteArgs, reply *tinnraftpb.RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	defer rf.persist()

	//发现候选人的term比自己小，拒绝投票
	if int(args.Term) < rf.currentTerm ||
		(args.Term == int64(rf.currentTerm) && rf.votedFor != -1 && rf.votedFor != int(args.CandidateId)) {
		DLog("[%v]: term %v | refuse vote for [%v]", rf.me, rf.currentTerm, args.CandidateId)
		reply.Term = int64(rf.currentTerm)
		reply.VoteGranted = false
		return
	}

	//发现候选人的term比自己大
	if int(args.Term) > rf.currentTerm {
		rf.ChangeRaftState(Follower)
		rf.currentTerm = int(args.Term)
		rf.votedFor = -1
	}

	follow_lastLog := rf.log.GetPersistLastEntry()

	//判断候选者日志是否是最新的
	upToDate := uint64(args.LastLogTerm) > follow_lastLog.Term ||
		(uint64(args.LastLogTerm) == follow_lastLog.Term &&
			args.LastLogIndex >= follow_lastLog.Index)

	if (rf.votedFor == -1 || rf.votedFor == int(args.CandidateId)) && upToDate && rf.me != int(args.CandidateId) {
		reply.VoteGranted = true
		rf.votedFor = int(args.CandidateId)

		rf.resetElectionTimer()

		//LOG
		raftlog := &tinnraftpb.LogArgs{
			Op:       tinnraftpb.LogOp_Vote,
			Contents: "vote for candidate",
			FromId:   strconv.Itoa(rf.me),
			ToId:     strconv.Itoa(rf.votedFor),
			CurState: "follower",
			Pid:      int64(syscall.Getpid()),
			Term:     int64(rf.currentTerm),
			Layer:    tinnraftpb.LogLayer_RAFT,
		}
		rf.apiGateClient.SendLogToGate(raftlog)

		DLog("[%v]: term %v | vote for [%v]", rf.me, rf.currentTerm, args.CandidateId)
	} else {
		reply.VoteGranted = false
		DLog("[%v]: term %v | refuse vote for [%v]", rf.me, rf.currentTerm, args.CandidateId)
	}
	reply.Term = int64(rf.currentTerm)
}

// func (rf *Raft) sendRequestVote(server int, args *tinnraftpb.RequestVoteArgs) (*tinnraftpb.RequestVoteReply, error) {
// 	//test 04-24
// 	return (*rf.peers[server].raftServiceCli).RequestVote(context.Background(), args)
// }
