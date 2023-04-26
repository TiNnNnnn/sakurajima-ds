package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sakurajima-ds/config_server"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"strings"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("usage: server [nodeId] [configserveraddr1,configserveraddr2,configserveraddr3]")
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	nodeIdStr := os.Args[1]
	nodeId, err := strconv.Atoi(nodeIdStr)
	if err != nil {
		panic(err)
	}

	configSvrAddrs := strings.Split(os.Args[2], ",")
	cfPeerMap := make(map[int]string)
	for i, addr := range configSvrAddrs {
		cfPeerMap[i] = addr
	}

	cfgSvr := config_server.MakeConfigServer(cfPeerMap, nodeId)
	lis, err := net.Listen("tcp", cfPeerMap[nodeId])
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	tinnraft.DLog("server listen on:  " + cfPeerMap[nodeId])
	s := grpc.NewServer()

	tinnraftpb.RegisterRaftServiceServer(s, cfgSvr)

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		cfgSvr.GetRaft().CloseAllConn()
		cfgSvr.StopApply()
		os.Exit(-1)
	}()

	reflection.Register(s)
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
