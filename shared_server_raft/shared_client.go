package shared_server

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
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
		BucketOpType:  tinnraftpb.BucketOpType(tinnraftpb.OpType_DeleteBuckets),
		GroupId:       int64(groupId),
		ConfigVersion: int64(kvCli.config.Version),
		BucketIds:     bucketIds,
	})
}

func (kvCli *ShardKVClient) InsertBucketDatas(gid int, bucketIds []int64, datas []byte) string {
	return kvCli.SendBucketRpcCommand(&tinnraftpb.BucketOpArgs{
		BucketOpType:  tinnraftpb.BucketOpType(tinnraftpb.OpType_InsertBuckets),
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
	groupId := kvCli.config.Buckets[bucket_id]
	if groupId == 0 {
		return "", errors.New("there is no shared in charge of this bucke")
	}
	if servers, ok := kvCli.config.Groups[groupId]; ok {

		for _, addrs := range servers {
			if kvCli.ReadConnsList(addrs) == nil {
				kvCli.client = tinnraft.MakeClientEnd(99, addrs)
			} else {
				kvCli.client = kvCli.ReadConnsList(addrs)
			}

			reply, err := (*kvCli.client.GetRaftServiceCli()).DoCommand(context.Background(), args)
			if err != nil {
				tinnraft.DLog("a node in cluster has down,begin try a another one")
				continue
			}

			switch reply.ErrCode {
			case common.ErrCodeNoErr:
				kvCli.commandId++
				return reply.Value, nil
			case common.ErrCodeWrongGroup:
				kvCli.config = kvCli.configCli.Query(-1)
				return "", errors.New("Wrong Group")
			case common.ErrCodeWrongLeader:
				kvCli.client = tinnraft.MakeClientEnd(999, servers[reply.LeaderId])
				if err != nil {
					fmt.Printf("err: %s", err.Error())
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
	return "", errors.New("Unkonw Code")
}

func (kvCli *ShardKVClient) SendBucketRpcCommand(args *tinnraftpb.BucketOpArgs) string {
	for {
		if servers, ok := kvCli.config.Groups[int(args.GroupId)]; ok {
			for _, serverAddr := range servers {
				kvCli.client = tinnraft.MakeClientEnd(99, serverAddr)
				reply, err := (*kvCli.client.GetRaftServiceCli()).DoBucket(context.Background(), args)
				if err == nil {
					if reply != nil {
						return string(reply.BucketData)
					} else {
						return ""
					}
				} else {
					tinnraft.DLog("send commend to server error: %v", err.Error())
					return ""
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
