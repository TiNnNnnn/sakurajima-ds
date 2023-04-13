package meta_server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sakurajima-ds/common"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type MetaServer struct {
	mu          sync.RWMutex
	dead        int32
	tinnRf      *tinnraft.Raft
	applyCh     chan *tinnraftpb.ApplyMsg
	stm         MetaStateMachine
	notifyChans map[int]chan *tinnraftpb.MetaReply
	stopApplyCh chan interface{}

	tinnraftpb.UnimplementedRaftServiceServer
}

func MakeMetaServer(peerMaps map[int]string, serverId int) *MetaServer {
	clientEnds := []*tinnraft.ClientEnd{}

	//创建三个rpc客户端,并添加到clientEnds数组
	for id, addr := range peerMaps {
		newClient := tinnraft.MakeClientEnd(uint64(id), addr)
		clientEnds = append(clientEnds, newClient)
	}
	//创建管道，应用层监听raft提交的操作与数据信息(ConfigArgs)
	applyCh := make(chan *tinnraftpb.ApplyMsg)

	Metaengine := storage_engine.EngineFactory("leveldb", "./meta_data"+"/node_"+strconv.Itoa(serverId))

	logEngine := storage_engine.EngineFactory("leveldb", "./log_data/"+"metaserver/node_"+strconv.Itoa(serverId))

	tinnRf := tinnraft.MakeRaft(clientEnds, serverId, logEngine, applyCh)

	metaServer := &MetaServer{
		dead:        0,
		applyCh:     applyCh,
		notifyChans: make(map[int]chan *tinnraftpb.MetaReply),
		stm:         *MakeMetaStm(Metaengine),
		tinnRf:      tinnRf,
	}

	metaServer.stopApplyCh = make(chan interface{})

	go metaServer.ApplingToStm(metaServer.stopApplyCh)

	return metaServer
}

func (cs *MetaServer) IsKilled() bool {
	return atomic.LoadInt32(&cs.dead) == 1
}

func (cs *MetaServer) getNotifyChan(index int) chan *tinnraftpb.MetaReply {
	if _, ok := cs.notifyChans[index]; !ok {
		cs.notifyChans[index] = make(chan *tinnraftpb.MetaReply, 1)
	}
	return cs.notifyChans[index]
}

// 实现 RequestVote interface
func (cs *MetaServer) RequestVote(ctx context.Context, args *tinnraftpb.RequestVoteArgs) (*tinnraftpb.RequestVoteReply, error) {
	reply := &tinnraftpb.RequestVoteReply{}
	cs.tinnRf.HandleRequestVote(args, reply)
	return reply, nil
}

// 实现 AppendEntries interface
func (cs *MetaServer) AppendEntries(ctx context.Context, args *tinnraftpb.AppendEntriesArgs) (*tinnraftpb.AppendEntriesReply, error) {
	reply := &tinnraftpb.AppendEntriesReply{}
	cs.tinnRf.HandleAppendEntries(args, reply)
	return reply, nil
}

func (ms *MetaServer) ApplingToStm(done <-chan interface{}) {
	for !ms.IsKilled() {
		select {
		case <-done:
			return
		case appliedMsg := <-ms.applyCh:
			args := &tinnraftpb.MetaArgs{}
			err := json.Unmarshal(appliedMsg.Command, args)
			if err != nil {
				tinnraft.DLog("Unmarshal MetaArgs failed")
				continue
			}

			reply := &tinnraftpb.MetaReply{}
			switch args.OpType {
			case tinnraftpb.MetaOpType_PutName: //putName
				bucketName := args.BucketName
				objectName := args.ObjectName
				objectId := args.ObjectId
				blockId, err := ms.stm.PutName(bucketName, objectName, objectId)
				if err != nil {
					log.Fatalf("put the key: %v to NameTable faield", bucketName+objectName)
					reply.ErrMsg = err.Error()
					return
				}
				reply.BlockId = blockId
				reply.ErrCode = common.ErrCodeNoErr
			case tinnraftpb.MetaOpType_PutObject: //putObject
				objecId := args.ObjectId
				blockList := args.BlockList
				bid, err := ms.stm.PutObject(objecId, blockList)
				if err != nil {
					log.Fatal("put the blocklist to ObjectTable field")
					reply.ErrMsg = err.Error()
					return
				}
				reply.BlockId = bid
				reply.ErrCode = common.ErrCodeNoErr
			case tinnraftpb.MetaOpType_GetObjectList: //GetObjectLIst
				key := args.BucketName + args.ObjectName
				num := args.Num
				list, err := ms.stm.GetBlockList(key, int(num))
				if err != nil {
					log.Fatal("get the blocklist from metaserver faield")
					reply.ErrMsg = err.Error()
					return
				}
				reply.BlockIdList = list
				reply.ErrCode = common.ErrCodeNoErr
			}
			ch := ms.getNotifyChan(int(appliedMsg.CommandIndex))
			ch <- reply
		}
	}
}

func (ms *MetaServer) DoMeta(ctx context.Context, args *tinnraftpb.MetaArgs) (*tinnraftpb.MetaReply, error) {
	reply := &tinnraftpb.MetaReply{}

	argsBytes, err := json.Marshal(args)
	if err != nil {
		reply.ErrMsg = err.Error()
		return reply, err
	}

	idx, _, isLeader := ms.tinnRf.Propose(argsBytes)
	if !isLeader {
		reply.ErrMsg = "not Leader node"
		reply.ErrCode = common.ErrCodeWrongLeader
		reply.LeaderId = ms.tinnRf.GetLeaderId()
		return reply, nil
	}

	ms.mu.Lock()
	ch := ms.getNotifyChan(idx)
	ms.mu.Unlock()

	select {
	case res := <-ch:
		reply.BlockId = res.BlockId
		reply.BlockIdList = res.BlockIdList
		reply.ErrMsg = res.ErrMsg
		reply.ErrCode = common.ErrCodeNoErr
	case <-time.After(3 * time.Second):
		reply.ErrMsg = "server timeout"
		reply.ErrCode = common.ErrCodeExecTimeout
		return reply, errors.New("ServerTimerout")
	}

	go func() {
		ms.mu.Lock()
		delete(ms.notifyChans, idx)
		ms.mu.RUnlock()
	}()

	return reply, nil

}