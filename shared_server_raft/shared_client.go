package shared_server

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sakurajima-ds/common"
	"sakurajima-ds/config_server"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"strings"
)

type ShardKVClient struct {
	client       *tinnraft.ClientEnd
	connectsList map[string]*tinnraft.ClientEnd

	configCli *config_server.ConfigClient
	config    *config_server.Config

	leaderId  int64
	clientId  int64
	commandId int64
}

// 初始化一个sharedserver
func MakeSharedKvClient(csAddrs string) *ShardKVClient {
	cfgSvrCli := config_server.MakeCfgSvrClient(999, strings.Split(csAddrs, ","))
	kvCli := &ShardKVClient{
		configCli: cfgSvrCli,
		client:    nil,
		leaderId:  0,
		clientId:  nrand(),
		commandId: 0,
	}
	//缓存最新的configserver 的分组信息
	kvCli.config = kvCli.configCli.Query(-1)
	return kvCli
}

func (kvCli *ShardKVClient) Get(key string) (string, error) {
	return kvCli.SendRpcCommand(&tinnraftpb.CommandArgs{
		Key:    key,
		OpType: tinnraftpb.OpType_Get,
	})
}

func (kvCli *ShardKVClient) Put(key string, value string) error {
	_, err := kvCli.SendRpcCommand(&tinnraftpb.CommandArgs{
		Key:    key,
		Value:  value,
		OpType: tinnraftpb.OpType_Put,
	})
	return err
}

func (kvCli *ShardKVClient) GetBucketDatas(groupId int, bucketIds []int64) string {
	return kvCli.SendBucketRpcCommand(&tinnraftpb.BucketOpArgs{
		BucketOpType:  tinnraftpb.BucketOpType_GetData,
		GroupId:       int64(groupId),
		ConfigVersion: int64(kvCli.config.Version),
		BucketIds:     bucketIds,
	})
}

func (kvCli *ShardKVClient) InsertBucketDatas(gid int, bucketIds []int64, datas []byte) string {
	return kvCli.SendBucketRpcCommand(&tinnraftpb.BucketOpArgs{
		BucketOpType:  tinnraftpb.BucketOpType_AddData,
		BucketsData:   datas,
		GroupId:       int64(gid),
		BucketIds:     bucketIds,
		ConfigVersion: int64(kvCli.config.Version),
	})
}

func (kvCli *ShardKVClient) DeleteBucketDatas(gid int, bucketIds []int64) string {
	return kvCli.SendBucketRpcCommand(&tinnraftpb.BucketOpArgs{
		BucketOpType:  tinnraftpb.BucketOpType_DelData,
		GroupId:       int64(gid),
		ConfigVersion: int64(kvCli.config.Version),
		BucketIds:     bucketIds,
	})
}

func (kvCli *ShardKVClient) SendRpcCommand(args *tinnraftpb.CommandArgs) (string, error) {
	bucket_id := common.KeyToBucketId(args.Key)
	log.Printf("caculated bucketId: %v\n", bucket_id)

	//根据bucketId 查询分组
	groupId := kvCli.config.Buckets[bucket_id]
	if groupId == 0 { //当前bucket没有挂载对应group
		tinnraft.DLog("there is no shared in charge of this bucket")
		return "", errors.New("there is no shared in charge of this bucket")
	}

	//向shared cluster 发送 put request
	if servers, ok := kvCli.config.Groups[groupId]; ok {

		for _, addrs := range servers {
			if kvCli.ReadConnsList(addrs) == nil {
				kvCli.client = tinnraft.MakeClientEnd(99, addrs)
			} else {
				kvCli.client = kvCli.ReadConnsList(addrs)
			}

			reply, err := (*kvCli.client.GetRaftServiceCli()).DoCommand(context.Background(), args)
			if err != nil {
				tinnraft.DLog("a node in config cluster has down,begin try a another one")
				continue
			}

			switch reply.ErrCode {
			case common.ErrCodeNoErr:
				kvCli.commandId++
				return reply.Value, nil
			case common.ErrCodeWrongGroup:
				kvCli.config = kvCli.configCli.Query(-1)
				return "", errors.New("find Wrong Group")
			case common.ErrCodeWrongLeader:
				kvCli.client = tinnraft.MakeClientEnd(99, servers[reply.LeaderId])
				reply, err := (*kvCli.client.GetRaftServiceCli()).DoCommand(context.Background(), args)
				if err != nil {
					fmt.Printf("the leader in config cluster has down,err: %s", err.Error())
					panic(err)
				}
				if reply.ErrCode == common.ErrCodeNoErr {
					kvCli.commandId++
					return reply.Value, nil
				}
			default:
				return "", errors.New("unkown code")
			}
		}
	} else {
		return "", errors.New("please join the server group first")
	}
	return "", errors.New("Unknow error")
}

