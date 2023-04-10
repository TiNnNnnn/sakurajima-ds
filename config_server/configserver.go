package config_server

import (
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"sync"
)

type ConfigServer struct {
	mu          sync.RWMutex
	dead        int32
	tinnRf      *tinnraft.Raft
	applyCh     chan *tinnraftpb.ApplyMsg
	stm         ConfigStateMachine
	notifyChans map[int]chan *tinnraftpb.CommandReply
	stopApplyCh chan interface{}

	tinnraftpb.UnimplementedRaftServiceServer
}

func MakeConfigServer(peerMaps map[int]string, serverId int) *ConfigServer {
	clientEnds := []*tinnraft.ClientEnd{}

	//创建三个rpc客户端,并添加到clientEnds数组
	for id, addr := range peerMaps {
		newClient := tinnraft.MakeClientEnd(uint64(id), addr)
		clientEnds = append(clientEnds, newClient)
	}

	applyCh := make(chan *tinnraftpb.ApplyMsg)

	Configengine := storage_engine.EngineFactory("leveldb", "./conf_data"+"/node_"+strconv.Itoa(serverId))

	logEngine := storage_engine.EngineFactory("leveldb", "./log_data/"+"configserver/node_"+strconv.Itoa(serverId))

	tinnRf := tinnraft.MakeRaft(clientEnds, serverId, logEngine, applyCh)

	configServer := &ConfigServer{
		dead:        0,
		applyCh:     applyCh,
		notifyChans: make(map[int]chan *tinnraftpb.ConfigReply),
		tinnRf:      tinnRf,
	}

	configServer.stopApplyCh = make(chan interface{})

	return configServer
}


