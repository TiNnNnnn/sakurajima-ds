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
	Candidate RaftState = "Candidate"
	Leader    RaftState = "Leader"
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

	isSnapshoting bool

	applyCh   chan *tinnraftpb.ApplyMsg
	applyCond *sync.Cond

	leaderId int
}

// 初始化一个raft主机
func MakeRaft(peers []*ClientEnd, me int, dbEngine storage_engine.KvStorage,
	applyCh chan *tinnraftpb.ApplyMsg) *Raft {
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
	rf.log = MakePersister(dbEngine)

	rf.commitIndex = 0
	rf.lastApplied = 0
	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))

	rf.isSnapshoting = false

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

/*
用户请求到来和raft交互的入口函数是Propose，这个函数首先会查询当前节点状态，
只有Leader节点才能处理提案（propose），之后会把用户操作的序列化之后的[]byte调用
Append追加到自己的日志中，之后appendEntries将日志内容发送给集群中的Follower节点
*/
func (rf *Raft) Propose(payload []byte) (int, int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.state != Leader {
		return -1, -1, false
	}

	newEntry := rf.AppendNewCommand(payload)
	rf.appendEntries(false)

	return int(newEntry.Index), int(newEntry.Term), true
}

// 添加一条新的日志到Leader的日志中
func (rf *Raft) AppendNewCommand(command []byte) *tinnraftpb.Entry {
	lastLog := rf.log.GetPersistLastEntry()
	newEntry := &tinnraftpb.Entry{
		Index: lastLog.Index + 1,
		Term:  uint64(rf.currentTerm),
		Data:  command,
	}
	rf.log.PersistAppend(newEntry)
	rf.persister.PersistRaftState(int64(rf.currentTerm), int64(rf.votedFor))
	return newEntry
}

func (rf *Raft) GetState() (int, bool) {

	rf.mu.Lock()
	defer rf.mu.Unlock()
	term := rf.currentTerm
	isleader := rf.state == Leader

	return term, isleader
}

// 将状态持久化
func (rf *Raft) persist() {
	rf.persister.PersistRaftState(int64(rf.currentTerm), int64(rf.votedFor))
}

// 恢复到之前的持久化状态
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
	for !rf.IsKilled() {
		rf.mu.Lock()

		if rf.commitIndex > rf.lastApplied && rf.log.GetPersistLastEntry().Index > int64(rf.lastApplied) {

			firstIndex := rf.log.FirstLogIdx()
			commitIndex := rf.commitIndex
			lastApplied := rf.lastApplied
			entries := make([]*tinnraftpb.Entry, commitIndex-lastApplied)
			copy(entries, rf.log.GetPersistInterLog(int64(lastApplied)+1-int64(firstIndex), int64(commitIndex)-int64(firstIndex)))

			rf.mu.Unlock()
			for _, entry := range entries {
				applyMsg := tinnraftpb.ApplyMsg{
					CommandValid: true,
					Command:      entry.Data,
					CommandIndex: int64(entry.Index),
					CommandTerm:  int64(entry.Term),
				}
				//向applyCh写入数据，提醒应用层将操作应用到状态机
				rf.applyCh <- &applyMsg
			}
			rf.mu.Lock()
			rf.lastApplied = max(rf.lastApplied, commitIndex)
			rf.mu.Unlock()
		} else {
			//阻塞等待commitindex发生变化
			rf.applyCond.Wait()
		}
	}
}

// 获取日志长度
func (rf *Raft) LogCount() int {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.log.LogPersistLen()
}

// 关闭所有rpc连接
func (rf *Raft) CloseAllConn() {
	for _, peer := range rf.peers {
		peer.CloseConns()
	}
}
