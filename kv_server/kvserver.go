package kv_server

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	api_gateway "sakurajima-ds/api_gateway_2"
	"sakurajima-ds/common"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

/*
kv存储服务器
*/

var DnsMap = map[int]string{
	0: "127.0.0.1:10010",
	1: "127.0.0.1:10011",
	2: "127.0.0.1:10012",
	3: "127.0.0.1:10013",
}

var PeersMap = map[int]string{
	0: ":10010",
	1: ":10011",
	2: ":10012",
	3: ":10013",
}

const ClientTimeOut = 1 * time.Second

type OperationContext struct {
	MaxAppliedCommandId int64
	LastResponse        tinnraftpb.CommandReply
}

type KvServer struct {
	mu      sync.Mutex
	dead    int32
	tinnRf  *tinnraft.Raft
	applyCh chan *tinnraftpb.ApplyMsg

	lastApplied int
	stm         StateMachine

	//lastOperations map[int64]OperationContext

	notifyChans map[int]chan *tinnraftpb.CommandReply //从ApplingToStm协程中获取操作结果

	stopApplyCh chan interface{}

	tinnraftpb.UnimplementedRaftServiceServer
}

func (kvs *KvServer) GetTinnRaft() *tinnraft.Raft {
	return kvs.tinnRf
}

// 初始化一个kvserver
func MakeKvServer(serverId int) *KvServer {
	clientEnds := []*tinnraft.ClientEnd{}

	//创建三个rpc客户端,并添加到clientEnds数组
	for id, addr := range DnsMap {
		newClient := tinnraft.MakeClientEnd(uint64(id), addr)
		clientEnds = append(clientEnds, newClient)
	}

	//创建levelDB存储引擎
	logEngine, err := storage_engine.MakeLevelDBKvStorage("./data/kv_server" + "/node_" + strconv.Itoa(serverId))
	if err != nil {
		panic(err)
	}

	//创建管道，应用层监听raft提交的操作与数据信息(CommandArgs)
	newApplyCh := make(chan *tinnraftpb.ApplyMsg)

	apigateclient := api_gateway.MakeApiGatwayClient(99, "127.0.0.1:10030")

	//实例化raft模块
	newRaft := tinnraft.MakeRaft(clientEnds, serverId, logEngine, newApplyCh, apigateclient)

	//实例化kvserver
	kvserver := &KvServer{
		tinnRf:      newRaft,
		applyCh:     newApplyCh,
		dead:        0,
		lastApplied: 0,
		stm:         NewMenKv(),
		notifyChans: make(map[int]chan *tinnraftpb.CommandReply),
	}

	kvserver.stopApplyCh = make(chan interface{})

	go kvserver.ApplingToStm(kvserver.stopApplyCh)
	return kvserver
}

// 读取快照并序列化到状态机
func (kvs *KvServer) restoreSnapshot(snapdata []byte) {
	if snapdata == nil {
		return
	}
	buf := bytes.NewBuffer(snapdata)
	data := gob.NewDecoder(buf)
	var tmpStm MemKv
	//将快照数据反序列化存储到tmpStm
	if data.Decode(&tmpStm) != nil {
		tinnraft.DLog("decode stm failed")
	}
	kvs.stm = &tmpStm
}

// 序列化后的状态机转为快照
func (kvs *KvServer) takeSnapshot(index int) {
	var bytesState bytes.Buffer
	enc := gob.NewEncoder(&bytesState)
	enc.Encode(kvs.stm)
	kvs.tinnRf.Snapshot(index, bytesState.Bytes())
}

// 实现 RequestVote interface
func (kvs *KvServer) RequestVote(ctx context.Context, args *tinnraftpb.RequestVoteArgs) (*tinnraftpb.RequestVoteReply, error) {
	reply := &tinnraftpb.RequestVoteReply{}
	kvs.tinnRf.HandleRequestVote(args, reply)
	return reply, nil
}

// 实现 AppendEntries interface
func (kvs *KvServer) AppendEntries(ctx context.Context, args *tinnraftpb.AppendEntriesArgs) (*tinnraftpb.AppendEntriesReply, error) {
	reply := &tinnraftpb.AppendEntriesReply{}
	kvs.tinnRf.HandleAppendEntries(args, reply)
	return reply, nil
}

