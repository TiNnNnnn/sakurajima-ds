package meta_server

/*
message MetaArgs{
	MetaOpType op = 1;
	int64 UserId = 2;
	string BucketName = 3;
	string ObjectName = 4;
	int64 Version = 5;
	int64 Size = 6;
	string Hash = 7;
}


message Metadata{
	int64 UserId = 2;
	string BucketName = 3;
	string ObjectName = 4;
	int64 Version = 5;
	int64 Size = 6;
	string Hash = 7;
}

message MetaReply{
	string ErrMsg = 1;
	int64 ErrCode = 2;
	int64 LeaderId = 3;
	repeated Metadata metas = 4;
	Metadata meta = 5;
}
*/

import (
	"context"
	"crypto/rand"
	"math/big"
	"sakurajima-ds/common"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
)

type MetaClient struct {
	//mu        sync.Mutex
	Clients   []*tinnraft.ClientEnd
	LeaderId  int64
	ClientId  int64
	CommandId int64
}

func MakeMetaSvrClient(targetId uint64, targetAddrs []string) *MetaClient {
	metaCli := &MetaClient{
		LeaderId:  0,
		ClientId:  nrand(),
		CommandId: 0,
	}

	for _, addr := range targetAddrs {
		cli := tinnraft.MakeClientEnd(targetId, addr)
		metaCli.Clients = append(metaCli.Clients, cli)
	}

	return metaCli
}

func (MetaCli *MetaClient) GetRpcClis() []*tinnraft.ClientEnd {
	return MetaCli.Clients
}

func (mc *MetaClient) GetMetadata(userId int64, BucketName string, ObjectName string, version int) *MetaData {
	metaArgs := &tinnraftpb.MetaArgs{
		Op:         tinnraftpb.MetaOpType_GetMetadata,
		UserId:     userId,
		BucketName: BucketName,
		ObjectName: ObjectName,
		Version:    int64(version),
	}
	reply := mc.CallDoMetaRpc(metaArgs)
	retMeta := &MetaData{}
	if reply != nil {
		retMeta.Hash = reply.Meta.Hash
		retMeta.Size = reply.Meta.Size
	}
	tinnraft.DLog("get the metadata success")
	return retMeta
}

func (mc *MetaClient) PutMetadata(userId int64, bucketName string, objectName string, version int, size int64, hash string) {
	metaArgs := &tinnraftpb.MetaArgs{
		Op:         tinnraftpb.MetaOpType_PutMetadata,
		UserId:     userId,
		BucketName: bucketName,
		ObjectName: objectName,
		Version:    int64(version),
		Size:       size,
		Hash:       hash,
	}
	reply := mc.CallDoMetaRpc(metaArgs)
	if reply.ErrMsg != "" {
		tinnraft.DLog("put the metadata failed: %v", reply.ErrMsg)
		return
	}
	tinnraft.DLog("put the metadata success")
}



func (MetaCli *MetaClient) CallDoMetaRpc(args *tinnraftpb.MetaArgs) *tinnraftpb.MetaReply {
	var err error
	metaReply := &tinnraftpb.MetaReply{}

	for _, end := range MetaCli.Clients {
		metaReply, err = (*end.GetRaftServiceCli()).DoMeta(context.Background(), args)
		if err != nil {
			tinnraft.DLog("a node int cluster is down")
			continue
		}
		switch metaReply.ErrCode { //返回无错误
		case common.ErrCodeNoErr:
			MetaCli.CommandId++
			return metaReply
		case common.ErrCodeWrongLeader: //发送信息到非Leader节点
			tinnraft.DLog("find leader id is %v", metaReply.LeaderId)
			//获取leader节点id,重新发送
			metaResp, err := (*MetaCli.Clients[metaReply.LeaderId].GetRaftServiceCli()).DoMeta(context.Background(), args)
			if err != nil {
				tinnraft.DLog("one node in cluster has down: %v", err.Error())
				continue
			}
			if metaResp.ErrCode == common.ErrCodeNoErr {
				MetaCli.CommandId++
				return metaResp
			}
			if metaResp.ErrCode == common.ErrCodeExecTimeout {
				tinnraft.DLog("exec has timeout")
			}
			return metaReply
		}
	}
	return metaReply
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	return bigx.Int64()
}
