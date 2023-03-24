package tinnraft

import (
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraftpb"
	"sync"
	"sync/atomic"
	"time"
)

type RaftState string

const (
	Follower  RaftState = "Follower"
	Candidate           = "Candidate"
	Leader              = "Leader"
)

type Raft struct {
	mu        sync.Mutex
	peers     []*ClientEnd
	persister *Log
	me        int
	dead      int32

	state RaftState //raft状态机
	//appendEntryCh chan *tinnraftpb.Entry
	heartBeat    time.Duration //心跳间隔时间
	electionTime time.Time     //选举超时时间

	//通用持久化状态（所有server）
	currentTerm int  //当前任期
	votedFor    int  //记录将选票投给了谁
	log         *Log //日志

	//通用易失性状态 (所有server)
	commitIndex int //已提交的日志下标
	lastApplied int //最后一条已经应用于状态机的log索引条目

	//通用易失性状态 (Leader)
	nextIndex  []int //记录 发送到follwer服务器的下一条日志条目的索引
	matchIndex []int //记录 发送到floower服务器的已知的已经复制到该服务器的最高日志的索引

	applyCh   chan tinnraftpb.ApplyMsg
	applyCond *sync.Cond
}

// 初始化一个raft主机
func MakeRaft(peers []*ClientEnd, me int, dbEngine storage_engine.KvStorage,
	persister *Persister, applyCh chan tinnraftpb.ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = MakePersister(dbEngine)
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).
	rf.state = Follower
	rf.currentTerm = 0
	rf.votedFor = -1
	rf.heartBeat = 50 * time.Millisecond
	rf.resetElectionTimer()

	//日志初始化,加入一个空日志
	rf.log = makeEmptyLog()

	rf.commitIndex = 0
	rf.lastApplied = 0
	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))

	rf.applyCh = applyCh
	rf.applyCond = sync.NewCond(&rf.mu)

	// initialize from state persisted before a crash
	rf.readPersist()

	//开启一个协程进行选举
	go rf.ticker()

	//开启一个协程去进行更新LastAppId
	go rf.applier()

	return rf
}

// 设置主机状态为失活
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
}

// 判断主机是否失活
func (rf *Raft) IsKilled() bool {
	return atomic.LoadInt32(&rf.dead) == 1
}

func (rf *Raft) Propose(payload []byte) (int, int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.state != Leader {
		return -1, -1, false
	}
	index := rf.log.lastLog().Index + 1
	term := rf.currentTerm

	log := tinnraftpb.Entry{
		Index: index,
		Term:  uint64(term),
		Data:  payload,
	}

	rf.log.append2(log)
	rf.persist()
	//fmt.Printf("[%v]: term %v Start %v", rf.me, term, log)
	//DPrintf("[%v]: term %v Start %v", rf.me, term, log)
	rf.appendEntries(false)

	return int(index), term, true
}

func (rf *Raft) GetState() (int, bool) {

	// Your code here (2A).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	term := rf.currentTerm
	isleader := rf.state == Leader

	return term, isleader
}

// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
func (rf *Raft) persist() {
	rf.persister.PersistRaftState(int64(rf.currentTerm), int64(rf.votedFor))
}

// restore previously persisted state.
func (rf *Raft) readPersist() {
	curtTerm, votedFor := rf.persister.ReadRaftState()
	rf.currentTerm = int(curtTerm)
	rf.votedFor = int(votedFor)
}

// 选举与心跳触发器
func (rf *Raft) ticker() {
	for !rf.IsKilled() {
		time.Sleep(rf.heartBeat)
		rf.mu.Lock()
		if rf.state == Leader {
			rf.appendEntries(true)
		}
		if time.Now().After(rf.electionTime) {
			rf.leaderElection()
		}
		rf.mu.Unlock()
	}
}

// 唤醒 更新LastAppId协程
func (rf *Raft) apply() {
	rf.applyCond.Broadcast()
}

// 更新LastAppId
func (rf *Raft) applier() {
	rf.mu.Lock()
	defer rf.mu.Lock()

	for !rf.IsKilled() {
		if rf.commitIndex > rf.lastApplied && rf.log.lastLog().Index > int64(rf.lastApplied) {
			rf.lastApplied++
			applyMsg := tinnraftpb.ApplyMsg{
				CommandValid: true,
				Command:      rf.log.at(rf.lastApplied).Data,
				CommandIndex: int64(rf.lastApplied),
			}
			rf.mu.Unlock()
			rf.applyCh <- applyMsg
			rf.mu.Lock()
		} else {
			//阻塞等待commitindex发生变化
			rf.applyCond.Wait()
		}
	}
}
