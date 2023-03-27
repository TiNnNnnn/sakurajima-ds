package kv_server

import (
	"context"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"sync"
	"sync/atomic"
)

/*
kv存储服务器
*/

var DnsMap = map[int]string{
	0: "127.0.0.1:10010",
	1: "127.0.0.1:10011",
	2: "127.0.0.1:10012",
}

type OperationContext struct {
	MaxAppliedCommandId int64
	LastResponse        tinnraftpb.CommandReply
}

type KvServer struct {
	mu      sync.Mutex
	dead    int32
	tinnRf  *tinnraft.Raft
	applyCh chan *tinnraftpb.ApplyMsg

	lastApplied    int
	stm            StateMachine
	lastOperations map[int64]OperationContext

	notifyChans map[int]chan *tinnraftpb.CommandReply

	stopApplyCh chan interface{}

	tinnraftpb.UnimplementedRaftServiceServer
}

func MakeKvServer(serverId int) *KvServer {
	clientEnds := []*tinnraft.ClientEnd{}
	for id, addr := range DnsMap {
		newClient := tinnraft.MakeClientEnd(uint64(id), addr)
		clientEnds = append(clientEnds, newClient)
	}

	logEngine, err := storage_engine.MakeLevelDBKvStorage("./data/kv_server" + "/node_" + strconv.Itoa(serverId))
	if err != nil {
		panic(err)
	}

	newApplyCh := make(chan *tinnraftpb.ApplyMsg)

	newRaft := tinnraft.MakeRaft(clientEnds, serverId, logEngine, newApplyCh)

	kvserver := &KvServer{
		tinnRf:      newRaft,
		applyCh:     newApplyCh,
		dead:        0,
		lastApplied: 0,
		stm:         NewMenKv(),
		notifyChans: make(map[int]chan *tinnraftpb.CommandReply),
	}

	return kvserver
}

func (kvs *KvServer) restoreSnapshot(snapdata []byte) {

}

func (kvs *KvServer) takeSnapshot(index int) {

}

//实现 RequestVote interface
func RequestVote(ctx context.Context, in *tinnraftpb.RequestVoteArgs) (*tinnraftpb.RequestVoteReply, error) {

}

// 实现 AppendEntries interface
func AppendEntries(ctx context.Context, in *tinnraftpb.AppendEntriesArgs) (*tinnraftpb.AppendEntriesReply, error) {
	
}

// 实现 Snapshot interface
func Snapshot(ctx context.Context, in *tinnraftpb.InstallSnapshotArgs) (*tinnraftpb.InstallSnapshotReply, error) {

}

func (kvs *KvServer) IsKilled() bool {
	return atomic.LoadInt32(&kvs.dead) == 1
}

func (kvs *KvServer) ApplingToStm(done <-chan interface{}) {

}

func (kvs *KvServer) DoCommand(ctx context.Context, args *tinnraftpb.CommandArgs) *tinnraftpb.CommandReply {

}
