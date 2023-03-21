package tinnraft

import (
	"sync"
	"sync/atomic"
)

type Raft struct {
	mu    sync.Mutex
	peers []*RaftClientEnd
	
}
