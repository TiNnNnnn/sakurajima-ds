package api_gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sakurajima-ds/tinnraftpb"
	"strconv"

	"github.com/gorilla/websocket"
)

type ApiGatwayServer struct {
	//mu         sync.Mutex
	LogChan  chan string
	MutiChan chan *tinnraftpb.LogArgs

	tinnraftpb.UnimplementedRaftServiceServer
}

func MakeApiGatwayServer(saddr string) *ApiGatwayServer {

	apiServer := &ApiGatwayServer{
		LogChan:  make(chan string),
		MutiChan: make(chan *tinnraftpb.LogArgs),
	}
	return apiServer
}

// 启动KvServer
func (as *ApiGatwayServer) StartKvServer(w http.ResponseWriter, r *http.Request) {

	sid := GetServerIdFromHeader(r.Header)
	if sid == "" {
		return
	}
	cmd := exec.Command("./../output/kvserver", sid)

	cmd.Stdin = os.Stdin
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	go func() {
		for {
			l, err := out.ReadString('\n')
			if err != nil && err.Error() != "EOF" {
				log.Print(err)
				continue
			}
			as.LogChan <- l

		}
	}()
	cmd.Run()
}

// 将不同类别的日志发送给客户端
func (as *ApiGatwayServer) SendMutiLog(c *websocket.Conn) {
	for {
		mutiLog := <-as.MutiChan
		logbytes, err := json.Marshal(mutiLog)
		log.Printf("mutilogbytes: %v", logbytes)
		if err != nil {
			log.Println("log json marshal failed")
		}
		c.WriteMessage(1, logbytes)
	}
}

func (as *ApiGatwayServer) DoLog(ctx context.Context, args *tinnraftpb.LogArgs) (*tinnraftpb.LogReply, error) {
	reply := tinnraftpb.LogReply{}
	if args != nil {
		reply.Errcode = 0
		reply.ErrMsg = ""
		reply.Success = true
		as.MutiChan <- args
		return &reply, nil
	} else {
		reply.Errcode = 10
		reply.ErrMsg = "args is empty"
		reply.Success = false
		return &reply, errors.New("args is empty")
	}
}

func GetServerIdFromHeader(h http.Header) string {
	kvs_id := h.Get("kvs_id")
	if id, _ := strconv.Atoi(kvs_id); id < 0 {
		log.Println("a illegal serverId! it should be more than 0")
		return ""
	}
	return kvs_id
}