// 实现 Snapshot interface
func (kvs *KvServer) Snapshot(ctx context.Context, args *tinnraftpb.InstallSnapshotArgs) (*tinnraftpb.InstallSnapshotReply, error) {
	reply := &tinnraftpb.InstallSnapshotReply{}
	kvs.tinnRf.HandleInstallSnapshot(args, reply)
	return reply, nil
}

func (kvs *KvServer) IsKilled() bool {
	return atomic.LoadInt32(&kvs.dead) == 1
}

func (kvs *KvServer) ReadNotifyChan(idx int) chan *tinnraftpb.CommandReply {
	_, ok := kvs.notifyChans[idx]
	if !ok {
		kvs.notifyChans[idx] = make(chan *tinnraftpb.CommandReply, 1)
	}
	return kvs.notifyChans[idx]
}

// 将提交的数据应用到状态机
func (kvs *KvServer) ApplingToStm(done <-chan interface{}) {
	for !kvs.IsKilled() {
		select {
		case <-done: //退出select
			return
		case appliedMsg := <-kvs.applyCh: //applyCh监听并接受 raft提交的数据
			if appliedMsg.CommandValid {
				kvs.mu.Lock()
				args := &tinnraftpb.CommandArgs{}
				//将Command信息反序列化到结构体args中
				if err := json.Unmarshal(appliedMsg.Command, args); err != nil {
					kvs.mu.Unlock()
					tinnraft.DLog("Umarshalal ComandArgs failed")
					continue
				}
				if appliedMsg.CommandIndex <= int64(kvs.lastApplied) {
					kvs.mu.Unlock()
					continue
				}

				//更新lastApplied
				kvs.lastApplied = int(appliedMsg.CommandIndex)

				var value string
				switch args.OpType {
				case tinnraftpb.OpType_Put:
					kvs.stm.Put(args.Key, args.Value)
				case tinnraftpb.OpType_Get:
					value, _ = kvs.stm.Get(args.Key)
				case tinnraftpb.OpType_Append:
					kvs.stm.Append(args.Key, args.Value)
				}

				//构建对客户端的响应
				reply := &tinnraftpb.CommandReply{}
				reply.Value = value

				//通过NotifyChan通知客户端 操作已完成
				ch := kvs.ReadNotifyChan(int(appliedMsg.CommandIndex))
				ch <- reply

				//当日志数量大于10,进行快照,压缩空间
				if kvs.tinnRf.LogCount() > 10 {
					kvs.takeSnapshot(int(appliedMsg.CommandIndex))
				}
				kvs.mu.Unlock()

			} else if appliedMsg.SnapshotValid {
				kvs.mu.Lock()
				if kvs.tinnRf.CondInstallSnapshot(int(appliedMsg.SnapshotTerm), int(appliedMsg.SnapshotIndex), appliedMsg.Snapshot) {
					kvs.restoreSnapshot(appliedMsg.Snapshot)
					kvs.lastApplied = int(appliedMsg.SnapshotIndex)
				}
				kvs.mu.Unlock()
			}

		}
	}
}

func (kvs *KvServer) DoCommand(ctx context.Context, args *tinnraftpb.CommandArgs) (*tinnraftpb.CommandReply, error) {

	reply := &tinnraftpb.CommandReply{}

	if args != nil {
		//将客户段消息序列化为json
		argBytes, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}
		//将json消息发给raft模块处理，且只有leader节点接受消息
		idx, _, isLeader := kvs.tinnRf.Propose(argBytes)
		if !isLeader {
			reply.ErrCode = common.ErrCodeWrongLeader
			return reply, nil
		}

		//获取idx对应的notifychan管道
		kvs.mu.Lock()
		ch := kvs.ReadNotifyChan(idx)
		kvs.mu.Unlock()

		select {
		case res := <-ch:
			reply.Value = res.Value
		case <-time.After(ClientTimeOut):
			reply.ErrCode = common.ErrCodeExecTimeout
			reply.Value = "timeout"
		}

		go func() {
			kvs.mu.Lock()
			delete(kvs.notifyChans, idx)
			kvs.mu.Unlock()
		}()
	}
	return reply, nil
}
