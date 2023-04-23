package objects

import (
	"fmt"
	"log"
	"net/http"
	shared_server "sakurajima-ds/shared_server_raft"
	"strconv"
	"strings"
)


func put(w http.ResponseWriter, r *http.Request) {
	if len(strings.Split(r.URL.EscapedPath(), "/")) < 4 {
		log.Println("wrong args in urls")
		return
	}
	userId, _ := strconv.Atoi(strings.Split(r.URL.EscapedPath(), "/")[2])
	bucketName := strings.Split(r.URL.EscapedPath(), "/")[3]
	objectName := strings.Split(r.URL.EscapedPath(), "/")[4]
	log.Printf("userId: %v,bucketName: %v,objectName:%v", userId, bucketName, objectName)

	// metaCli := meta_server.MakeMetaSvrClient(99, PeersMap)

	// metaCli.PutMetadata(int64(userId), bucketName, objectName, 1, 1234, "abcdef")

	//向sharedServer 存储数据
	skvclient := shared_server.MakeSharedKvClient(ShardPeersMap)
	err := skvclient.Put("232343882", "hahahaha")
	if err != nil {
		fmt.Println("err: " + err.Error())
	}


}
