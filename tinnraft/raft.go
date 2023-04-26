package tinnraft

import (
	"fmt"
	api_gateway "sakurajima-ds/api_gateway_2"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
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

	apiGateClient *api_gateway.ApiGatwayClient //api_server客户端

	state RaftState //raft状态机
	//appendEntryCh chan *tinnraftpb.Entry
	heartBeatTimer *time.Timer //心跳间隔时间
	electionTimer  *time.Timer //选举超时时间

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
	applyCh chan *tinnraftpb.ApplyMsg, apigateclient *api_gateway.ApiGatwayClient) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = MakePersister(dbEngine)
	rf.me = me

	rf.state = Follower
	rf.currentTerm = 0
	rf.votedFor = -1
	rf.heartBeatTimer = time.NewTimer(time.Millisecond * time.Duration(2000))
	rf.electionTimer = time.NewTimer(time.Microsecond * time.Duration(MakeAnRandomElectionTimeout()))

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

	rf.apiGateClient = apigateclient

	// initialize from state persisted before a crash
	rf.readPersist()

	fmt.Println("-----------------------------------")
	for _, peer := range peers {
		DLog("peer addr %s id %d", peer.addr, peer.id)
	}
	fmt.Println("-----------------------------------")

	// LOG
	raftlog := &tinnraftpb.LogArgs{
		Op:       tinnraftpb.LogOp_StartSucess,
		Contents: "start server success",
		FromId:   strconv.Itoa(rf.me),
		CurState: "follower",
		Pid:      int64(syscall.Getpid()),
		Term:     int64(rf.currentTerm),
		Layer:    tinnraftpb.LogLayer_RAFT,
	}
	rf.apiGateClient.SendLogToGate(raftlog)

	DLog("[%v]: term %v | the last log idx is %v", rf.me, rf.currentTerm, rf.log.GetPersistLastEntry().Index)

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

	if rf.isSnapshoting {
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
	//DLog("[%v]: term %v start %v", rf.me, rf.currentTerm, rf.log)
	return newEntry
}

func (rf *Raft) GetState() (int, bool) {

	rf.mu.Lock()
	defer rf.mu.Unlock()
	term := rf.currentTerm
	isleader := rf.state == Leader

	return term, isleader
}

// 获取当前节点状态/角色
func (rf *Raft) GetNowState() string {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return string(rf.state)
}

// 获取当前节点任期
func (rf *Raft) GetNowTerm() int {
	rf.mu.Lock()
	defer rf.mu.Lock()
	return rf.currentTerm
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
	// for !rf.IsKilled() {
	// 	time.Sleep(rf.heartBeat)
	// 	rf.mu.Lock()
	// 	if rf.state == Leader {
	// 		rf.appendEntries(true)
	// 	}
	// 	if time.Now().After(rf.electionTime) {
	// 		rf.leaderElection()
	// 	}
	// 	rf.mu.Unlock()
	// }

	for !rf.IsKilled() {
		select {
		case <-rf.electionTimer.C:
			{
				rf.leaderElection()
			}
		case <-rf.heartBeatTimer.C:
			{
				if rf.state == Leader {
					//DLog("hahahhahahah")
					rf.appendEntries(true)
					rf.resetHeartTimer()
				}
			}
		}
	}
}

// 唤醒 更新LastAppId协程
func (rf *Raft) apply() {
	rf.applyCond.Signal()
}

// 更新LastAppId
func (rf *Raft) applier() {
	for !rf.IsKilled() {
		rf.mu.Lock()
		for rf.lastApplied >= rf.commitIndex {
			DLog("[%v]: waiting for applier...", rf.me)
			//阻塞等待commitindex发生变化
			rf.applyCond.Wait()
		}

		if rf.commitIndex > rf.lastApplied && rf.log.GetPersistLastEntry().Index > int64(rf.lastApplied) {

			firstIndex := rf.log.GetPersistFirstEntry().Index
			commitIndex := rf.commitIndex
			lastApplied := rf.lastApplied
			entries := make([]*tinnraftpb.Entry, commitIndex-lastApplied)
			copy(entries, rf.log.GetPersistInterLog(int64(lastApplied)+1-int64(firstIndex), int64(commitIndex)+1-int64(firstIndex)))
			//DLog("[%v]: firstIndex = %v", rf.me, firstIndex)
			DLog("[%v]: term %v | applied entries from %d to %d ", rf.me, rf.currentTerm, rf.lastApplied, commitIndex)

			rf.mu.Unlock()

			//DLog("entries len: %d", len(entries))
			for _, entry := range entries {

				//向applyCh写入数据，提醒应用层将操作应用到状态机
				rf.applyCh <- &tinnraftpb.ApplyMsg{
					CommandValid: true,
					Command:      entry.Data,
					CommandIndex: int64(entry.Index),
					CommandTerm:  int64(entry.Term),
				}

			}
			rf.mu.Lock()
			rf.lastApplied = max(rf.lastApplied, commitIndex)
			rf.mu.Unlock()
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

func (rf *Raft) GetLeaderId() int64 {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return int64(rf.leaderId)
}

func (rf *Raft) ChangeRaftState(state RaftState) {
	if state == rf.state {
		return
	}
	rf.state = state
	DLog("[%v] term %v | change state to %v", rf.me, rf.currentTerm, state)

	switch state {
	case Follower:
		rf.heartBeatTimer.Stop()
		rf.resetElectionTimer()
	case Candidate:
	case Leader:
		lastLog := rf.log.GetPersistLastEntry()
		rf.leaderId = rf.me
		for i := 0; i < len(rf.peers); i++ {
			rf.matchIndex[i] = 0
			rf.nextIndex[i] = int(lastLog.Index) + 1
		}
		rf.electionTimer.Stop()
		rf.resetHeartTimer()
	}
}