// func (kvCli *ShardKVClient) SendBucketRpcCommand(args *tinnraftpb.BucketOpArgs) string {
// 	for {
// 		if servers, ok := kvCli.config.Groups[int(args.GroupId)]; ok {
// 			for _, serverAddr := range servers { //获取当前分组对应的集群server addr，尝试发送
// 				kvCli.client = tinnraft.MakeClientEnd(99, serverAddr)
// 				reply, err := (*kvCli.client.GetRaftServiceCli()).DoBucket(context.Background(), args)
// 				if err == nil {
// 					if reply != nil {
// 						return string(reply.BucketData)
// 					} else {
// 						return ""
// 					}
// 				} else {
// 					tinnraft.DLog("send commend to server error: %v", err.Error())
// 					continue
// 				}
// 			}
// 		}
// 	}
// }

func (kvCli *ShardKVClient) SendBucketRpcCommand(args *tinnraftpb.BucketOpArgs) string {
	for {
		if servers, ok := kvCli.config.Groups[int(args.GroupId)]; ok {
			for _, serverAddr := range servers { //获取当前分组对应的集群server addr，尝试发送
				kvCli.client = tinnraft.MakeClientEnd(99, serverAddr)
				reply, err := (*kvCli.client.GetRaftServiceCli()).DoBucket(context.Background(), args)
				if reply.ErrCode == common.ErrCodeNoErr {
					if reply != nil {
						return string(reply.BucketData)
					} else {
						return ""
					}
				} else {
					//tinnraft.DLog("send commend to server error,", err.Error())

					if reply == nil {
						tinnraft.DLog("reply is nil,", reply)
						return ""
					}

					if reply.ErrCode == common.ErrCodeWrongLeader {
						tinnraft.DLog("get the leader id: %v", reply.LeaderId)
						kvCli.client = tinnraft.MakeClientEnd(99, servers[reply.LeaderId])
						reply, err = (*kvCli.client.GetRaftServiceCli()).DoBucket(context.Background(), args)
						if err != nil {
							tinnraft.DLog("leader in sharedserver has down: %v", err.Error())
							return ""
						}
						if reply.ErrCode == common.ErrCodeNoErr && reply != nil {
							return string(reply.BucketData)
						} else if reply.ErrCode == common.ErrCodeNotReady {
							tinnraft.DLog("ErrNotReady: %v", err.Error())
							return ""
						} else if reply.ErrCode == common.ErrCodeCopyBuckets {
							tinnraft.DLog("ErrorCopyError: %v", err.Error())
							return ""
						}
					} else {
						return ""
					}

				}
			}
		}
	}
}

func (cli *ShardKVClient) WriteConnsList(serverAddr string, rpcCli *tinnraft.ClientEnd) {
	cli.connectsList[serverAddr] = rpcCli
}

func (cli *ShardKVClient) ReadConnsList(serverAddrs string) *tinnraft.ClientEnd {
	if conn, ok := cli.connectsList[serverAddrs]; ok {
		return conn
	}
	return nil
}

func (cli *ShardKVClient) CloseRpcConn() {
	cli.client.CloseConns()
}

func (cli *ShardKVClient) GetClient() *config_server.ConfigClient {
	return cli.configCli
}

func (cli *ShardKVClient) GetRpcClient() *tinnraft.ClientEnd {
	return cli.client
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	return bigx.Int64()
}
