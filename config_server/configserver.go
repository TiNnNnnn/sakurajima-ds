package config_server

import (
	"context"
	"encoding/json"
	"errors"
	api_gateway "sakurajima-ds/api_gateway"
	"sakurajima-ds/common"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type ConfigServer struct {
	mu            sync.RWMutex
	dead          int32
	tinnRf        *tinnraft.Raft
	applyCh       chan *tinnraftpb.ApplyMsg
	stm           ConfigStateMachine
	notifyChans   map[int]chan *tinnraftpb.ConfigReply
	stopApplyCh   chan interface{}
	apiGateClient *api_gateway.ApiGatwayClient

	tinnraftpb.UnimplementedRaftServiceServer
}

func MakeConfigServer(peerMaps map[int]string, serverId int) *ConfigServer {
	clientEnds := []*tinnraft.ClientEnd{}

	//创建三个rpc客户端,并添加到clientEnds数组
	for id, addr := range peerMaps {
		newClient := tinnraft.MakeClientEnd(uint64(id), addr)
		clientEnds = append(clientEnds, newClient)
	}
	//创建管道，应用层监听raft提交的操作与数据信息(ConfigArgs)
	applyCh := make(chan *tinnraftpb.ApplyMsg)

	Configengine := storage_engine.EngineFactory("leveldb", "../SALOG/conf_data"+"/node_"+strconv.Itoa(serverId))

	logEngine := storage_engine.EngineFactory("leveldb", "../SALOG/log_data/"+"configserver/node_"+strconv.Itoa(serverId))

	apigateclient := api_gateway.MakeApiGatwayClient(99, "127.0.0.1:10030")

	tinnRf := tinnraft.MakeRaft(clientEnds, serverId, logEngine, applyCh, apigateclient)

	configServer := &ConfigServer{
		dead:          0,
		applyCh:       applyCh,
		notifyChans:   make(map[int]chan *tinnraftpb.ConfigReply),
		stm:           *MakeConfigStm(Configengine),
		tinnRf:        tinnRf,
		apiGateClient: apigateclient,
	}

	configServer.stopApplyCh = make(chan interface{})

	go configServer.ApplingToStm(configServer.stopApplyCh)

	return configServer
}

func (cs *ConfigServer) IsKilled() bool {
	return atomic.LoadInt32(&cs.dead) == 1
}

func (s *ConfigServer) GetRaftNodesStat() string {
	return s.tinnRf.GetNowState()
}

func (s *ConfigServer) GetRaftNodeTerm() int {
	return s.tinnRf.GetNowTerm()
}

func (s *ConfigServer) StopApply() {
	close(s.applyCh)
}

func (s *ConfigServer) GetRaft() *tinnraft.Raft {
	return s.tinnRf
}

func (cs *ConfigServer) getNotifyChan(index int) chan *tinnraftpb.ConfigReply {
	if _, ok := cs.notifyChans[index]; !ok {
		cs.notifyChans[index] = make(chan *tinnraftpb.ConfigReply, 1)
	}
	return cs.notifyChans[index]
}

// 实现 RequestVote interface
func (cs *ConfigServer) RequestVote(ctx context.Context, args *tinnraftpb.RequestVoteArgs) (*tinnraftpb.RequestVoteReply, error) {
	reply := &tinnraftpb.RequestVoteReply{}
	cs.tinnRf.HandleRequestVote(args, reply)
	return reply, nil
}

// 实现 AppendEntries interface
func (cs *ConfigServer) AppendEntries(ctx context.Context, args *tinnraftpb.AppendEntriesArgs) (*tinnraftpb.AppendEntriesReply, error) {
	reply := &tinnraftpb.AppendEntriesReply{}
	cs.tinnRf.HandleAppendEntries(args, reply)
	return reply, nil
}

func (cs *ConfigServer) ApplingToStm(done <-chan interface{}) {
	for !cs.IsKilled() {
		select {
		case <-done:
			return
		case appliedMsg := <-cs.applyCh:
			args := &tinnraftpb.ConfigArgs{}
			uerr := json.Unmarshal(appliedMsg.Command, args)
			if uerr != nil {
				tinnraft.DLog("Unmarshal ConfigArgs failed")
				continue
			}
			var conf Config
			var err error

			reply := &tinnraftpb.ConfigReply{}
			switch args.OpType {
			case tinnraftpb.ConfigOpType_Join: //Join操作
				groups := map[int][]string{}
				for groupId, addrs := range args.Servers {
					groups[int(groupId)] = strings.Split(addrs, ",")
				}
				err = cs.stm.Join(groups)

			case tinnraftpb.ConfigOpType_Leave: //Leave操作
				groupIds := []int{}
				for _, id := range args.GroupIds {
					groupIds = append(groupIds, int(id))
				}
				err = cs.stm.Leave(groupIds)
			case tinnraftpb.ConfigOpType_Move: //Move操作
				conf, err = cs.stm.Query(int(args.ConfigVersion))
				if err != nil {
					reply.ErrMsg = err.Error()
				}
				// //judge isneeded to migrate buckets(slots)
				// if conf.Buckets[args.BucketId] != 0 && conf.Buckets[args.BucketId] != int(args.GroupId) {
				// 	reply.Mb = &tinnraftpb.MigrateBucket{
				// 		BucketId: args.BucketId,
				// 		From:     int64(conf.Buckets[args.BucketId]),
				// 		To:       args.GroupId,
				// 	}
				// 	reply.Ismigrate = true
				// }
				err = cs.stm.Move(int(args.BucketId), int(args.GroupId))
			case tinnraftpb.ConfigOpType_Query: //Query操作
				conf, err = cs.stm.Query(int(args.ConfigVersion))
				if err != nil {
					reply.ErrMsg = err.Error()
				}
				confBytes, err := json.Marshal(conf)
				if err != nil {
					reply.ErrMsg = err.Error()
				}
				tinnraft.DLog("query configs: %v", string(confBytes))

				reply.Config = &tinnraftpb.ServerConfig{}
				reply.Config.ConfigVersion = int64(conf.Version)
				for _, b := range conf.Buckets {
					reply.Config.Buckets = append(reply.Config.Buckets, int64(b))
				}
				reply.Config.Groups = make(map[int64]string)
				for groupId, addrs := range conf.Groups {
					reply.Config.Groups[int64(groupId)] = strings.Join(addrs, ",")
				}
				reply.LeaderId = cs.tinnRf.GetLeaderId()
			}
			if err != nil {
				reply.ErrMsg = err.Error()
			}

			ch := cs.getNotifyChan(int(appliedMsg.CommandIndex))
			ch <- reply
		}
	}
}

func (cs *ConfigServer) DoConfig(ctx context.Context, args *tinnraftpb.ConfigArgs) (*tinnraftpb.ConfigReply, error) {
	reply := &tinnraftpb.ConfigReply{}

	reqBytes, err := json.Marshal(args)
	if err != nil {
		reply.ErrMsg = err.Error()
		return reply, err
	}

	idx, _, isLeader := cs.tinnRf.Propose(reqBytes)
	if !isLeader {
		reply.ErrMsg = "not Leader node"
		reply.ErrCode = common.ErrCodeWrongLeader
		reply.LeaderId = cs.tinnRf.GetLeaderId()
		return reply, nil
	}

	cs.mu.Lock()
	ch := cs.getNotifyChan(idx)
	cs.mu.Unlock()

	select {
	case res := <-ch:
		reply.Config = res.Config
		reply.ErrMsg = res.ErrMsg
		reply.ErrCode = common.ErrCodeNoErr
		reply.LeaderId = res.LeaderId
	case <-time.After(3 * time.Second):
		reply.ErrMsg = "server timeout"
		reply.ErrCode = common.ErrCodeExecTimeout
		return reply, errors.New("ServerTimerout")
	}

	go func() {
		cs.mu.RLock()
		delete(cs.notifyChans, idx)
		cs.mu.RUnlock()
	}()

	return reply, nil

}
