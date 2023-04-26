/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"errors"
	"log"
	"net"

	pb "sakurajima-ds/test/study/rpc_test1/genpb"

	"google.golang.org/grpc"
)

const (
	port = ":10051"
)

type server struct {
	stm map[string]string
	pb.UnimplementedKvServiceServer
}

func MakeServer() *server {
	newSvr := server{
		stm: map[string]string{},
	}
	return &newSvr
}

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
