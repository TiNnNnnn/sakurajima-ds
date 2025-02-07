package api_gateway

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"sakurajima-ds/tinnraftpb"
	"time"
)

type ApiGatwayClient struct {
	//mu         sync.Mutex
	clientend ClientEnd
	clientId  int64
	commandId int64
}

func MakeApiGatwayClient(targetId uint64, targetAddrs string) *ApiGatwayClient {

	apiCli := &ApiGatwayClient{
		clientId:  nrand(),
		commandId: 0,
	}

	newcli := MakeClientEnd(targetId, targetAddrs)
	apiCli.clientend = *newcli

	return apiCli
}

// raft层日志
func (ac *ApiGatwayClient) SendLogToGate(args *tinnraftpb.LogArgs) {

	now := time.Now()
	nowTime := now.UnixNano() / 1e6

	args.Time = nowTime

	reply := ac.CallDoLog(args)
	if reply == nil {
		return
	}
	if !reply.Success {
		log.Println("call dolog Rpc failed")
		return
	}
	log.Println("call dolog Rpc success")
}


// 将log发送给api_server
func (ac *ApiGatwayClient) CallDoLog(args *tinnraftpb.LogArgs) *tinnraftpb.LogReply {
	logReply, err := (*ac.clientend.GetRaftServiceCli()).DoLog(context.Background(), args)
	if err != nil {
		log.Println("api server has down!")
	}
	return logReply
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	return bigx.Int64()
}
