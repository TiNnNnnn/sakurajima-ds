package main

import (
	"context"
	"fmt"
	"log"
	"os"
	pb "sakurajima-ds/test/study/rpc_test1/lab1/genpb"

	"google.golang.org/grpc"
)

const (
	address = "localhost:10051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKvServiceClient(conn)

	op := ""
	key := ""
	value := ""

	op = os.Args[1]
	if op == "put" {
		key = os.Args[2]
		value = os.Args[3]
		args := &pb.PutArgs{
			Key:   key,
			Value: value,
		}
		_, err := c.DoPut(context.Background(), args)
		if err != nil {
			log.Printf("err: %v", err)
			return
		}
		fmt.Printf("success")

	} else if op == "get" {
		key = os.Args[2]
		args := &pb.GetArgs{
			Key: key,
		}
		r, err := c.DoGet(context.Background(), args)
		if err != nil {
			log.Printf("err: %v", err)
			return
		}
		fmt.Printf(r.Value)
	}
}
