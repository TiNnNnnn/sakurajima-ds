package objects

import (
	"log"
	"net/http"
	"sakurajima-ds/meta_server"
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

	metaCli := meta_server.MakeMetaSvrClient(99, PeersMap)

	metaCli.PutMetadata(int64(userId), bucketName, objectName, 1, 1234, "abcdef")

}
