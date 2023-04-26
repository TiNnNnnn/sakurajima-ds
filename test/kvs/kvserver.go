package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	ksv "sakurajima-ds/kv_server"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("kv_server usageL: server [serverId]")
		return
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	serverId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	kvserver := ksv.MakeKvServer(serverId)
	listen, err := net.Listen("tcp", ksv.PeersMap[serverId])
	if err != nil {
		fmt.Printf("faield to listen: %v\n", err)
		return
	}
	
	fmt.Printf("[server %v] listening on %s\n", serverId, ksv.PeersMap[serverId])
	s := grpc.NewServer()
	tinnraftpb.RegisterRaftServiceServer(s, kvserver)

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		kvserver.GetTinnRaft().CloseAllConn()
		os.Exit(-1)
	}()

	reflection.Register(s)
	err = s.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}

}
