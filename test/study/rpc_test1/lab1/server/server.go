package main

import (
	"context"
	"errors"
	"log"
	"net"

	pb "sakurajima-ds/test/study/rpc_test1/lab1/genpb"

	"google.golang.org/grpc"
)

const (
	port = ":10051"
)


type server struct {
	stm map[string]string
	pb.UnimplementedKvServiceServer
}

/*
实例化一个server
*/
func MakeServer() *server {
	newSvr := server{
		stm: map[string]string{},
	}
	return &newSvr
}

/*
上传函数
*/
func (s *server) DoPut(ctx context.Context, args *pb.PutArgs) (*pb.PutReply, error) {

	reply := pb.PutReply{}
	key, value := args.Key, args.Value

	if key == "" || value == "" {
		reply.Success = "false"
		reply.ErrMsg = "empty key or value"
		reply.ErrCode = 3
		return &reply, errors.New(reply.ErrMsg)
	}

	log.Printf("recv the key: %v value: %v", key, value)
	s.stm[args.Key] = args.Value

	return &reply, nil

}

/*
启动server
*/
func startServer() {

	newSvr := MakeServer()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterKvServiceServer(s, newSvr)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	startServer()
}
