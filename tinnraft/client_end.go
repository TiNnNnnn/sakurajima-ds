package tinnraft

import (
	"fmt"
	"google.golang.org/grpc"
	"sakurajima-ds/tinnraftpb"
)

// raft节点客户端，抽象化为ClientEnd结构
type ClientEnd struct {
	id             uint64                        //节点编号
	addr           string                        //节点地址
	conns          []*grpc.ClientConn            //客户端连接句柄
	raftServiceCli *tinnraftpb.RaftServiceClient //rpc客户端(用于函数调用)
}

// 初始化ClientEnd
func MakeClientEnd(id uint64, addr string) *ClientEnd {
	//创建客户端连接
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to connect rpc: %v", err)
	}
	//将客户端连接句柄加入conns数组
	conns := []*grpc.ClientConn{}
	conns = append(conns, conn)

	//创建远程连接的客户端，客户端内封装了可用函数
	serviceClient := tinnraftpb.NewRaftServiceClient(conn)

	return &ClientEnd{
		id:             id,
		addr:           addr,
		conns:          conns,
		raftServiceCli: &serviceClient,
	}
}

// 返回节点编号
func (ce *ClientEnd) Id() uint64 {
	return ce.id
}

// 获取rpc客户端
func (ce *ClientEnd) GetRaftServiceCli() *tinnraftpb.RaftServiceClient {
	return ce.raftServiceCli
}

