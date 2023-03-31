package tinnraft

import (
	"math/rand"
	"sakurajima-ds/tinnraftpb"
	"sync"
	"time"
)

// leader选举
func (rf *Raft) leaderElection() {
	rf.currentTerm++
	rf.ChangeRaftState(Candidate)
	rf.votedFor = rf.me

	rf.resetElectionTimer()
	rf.persist()
	
	term := rf.currentTerm
	voteCounter := 1

	candiate_lastLog := rf.log.GetPersistLastEntry()

	DLog("[%v]: term: %v | start leader election\n", rf.me, rf.currentTerm)

	args := tinnraftpb.RequestVoteArgs{
		Term:         int64(term),
		CandidateId:  int64(rf.me),
		LastLogIndex: int64(candiate_lastLog.Index),
		LastLogTerm:  int64(candiate_lastLog.Term),
	}

	//sync原语，表示becomeLeader对象只执行一次(类似单例)
	var becomeLeader sync.Once
	//向集群内其他server发送投票请求
	for _, peer := range rf.peers {
		if int(peer.id) != rf.me {
			go rf.candidateRequestVote(int(peer.id), &args, &voteCounter, &becomeLeader)
		}
	}
}

// 重置选举超时计时
func (rf *Raft) resetElectionTimer() {
	rf.electionTimer.Reset(time.Duration(MakeAnRandomElectionTimeout()) * time.Millisecond)
}

// 重置心跳超时计时
func (rf *Raft) resetHeartTimer() {
	rf.heartBeatTimer.Reset(time.Millisecond * time.Duration(2000))
}

func RandIntRange(min int, max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(max-min) + int(min)
}

func MakeAnRandomElectionTimeout() int {
	return RandIntRange(6000, 6000*2)
}

// 设置新的任期
func (rf *Raft) setNewTerm(term int) {
	if term > rf.currentTerm || rf.currentTerm == 0 {
		rf.state = Follower
		rf.currentTerm = term
		rf.votedFor = -1
		rf.persist()
	}
}
