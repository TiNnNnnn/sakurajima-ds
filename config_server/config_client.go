package config_server

import (
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"sync"
)

type ConfigClient struct {
	mu         sync.Mutex
	clientends []*tinnraft.ClientEnd
	leaderId   int64
	clientId   int64
	commandId  int64
}

func (cfc *ConfigClient) CallDoConfig(args *tinnraftpb.ConfigArgs) *tinnraftpb.ConfigReply {
	return nil
}
