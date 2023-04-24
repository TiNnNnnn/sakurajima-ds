package api_gateway

import (
	"fmt"
	"log"
	"net/http"
	shared_server "sakurajima-ds/shared_server_raft"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	if len(strings.Split(r.URL.EscapedPath(), "/")) < 4 {
		log.Println("wrong args in urls")
		return
	}
	userId, _ := strconv.Atoi(strings.Split(r.URL.EscapedPath(), "/")[2])
	bucketName := strings.Split(r.URL.EscapedPath(), "/")[3]
	objectName := strings.Split(r.URL.EscapedPath(), "/")[4]
	log.Printf("userId: %v,bucketName: %v,objectName:%v", userId, bucketName, objectName)

	skvclient := shared_server.MakeSharedKvClient(ShardPeersMap)
	v, err := skvclient.Get("232343882")
	if err != nil {
		fmt.Println("err: " + err.Error())
	}
	log.Printf("get the value %v success", v)
}
