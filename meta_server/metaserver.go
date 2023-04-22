package meta_server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
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

	apigateclient := api_gateway.MakeApiGatwayClient(99, "127.0.0.1:10030")

	tinnRf := tinnraft.MakeRaft(clientEnds, serverId, logEngine, applyCh, apigateclient)

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

func (kvs *MetaServer) GetTinnRaft() *tinnraft.Raft {
	return kvs.tinnRf
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

			tinnraft.DLog("appliingtostm recvie applyMsg: %v", args)

			reply := &tinnraftpb.MetaReply{}
			switch args.Op {
			case tinnraftpb.MetaOpType_PutMetadata: //put
				userId := args.UserId
				bucketName := args.BucketName
				objectName := args.ObjectName
				version := args.Version
				size := args.Size
				hash := args.Hash

				err := ms.stm.PutMetadata(userId, bucketName, objectName, int(version), size, hash)
				if err != nil {
					log.Printf("put the key: %v to NameTable faield", bucketName+objectName)
					reply.ErrMsg = err.Error()
					continue
				}
				log.Printf("put the key: %v to NameTable success", bucketName+objectName)
				reply.ErrCode = common.ErrCodeNoErr
			case tinnraftpb.MetaOpType_GetMetadata: //get
				userId := args.UserId
				bucketName := args.BucketName
				objectName := args.ObjectName
				version := args.Version
				meta, err := ms.stm.GetMetadata(userId, bucketName, objectName, int(version))
				if err != nil {
					log.Printf("get the meta server failed")
					reply.ErrMsg = err.Error()
					continue
				}
				m := tinnraftpb.Metadata{
					UserId:     meta.UserId,
					BucketName: meta.BucketName,
					ObjectName: meta.ObjectName,
					Version:    meta.Version,
					Size:       meta.Size,
					Hash:       meta.Hash,
				}
				reply.Meta = &m
				reply.ErrCode = common.ErrCodeNoErr
			case tinnraftpb.MetaOpType_SearchLastestVersion: //SearchLastestVersion
				userId := args.UserId
				bucketName := args.BucketName
				objectName := args.ObjectName
				meta, err := ms.stm.SearchLastestVersion(userId, bucketName, objectName)
				if err != nil {
					log.Printf("search lastest version failed")
					reply.ErrMsg = err.Error()
					continue
				}
				m := tinnraftpb.Metadata{
					UserId:     meta.UserId,
					BucketName: meta.BucketName,
					ObjectName: meta.ObjectName,
					Version:    meta.Version,
					Size:       meta.Size,
					Hash:       meta.Hash,
				}
				reply.Meta = &m
				reply.ErrCode = common.ErrCodeNoErr

			case tinnraftpb.MetaOpType_SearchAllVersion: //SearchAllVersion
				userId := args.UserId
				bucketName := args.BucketName
				objectName := args.ObjectName
				metaList, err := ms.stm.SearvhAllVersion(userId, bucketName, objectName, 0, 0)
				if err != nil {
					log.Printf("search all version failed")
					reply.ErrMsg = err.Error()
					continue
				}
				for _, meta := range metaList {
					m := tinnraftpb.Metadata{
						UserId:     meta.UserId,
						BucketName: meta.BucketName,
						ObjectName: meta.ObjectName,
						Version:    meta.Version,
						Size:       meta.Size,
						Hash:       meta.Hash,
					}
					reply.Metas = append(reply.Metas, &m)
				}
			case tinnraftpb.MetaOpType_AddVersion: //AddVersion
				userId := args.UserId
				bucketName := args.BucketName
				objectName := args.ObjectName
				size := args.Size
				hash := args.Hash
				err := ms.stm.AddVersion(userId, bucketName, objectName, hash, size)
				if err != nil {
					log.Printf("add version failed")
					reply.ErrMsg = err.Error()
					continue
				}
				log.Printf("add version success")
			case tinnraftpb.MetaOpType_DelMetadata:
				userId := args.UserId
				bucketName := args.BucketName
				objectName := args.ObjectName
				version := args.Version

				err := ms.stm.DelMetadata(userId, bucketName, objectName, int(version))
				if err != nil {
					log.Printf("delete failed")
					reply.ErrMsg = err.Error()
					continue
				}
				log.Printf("delete success")
			}
			ch := ms.getNotifyChan(int(appliedMsg.CommandIndex))
			ch <- reply
		}
	}
}

func (ms *MetaServer) DoMeta(ctx context.Context, args *tinnraftpb.MetaArgs) (*tinnraftpb.MetaReply, error) {
	tinnraft.DLog("DoMeta: args: %s", args.String())

	reply := &tinnraftpb.MetaReply{}

	argsBytes, err := json.Marshal(args)
	if err != nil {
		reply.ErrMsg = err.Error()
		return reply, err
	}

	//向raft层提交args
	idx, _, isLeader := ms.tinnRf.Propose(argsBytes)
	if !isLeader {
		reply.ErrMsg = "not Leader node"
		reply.ErrCode = common.ErrCodeWrongLeader
		reply.LeaderId = ms.tinnRf.GetLeaderId()
		return reply, nil
	}

	//等待获取stm层的响应
	ms.mu.Lock()
	ch := ms.getNotifyChan(idx)
	ms.mu.Unlock()

	select {
	case res := <-ch:
		reply.Meta = res.Meta
		reply.Metas = res.Metas
		reply.ErrMsg = res.ErrMsg
		reply.ErrCode = res.ErrCode
	case <-time.After(3 * time.Second):
		reply.ErrMsg = "server timeout"
		reply.ErrCode = common.ErrCodeExecTimeout
		return reply, errors.New("ServerTimerout")
	}

	go func() {
		ms.mu.Lock()
		delete(ms.notifyChans, idx)
		ms.mu.Unlock()
	}()

	return reply, nil

}
