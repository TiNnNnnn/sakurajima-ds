package shared_server

import (
	"context"
	"encoding/json"
	"errors"
	api_gateway "sakurajima-ds/api_gateway_2"
	"sakurajima-ds/common"
	"sakurajima-ds/config_server"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type ShardKV struct {
	mu        sync.RWMutex
	dead      int32
	tinnrf    *tinnraft.Raft
	applyCh   chan *tinnraftpb.ApplyMsg
	nodeId    int
	groupId   int
	configCli *config_server.ConfigClient

	lastApplied int
	lastConfig  config_server.Config //上一版本的集群配置表
	curConfig   config_server.Config //当前版本的集群配置表

	stm map[int]*Bucket //当前服务器的桶数据

	engine storage_engine.KvStorage

	notifyChans map[int]chan *tinnraftpb.CommandReply

	stopApplyCh chan interface{}

	apiGateClient *api_gateway.ApiGatwayClient

	tinnraftpb.UnimplementedRaftServiceServer
}

func MakeShardKVServer(peerMaps map[int]string, nodeId int, groupId int, configSvrAddrs string) *ShardKV {
	clients := []*tinnraft.ClientEnd{}
	for id, addr := range peerMaps {
		newcli := tinnraft.MakeClientEnd(uint64(id), addr)
		clients = append(clients, newcli)
	}

	newApplyCh := make(chan *tinnraftpb.ApplyMsg)

	logengine := storage_engine.EngineFactory("leveldb",
		"../SALOG/log_data/shared_svr/group_"+strconv.Itoa(groupId)+"/node_"+strconv.Itoa(nodeId))

	apigateclient := api_gateway.MakeApiGatwayClient(99, "127.0.0.1:10030")

	tinnRf := tinnraft.MakeRaft(clients, nodeId, logengine, newApplyCh, apigateclient)

	newengine := storage_engine.EngineFactory("leveldb",
		"../SALOG/shared_data/group_"+strconv.Itoa(groupId)+"/node_"+strconv.Itoa(nodeId))

	shardkv := &ShardKV{
		tinnrf:        tinnRf,
		dead:          0,
		lastApplied:   0,
		applyCh:       newApplyCh,
		nodeId:        nodeId,
		groupId:       groupId,
		lastConfig:    config_server.MakeDefaultConfig(),
		curConfig:     config_server.MakeDefaultConfig(),
		stm:           make(map[int]*Bucket),
		engine:        newengine,
		notifyChans:   map[int]chan *tinnraftpb.CommandReply{},
		configCli:     config_server.MakeCfgSvrClient(99, strings.Split(configSvrAddrs, ",")),
		apiGateClient: apigateclient,
	}

	//初始化buckets
	shardkv.MakeStm(shardkv.engine)

	shardkv.curConfig = *shardkv.configCli.Query(-1)
	shardkv.lastConfig = *shardkv.configCli.Query(-1)

	shardkv.stopApplyCh = make(chan interface{})

	go shardkv.ApplingToStm(shardkv.stopApplyCh)

	go shardkv.ConfigAction()

	return shardkv

}

func (s *ShardKV) IsKilled() bool {
	return atomic.LoadInt32(&s.dead) == 1
}

func (s *ShardKV) GetRaft() *tinnraft.Raft {
	return s.tinnrf
}

func (s *ShardKV) CloseApply() {
	close(s.stopApplyCh)
}

// 初始化状态机,创建BucketsNum个桶
func (s *ShardKV) MakeStm(engine storage_engine.KvStorage) {
	for i := 0; i < common.BucketsNum; i++ {
		if _, ok := s.stm[i]; !ok {
			s.stm[i] = MakeNewBucket(engine, i)
		}
	}
}

func (s *ShardKV) getNotifyChan(index int) chan *tinnraftpb.CommandReply {
	if _, ok := s.notifyChans[index]; !ok {
		s.notifyChans[index] = make(chan *tinnraftpb.CommandReply, 1)
	}
	return s.notifyChans[index]
}

func (s *ShardKV) LegToServe(bucketId int) bool {
	if s.curConfig.Buckets[bucketId] == s.groupId && (s.stm[bucketId].Status == Runing) {
		return true
	}
	return false
}

// 实现 RequestVote interface
func (cs *ShardKV) RequestVote(ctx context.Context, args *tinnraftpb.RequestVoteArgs) (*tinnraftpb.RequestVoteReply, error) {
	reply := &tinnraftpb.RequestVoteReply{}
	cs.tinnrf.HandleRequestVote(args, reply)
	return reply, nil
}

// 实现 AppendEntries interface
func (cs *ShardKV) AppendEntries(ctx context.Context, args *tinnraftpb.AppendEntriesArgs) (*tinnraftpb.AppendEntriesReply, error) {
	reply := &tinnraftpb.AppendEntriesReply{}
	cs.tinnrf.HandleAppendEntries(args, reply)
	return reply, nil
}

// 监听 configserver的配置变化
func (s *ShardKV) ConfigAction() {
	for !s.IsKilled() {
		_, isLeader := s.tinnrf.GetState()
		if isLeader {
			tinnraft.DLog("go into config action")
			s.mu.RLock()
			allowedUpdtaeConf := true
			for _, bucket := range s.stm {
				if bucket.Status != Runing {
					allowedUpdtaeConf = false
					tinnraft.DLog("can't perform next conf")
					break
				}
			}

			// log.Printf("-------------------------------------\n")
			// for i, b := range s.stm {
			// 	log.Printf("{%v,%v,%v} ", i, b.ID, b.Status)
			// }
			// log.Printf("-------------------------------------\n")

			if allowedUpdtaeConf {
				tinnraft.DLog("allowed to perfrom next conf")
			}
			curConfVersion := s.curConfig.Version
			s.mu.RUnlock()

			if allowedUpdtaeConf {
				//向configServer请求最新版本
				nextConfig := s.configCli.Query(int64(curConfVersion) + 1)
				if nextConfig == nil || nextConfig.Groups == nil || len(nextConfig.Buckets) == 0 {
					//LOG
					servicelog := &tinnraftpb.LogArgs{
						Op:       tinnraftpb.LogOp_ListenConfigFailed,
						Contents: "get config from configserver failed",
						FromId:   strconv.Itoa(s.nodeId),
						ToId:     strconv.Itoa(nextConfig.LeaderId),
						CurState: "leader",
						Pid:      int64(syscall.Getpid()),
						Layer:    tinnraftpb.LogLayer_SERVICE,
					}
					s.apiGateClient.SendLogToGate(servicelog)

					time.Sleep(2 * time.Second)
					continue
				}

				//LOG
				servicelog := &tinnraftpb.LogArgs{
					Op:       tinnraftpb.LogOp_ListenConfigSuccess,
					Contents: "get config from configserver success",
					FromId:   strconv.Itoa(s.nodeId),
					ToId:     strconv.Itoa(nextConfig.LeaderId),
					CurState: "leader",
					Pid:      int64(syscall.Getpid()),
					Layer:    tinnraftpb.LogLayer_SERVICE,
				}
				s.apiGateClient.SendLogToGate(servicelog)

				nextConfigBytes, _ := json.Marshal(nextConfig)
				curConfigBytes, _ := json.Marshal(s.curConfig)
				lastConfigBytes, _ := json.Marshal(s.lastConfig)

				tinnraft.DLog("configserver last conf: %v", string(nextConfigBytes))
				tinnraft.DLog("my cur conf: %v", string(curConfigBytes))
				tinnraft.DLog("my last conf: %v", string(lastConfigBytes))

				if nextConfig.Version == curConfVersion+1 {
					//Leader通过Propose向raft层提交一个OpType_ConfigChange的提案
					args := &tinnraftpb.CommandArgs{}
					nextConfigBytes, _ := json.Marshal(nextConfig)

					args.Context = nextConfigBytes
					args.OpType = tinnraftpb.OpType_ConfigChange

					argsBytes, _ := json.Marshal(args)
					idx, _, isLeader := s.tinnrf.Propose(argsBytes)
					if !isLeader {
						return
					}

					s.mu.Lock()
					ch := s.getNotifyChan(idx)
					s.mu.Unlock()

					reply := &tinnraftpb.CommandReply{}

					//监听读取ch(此时已经应用到stm),填写reply返回给client
					select {
					case res := <-ch:
						reply.Value = res.Value
					case <-time.After(3 * time.Second):
					}

					tinnraft.DLog("propose config change success")

					go func() {
						s.mu.Lock()
						delete(s.notifyChans, idx)
						s.mu.Unlock()
					}()
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}

// 读取applyCh，持久化信息到状态机
func (s *ShardKV) ApplingToStm(done <-chan interface{}) {
	for !s.IsKilled() {
		select {
		case <-done:
			return
		case appliedMsg := <-s.applyCh:
			tinnraft.DLog("appling msg: %v", appliedMsg.String())
			args := &tinnraftpb.CommandArgs{}

			uerr := json.Unmarshal(appliedMsg.Command, args)
			if uerr != nil {
				tinnraft.DLog("Unmarshal CommandArgs failed")
				continue
			}
			s.lastApplied = int(appliedMsg.CommandIndex)

			reply := &tinnraftpb.CommandReply{}
			value := ""
			var err error

			switch args.OpType {
			case tinnraftpb.OpType_Put: //上传数据
				bucketId := common.KeyToBucketId(args.Key)
				if s.LegToServe(bucketId) {
					tinnraft.DLog("put " + args.Key + " value " + args.Value + "to bucket " + strconv.Itoa(bucketId))
					err = s.stm[bucketId].Put(args.Key, args.Value)
					if err != nil {
						tinnraft.DLog("applingtostm err: %v", err.Error())
					}
					// LOG
					raftlog := &tinnraftpb.LogArgs{
						Op:       tinnraftpb.LogOp_PutKv,
						Contents: "put the kv to the stm success",
						FromId:   strconv.Itoa(s.nodeId),
						Pid:      int64(syscall.Getpid()),
						BucketId: strconv.Itoa(bucketId),
						GroupId:  strconv.Itoa(s.groupId),
						Layer:    tinnraftpb.LogLayer_PERSIST,
					}
					s.apiGateClient.SendLogToGate(raftlog)

				}
			case tinnraftpb.OpType_Get: //下载数据
				bucketId := common.KeyToBucketId(args.Key)
				if s.LegToServe(bucketId) {
					value, err = s.stm[bucketId].Get(args.Key)
					if err != nil {
						tinnraft.DLog("applingtostm err: %v", err.Error())
					}
					// LOG
					raftlog := &tinnraftpb.LogArgs{
						Op:       tinnraftpb.LogOp_GetKv,
						Contents: "get the key from the stm success",
						FromId:   strconv.Itoa(s.nodeId),
						Pid:      int64(syscall.Getpid()),
						BucketId: strconv.Itoa(bucketId),
						GroupId:  strconv.Itoa(s.groupId),
						Layer:    tinnraftpb.LogLayer_PERSIST,
					}
					s.apiGateClient.SendLogToGate(raftlog)
					tinnraft.DLog("get " + args.Key + " value " + value + " from bucket " + strconv.Itoa(bucketId))
				}
				reply.Value = value
			case tinnraftpb.OpType_ConfigChange: //ConfigServer内配置发生变化，SharedServer Config对应变化
				nextConfig := &config_server.Config{}
				json.Unmarshal(args.Context, nextConfig)

				if nextConfig.Version == s.curConfig.Version+1 {
					for i := 0; i < common.BucketsNum; i++ {
						//更新bucketidx和groupId的映射关系
						if s.curConfig.Buckets[i] != s.groupId && nextConfig.Buckets[i] == s.groupId {
							groupId := s.curConfig.Buckets[i]
							if groupId != 0 {
								s.stm[i].Status = Runing //启动该桶
								tinnraft.DLog("chang the bucket %d status to RUNNING", i)
							}
						}
						if s.curConfig.Buckets[i] == s.groupId && nextConfig.Buckets[i] != s.groupId {
							groupId := nextConfig.Buckets[i]
							if groupId != 0 {
								s.stm[i].Status = Stopped //关闭该桶
								tinnraft.DLog("chang the bucket %d status to STOPPED", i)
							}
						}
					}
				}

				s.lastConfig = s.curConfig
				s.curConfig = *nextConfig

				configBytes, _ := json.Marshal(s.curConfig)
				tinnraft.DLog("appiled config to server: %v", string(configBytes))

			case tinnraftpb.OpType_InsertBuckets:
				bucketOpArgs := &tinnraftpb.BucketOpArgs{}
				json.Unmarshal(args.Context, bucketOpArgs)
				bucketdatas := map[int]map[string]string{}
				json.Unmarshal(bucketOpArgs.BucketsData, &bucketdatas)
				for bucket_id, kvs := range bucketdatas {
					s.stm[bucket_id] = MakeNewBucket(s.engine, bucket_id)
					for k, v := range kvs {
						s.stm[bucket_id].Put(k, v)
						tinnraft.DLog("insert kv data to bucket,k: %v v: %v", k, v)
					}
				}
			case tinnraftpb.OpType_DeleteBuckets:
				bucketOpArgs := &tinnraftpb.BucketOpArgs{}
				json.Unmarshal(args.Context, bucketOpArgs)
				for _, bid := range bucketOpArgs.BucketIds {
					s.stm[int(bid)].DelBucketData()
					tinnraft.DLog("delete buckets data list" + strconv.Itoa(int(bid)))
				}
			}

			if err != nil {
				tinnraft.DLog("applingtostm err: %v", err.Error())
			}

			ch := s.getNotifyChan(int(appliedMsg.CommandIndex))
			ch <- reply
		}
	}
}

// 重写DoCommand方法 （被客户端直接调用）
func (s *ShardKV) DoCommand(ctx context.Context, args *tinnraftpb.CommandArgs) (*tinnraftpb.CommandReply, error) {

	reply := &tinnraftpb.CommandReply{}

	argsBytes, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	if !s.LegToServe(common.KeyToBucketId(args.Key)) {
		reply.ErrCode = common.ErrCodeWrongGroup
		return reply, nil
	}

	idx, _, isLeader := s.tinnrf.Propose(argsBytes)
	if !isLeader {
		reply.ErrCode = common.ErrCodeWrongLeader
		reply.LeaderId = s.tinnrf.GetLeaderId()
		return reply, nil
	}

	s.mu.Lock()
	ch := s.getNotifyChan(idx)
	s.mu.Unlock()

	select {
	case res := <-ch:
		if res != nil {
			reply.ErrCode = common.ErrCodeNoErr
			reply.Value = res.Value
		}
	case <-time.After(3 * time.Second):
		return reply, errors.New("exec has timeout")
	}

	go func() {
		s.mu.Lock()
		delete(s.notifyChans, idx)
		s.mu.Unlock()
	}()

	return reply, nil
}

// 重写DoBucketRpc方法（被客户端直接调用）
func (s *ShardKV) DoBucket(ctx context.Context, args *tinnraftpb.BucketOpArgs) (*tinnraftpb.BucketOpReply, error) {
	
	reply := &tinnraftpb.BucketOpReply{}
	if _, isLeader := s.tinnrf.GetState(); !isLeader {
		return reply, errors.New("ErrorWrongLeader")
	}

	switch args.BucketOpType {
	case tinnraftpb.BucketOpType_GetData: //获取bucket数据
		{
			s.mu.RLock()
			if s.curConfig.Version < int(args.ConfigVersion) {
				s.mu.RUnlock()
				return reply, errors.New("ErrorNotReady")
			}

			bucketdatas := map[int]map[string]string{}
			for _, bucketId := range args.BucketIds {
				datas, err := s.stm[int(bucketId)].DeepCopy() //DeepCopy
				if err != nil {
					s.mu.RUnlock()
					return reply, err
				}
				bucketdatas[int(bucketId)] = datas
			}
			bucketdataBytes, _ := json.Marshal(bucketdatas)
			reply.BucketData = bucketdataBytes
			reply.ConfigVersion = args.ConfigVersion
			s.mu.RUnlock()
		}
	case tinnraftpb.BucketOpType_AddData:
		{
			s.mu.RLock()
			if int64(s.curConfig.Version) > args.ConfigVersion {
				s.mu.RUnlock()
				return reply, nil
			}
			s.mu.RUnlock()
			comandArgs := &tinnraftpb.CommandArgs{}
			bucketOpArgsBytes, _ := json.Marshal(args)
			comandArgs.Context = bucketOpArgsBytes
			comandArgs.OpType = tinnraftpb.OpType_InsertBuckets

			comandArgsBytes, _ := json.Marshal(comandArgs)

			_, _, isLeader := s.tinnrf.Propose(comandArgsBytes)
			if !isLeader {
				return reply, nil
			}
		}
	case tinnraftpb.BucketOpType_DelData:
		{
			s.mu.RLock()
			if int64(s.curConfig.Version) > args.ConfigVersion {
				s.mu.RUnlock()
				return reply, nil
			}
			s.mu.RUnlock()
			commandArgs := &tinnraftpb.CommandArgs{}
			bucketOPArgsBytes, _ := json.Marshal(args)
			commandArgs.Context = bucketOPArgsBytes
			commandArgs.OpType = tinnraftpb.OpType_DeleteBuckets

			commandArgsBytes, _ := json.Marshal(commandArgs)

			_, _, isLeader := s.tinnrf.Propose(commandArgsBytes)
			if !isLeader {
				return reply, nil
			}
		}
	}
	return reply, nil
}
