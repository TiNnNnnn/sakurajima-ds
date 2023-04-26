package main

import (
	"context"
	"log"
	"os"
	pb "sakurajima-ds/test/study/rpc_test1/genpb"

	"google.golang.org/grpc"
)

const (
	address = "localhost:10051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKvServiceClient(conn)

	key := ""
	value := ""

	if len(os.Args) > 1 {
		key = os.Args[1]
	}
	if len(os.Args) > 2 {
		value = os.Args[2]
	}

	args := &pb.PutArgs{
		Key:   key,
		Value: value,
	}

	r, err := c.DoPut(context.Background(), args)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	log.Printf("success: %v", r.Success)
}
