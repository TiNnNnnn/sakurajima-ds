package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"strings"
	"syscall"

	shared_server "sakurajima-ds/shared_server_raft"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("usage: server [nodeId] [gId] [csAddr] [server1addr,server2addr,server3addr]")
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	nodeIdStr := os.Args[1]
	nodeId, err := strconv.Atoi(nodeIdStr)
	if err != nil {
		panic(err)
	}

	gIdStr := os.Args[2]
	gId, err := strconv.Atoi(gIdStr)
	if err != nil {
		panic(err)
	}

	svrAddrs := strings.Split(os.Args[4], ",")
	svrPeerMap := make(map[int]string)
	for i, addr := range svrAddrs {
		svrPeerMap[i] = addr
	}

	shardSvr := shared_server.MakeShardKVServer(svrPeerMap, nodeId, gId, os.Args[3])
	lis, err := net.Listen("tcp", svrPeerMap[nodeId])
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	fmt.Printf("server listen on: %s \n", svrPeerMap[nodeId])
	s := grpc.NewServer()
	tinnraftpb.RegisterRaftServiceServer(s, shardSvr)

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		shardSvr.GetRaft().CloseAllConn()
		shardSvr.CloseApply()
		os.Exit(-1)
	}()

	reflection.Register(s)
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}

}
