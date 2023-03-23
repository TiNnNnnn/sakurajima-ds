package tinnraft

import (
	"fmt"
	"math/rand"
	"sakurajima-ds/tinnraftpb"
	"sync"
	"time"
)

// leader选举
func (rf *Raft) leaderElection() {
	rf.currentTerm++
	rf.state = Candidate
	rf.votedFor = rf.me

	rf.persist()
	rf.resetElectionTimer()

	term := rf.currentTerm
	voteCounter := 1

	candiate_lastLog := rf.log.lastLog()

	fmt.Printf("[%v]: start leader election, term %d\n", rf.me, rf.currentTerm)

	args := tinnraftpb.RequestVoteArgs{
		Term:         int64(term),
		CandidateId:  int64(rf.me),
		LastLogIndex: int64(candiate_lastLog.Index),
		LastLogTerm:  int64(candiate_lastLog.Term),
	}
	//sync原语，表示becomeLeader对象只执行一次(类似单例)
	var becomeLeader sync.Once
	//向集群内其他server发送投票请求
	for serverId := range rf.peers {
		if serverId == rf.me {
			continue
		}
		go rf.candidateRequestVote(serverId, &args, &voteCounter, &becomeLeader)
	}
}

// 重置选举超时计时
func (rf *Raft) resetElectionTimer() {
	t := time.Now()
	electionTimeout := time.Duration(150+rand.Intn(150)) * time.Millisecond
	rf.electionTime = t.Add(electionTimeout)
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
