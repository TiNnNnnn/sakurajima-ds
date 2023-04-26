package config_server

import (
	"context"
	"crypto/rand"
	"math/big"
	"sakurajima-ds/common"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"strings"
)

type ConfigClient struct {
	//mu         sync.Mutex
	clientends []*tinnraft.ClientEnd
	leaderId   int64
	clientId   int64
	commandId  int64
}

func MakeCfgSvrClient(targetId uint64, targetAddrs []string) *ConfigClient {

	cfgCli := &ConfigClient{
		leaderId:  0,
		clientId:  nrand(),
		commandId: 0,
	}

	for _, addr := range targetAddrs {
		newcli := tinnraft.MakeClientEnd(targetId, addr)
		cfgCli.clientends = append(cfgCli.clientends, newcli)
	}

	return cfgCli
}

func (cfg *ConfigClient) GetRpcClients() []*tinnraft.ClientEnd {
	return cfg.clientends
}

func (cfgCli *ConfigClient) Query(version int64) *Config {
	confArgs := &tinnraftpb.ConfigArgs{
		OpType:        tinnraftpb.ConfigOpType_Query,
		ConfigVersion: version,
	}
	reply := cfgCli.CallDoConfig(confArgs)
	config := &Config{}
	if reply != nil && reply.Config != nil {
		config.Version = int(reply.Config.ConfigVersion)
		for i := 0; i < common.BucketsNum; i++ {
			config.Buckets[i] = int(reply.Config.Buckets[i])
		}
		config.Groups = make(map[int][]string)
		for k, v := range reply.Config.Groups {
			slist := strings.Split(v, ",")
			config.Groups[int(k)] = slist
		}
	} else {
		tinnraft.DLog("query config reply failed , all configserver down")
	}
	return config
}

func (cfgCli *ConfigClient) Join(servers map[int64]string) bool {
	confArgs := &tinnraftpb.ConfigArgs{
		OpType:  tinnraftpb.ConfigOpType_Join,
		Servers: servers,
	}
	confReply := cfgCli.CallDoConfig(confArgs)
	if confReply == nil {
		tinnraft.DLog("join config reply failed , all configserver down")
		return false
	}
	tinnraft.DLog("join config reply ok")
	return true
}

func (cfgCli *ConfigClient) Move(bucketId int, groupId int) bool {
	confArgs := &tinnraftpb.ConfigArgs{
		OpType:   tinnraftpb.ConfigOpType_Move,
		BucketId: int64(bucketId),
		GroupId:  int64(groupId),
	}
	confReply := cfgCli.CallDoConfig(confArgs)
	if confReply == nil {
		tinnraft.DLog("move bucket reply failed,move bid %v to group %v failed , all configserver down", bucketId, groupId)
		return false
	}
	tinnraft.DLog("move bucket reply ok,move bid %v to group %v", bucketId, groupId)
	return true
}

func (cfgCli *ConfigClient) Leave(gids []int64) {
	confReq := &tinnraftpb.ConfigArgs{
		OpType:   tinnraftpb.ConfigOpType_Leave,
		GroupIds: gids,
	}
	resp := cfgCli.CallDoConfig(confReq)
	if resp == nil {
		tinnraft.DLog("leave config reply failed , all configserver down")
		return
	}
	tinnraft.DLog("leave cfg resp ok" + resp.ErrMsg)
}

func (cfg *ConfigClient) CallDoConfig(args *tinnraftpb.ConfigArgs) *tinnraftpb.ConfigReply {
	var err error
	confReply := &tinnraftpb.ConfigReply{}
	confReply.Config = &tinnraftpb.ServerConfig{}

	for _, clientend := range cfg.clientends {
		confReply, err = (*clientend.GetRaftServiceCli()).DoConfig(context.Background(), args)
		if err != nil {
			tinnraft.DLog("a node in config cluster is down,try next,err:%v\n", err.Error())
			continue
		}
		switch confReply.ErrCode {
		case common.ErrCodeNoErr:
			cfg.commandId++
			return confReply
		case common.ErrCodeWrongLeader:
			tinnraft.DLog("find leader id is %d", confReply.LeaderId)
			confReply, err := (*cfg.clientends[confReply.LeaderId].GetRaftServiceCli()).DoConfig(context.Background(), args)
			if err != nil {
				tinnraft.DLog("leader in config cluster is down")
				continue
			}
			if confReply.ErrCode == common.ErrCodeNoErr {
				cfg.commandId++
				return confReply
			}
			if confReply.ErrCode == common.ErrCodeExecTimeout {
				tinnraft.DLog("exec has timeout")
				return confReply
			}
			return confReply
		}
	}
	return confReply
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	return bigx.Int64()
}
