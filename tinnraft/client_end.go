package tinnraft

import (
	"google.golang.org/grpc"
)

type ClientEnd struct {
	id    uint64
	addr  string
	conns []*grpc.ClientConn
}


