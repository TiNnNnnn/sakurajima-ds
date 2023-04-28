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
	stm      map[string]string
	pb.UnimplementedKvServiceServer
}

/*
实例化一个server
*/
func MakeServer() *server {
	newSvr := server{
		stm:      map[string]string{},
	}
	return &newSvr
}

/*
上传函数
*/
func (s *server) DoPut(ctx context.Context, args *pb.PutArgs) (*pb.PutReply, error) {

	//user write
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

	reply.Success = "true"

	return &reply, nil

}

func (s *server) DoGet(ctx context.Context, args *pb.GetArgs) (*pb.GetReply, error) {
	/*
		user write
	*/
	reply := pb.GetReply{}
	key := args.Key

	value := s.stm[key]
	log.Printf("get the value:%v ", value)

	reply.Value = value

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
