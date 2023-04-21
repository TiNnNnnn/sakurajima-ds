package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"sakurajima-ds/common"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	_"time"
)

type KvClient struct {
	Client    *tinnraft.ClientEnd
	LeadId    int64
	CLientId  int64
	CommandId int64
}

func MakeKvClient(serverId int, serverAddr string) *KvClient {
	client := tinnraft.MakeClientEnd(uint64(serverId), serverAddr)
	return &KvClient{
		Client:    client,
		LeadId:    0,
		CLientId:  getrandInt(),
		CommandId: 0,
	}
}



func (kvc *KvClient) Get(key string) string {
	args := &tinnraftpb.CommandArgs{
		Key:      key,
		OpType:   tinnraftpb.OpType_Get,
		ClientId: kvc.CLientId,
	}
	reply, err := (*kvc.Client.GetRaftServiceCli()).DoCommand(context.Background(), args)
	if err != nil {
		return "err"
	}
	return reply.Value
}

func (kvc *KvClient) Put(key string, value string) string {
	args := &tinnraftpb.CommandArgs{
		Key:      key,
		Value:    value,
		ClientId: kvc.CLientId,
		OpType:   tinnraftpb.OpType_Put,
	}
	_, err := (*kvc.Client.GetRaftServiceCli()).DoCommand(context.Background(), args)
	if err != nil {
		return "err"
	}
	return "ok"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: kvclient [serverAddr] [count] [op] [options]")
		return
	}
	sigs := make(chan os.Signal, 1)

	kvclient := MakeKvClient(99, os.Args[1])

	count, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	op := os.Args[3]

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		kvclient.Close()
		os.Exit(-1)
	}()

	keys := make([]string, count)
	vals := make([]string, count)

	for i := 0; i < count; i++ {
		randK := common.RandStringRunes(8)
		randV := common.RandStringRunes(8)

		keys[i] = randK
		vals[i] = randV
	}

	// startTs := time.Now()
	// for i := 0; i < count; i++ {
	// 	kvclient.Put(keys[i], vals[i])
	// }
	// elapsed := time.Since(startTs).Seconds()
	// fmt.Printf("total cost %f s\n", elapsed)

	switch op {
	case "get":
		fmt.Println(kvclient.Get(os.Args[4]))
	case "put":
		fmt.Println(kvclient.Put(os.Args[4], os.Args[5]))
	}

}

func (kvc *KvClient) Close() {
	kvc.Client.CloseConns()
}

func getrandInt() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	return bigx.Int64()
}
