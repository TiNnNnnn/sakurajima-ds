package tinnraft

import (
	"sakurajima-ds/tinnraftpb"
	"sync"
	_ "sync/atomic"
	"time"
)

type RaftState string

const (
	Follower  RaftState = "Follower"
	Candidate           = "Candidate"
	Leader              = "Leader"
)

type Raft struct {
	mu    sync.Mutex
	peers []*ClientEnd
	me    int
	dead  int32

	state         RaftState //raft状态机
	appendEntryCh chan *tinnraftpb.Entry
	heartBeat     time.Duration //心跳间隔时间
	electionTime  time.Time     //选举超时时间

	//通用持久化状态（所有server）
	currentTerm int //当前任期
	votedFor    int //记录将选票投给了谁
	log         Log //日志

	//通用易失性状态 (所有server)
	commitIndex int //已提交的日志下标
	lastApplied int //最后一条已经应用于状态机的log索引条目

	//通用易失性状态 (Leader)
	nextIndex  []int //记录 发送到follwer服务器的下一条日志条目的索引
	matchIndex []int //记录 发送到floower服务器的已知的已经复制到该服务器的最高日志的索引

	applyCh   chan tinnraftpb.ApplyMsg
	applyCond *sync.Cond
}

func Make(peers []*ClientEnd, me int,
	persister *Persister, applyCh chan tinnraftpb.ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).
	rf.state = Follower
	rf.currentTerm = 0
	rf.votedFor = -1
	rf.heartBeat = 50 * time.Millisecond
	rf.resetElectionTimer()

	//日志初始化
	rf.log = makeEmptyLog()
	rf.log.append(tinnraftpb.Entry{0, 0, 0})

	rf.commitIndex = 0
	rf.lastApplied = 0
	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))

	rf.applyCh = applyCh
	rf.applyCond = sync.NewCond(&rf.mu)

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	//开启一个协程进行选举
	go rf.ticker()

	//开启一个协程去进行更新LastAppId
	go rf.applier()

	return rf
}
