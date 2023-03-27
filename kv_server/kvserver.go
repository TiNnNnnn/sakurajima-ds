package kv_server

import (
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"sync"
)

/*
kv存储服务器
*/

type KvServer struct {
	mu     sync.Mutex
	dead   int32
	tinnRf *tinnraft.Raft
	apply  chan *tinnraftpb.ApplyMsg

	lastApplied    int
	stm            int
	lastOperations map[int64]int
}

func MakeKvServer(serverId int) *KvServer {
	
}
