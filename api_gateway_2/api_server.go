package api_gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"strings"
	"sync"

	"github.com/go-cmd/cmd"
	"github.com/gorilla/websocket"
)

type ApiLogServer struct {
	Mu       sync.Mutex
	BufferMu sync.RWMutex
	LogChan  chan string              //unprduce log channel
	MutiChan chan *tinnraftpb.LogArgs // format json log channel
	stm      *AddrStateMachine        //restore configsever,sharedserver addrs
	tinnraftpb.UnimplementedRaftServiceServer
}

func MakeApiGatwayServer(saddr string) *ApiLogServer {

	addrEngine := storage_engine.EngineFactory("leveldb", "./saddr_data/"+"api_server/")
	apiServer := &ApiLogServer{
		LogChan:  make(chan string),
		MutiChan: make(chan *tinnraftpb.LogArgs),
		stm:      MakeAddrConfigStm(addrEngine),
	}
	return apiServer
}

// 将不同类别的日志发送给客户端
func (as *ApiLogServer) SendMutiLog(c *websocket.Conn) {
	for {
		mutiLog := <-as.MutiChan

		logbytes, err := json.Marshal(mutiLog)
		//log.Printf("mutilogbytes: {%v %v %s %s}\n", mutiLog.Op, mutiLog.Contents, mutiLog.FromId, mutiLog.ToId)
		log.Printf("mutiLog: %v", mutiLog)
		if err != nil {
			log.Println("log json marshal failed")
		}
		as.Mu.Lock()
		c.WriteMessage(1, logbytes)
		as.Mu.Unlock()
	}
}

// 将未处理日志发送给客户端
func (as *ApiLogServer) SendCommonLog(c *websocket.Conn) {
	for {
		cLog := <-as.LogChan
		if len(cLog) > 0 {
			log.Printf("*cmd*: %v\n", cLog)
			as.Mu.Lock()
			c.WriteMessage(1, []byte(cLog))
			as.Mu.Unlock()
		}
	}
}

func (as *ApiLogServer) ReadStdoutAndStderr(out bytes.Buffer) {
	for {
		l, err := out.ReadBytes('\n')
		if err != nil && err.Error() != "EOF" {
			log.Print(err)
			continue
		}
		if len(l) == 0 {
			continue
		}
		as.LogChan <- string(l)
	}
}

func (as *ApiLogServer) DoLog(ctx context.Context, args *tinnraftpb.LogArgs) (*tinnraftpb.LogReply, error) {
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

// 启动KvServer
func (as *ApiLogServer) StartKvServer(w http.ResponseWriter, r *http.Request) {

	sid := GetServerIdFromHeader(r.Header)
	if sid == "" {
		return
	}
	//./../../output/kvserver
	cmd := exec.Command("./../../output/kvserver", sid)

	cmd.Stdin = os.Stdin
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	go func() {
		for {
			l, err := out.ReadBytes('\n')
			if err != nil && err.Error() != "EOF" {
				log.Print(err)
				continue
			}
			if len(l) == 0 {
				continue
			}
			as.LogChan <- string(l)
		}
	}()

	err := cmd.Run()
	if err != nil {
		fmt.Println("failed to begin the kvserver!")
	}
}

// 启动ConfigServer
func (as *ApiLogServer) StartConfigServer(w http.ResponseWriter, r *http.Request) {
	sid := GetServerIdFromHeader(r.Header)
	if sid == "" {
		return
	}
	cfg_addrs := GetGroupAddrsFromHeader(r.Header, "cfg_addrs")

	cfgAddrs := strings.Split(cfg_addrs, ",")
	//集群中节点数量过少
	if len(cfgAddrs) < 3 {
		fmt.Println("too low nodes in current config cluster")
		w.Write([]byte("too low nodes in current config cluster"))
		return
	}

	cfgAddrMap := make(map[int]string)
	for i, addr := range cfgAddrs {
		cfgAddrMap[i] = addr
	}
	curConfig, err := as.stm.Query(-1)
	if err != nil {
		fmt.Println("read lastest config failed")
	}
	//更新 configaddrs数据 并持久化
	newConfig := &AddrConfig{
		Cfg_server_addr: cfgAddrMap,
		CurVerison:      curConfig.CurVerison + 1,
	}
	as.stm.Update(*newConfig)

	//启动configserver
	fmt.Println("START: [./../../output/cfgserver " + sid + " " + cfg_addrs + "]")
	cmd := exec.Command("./../../output/cfgserver", sid, cfg_addrs)

	cmd.Stdin = os.Stdin
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	//获取server stdout
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

	runErr := cmd.Run()
	if runErr != nil {
		fmt.Println("failed to begin the configserver!")
	}
}

// 启动SharedServer
func (as *ApiLogServer) StartSharedServer(w http.ResponseWriter, r *http.Request) {
	sid := GetServerIdFromHeader(r.Header)
	gid := GetGroupIdFromHeader(r.Header)
	if sid == "" || gid == "" {
		return
	}

	cfg_addrs := GetGroupAddrsFromHeader(r.Header, "cfg_addrs")
	shared_addrs := GetGroupAddrsFromHeader(r.Header, "shared_addrs")

	cfgAddrs := strings.Split(cfg_addrs, ",")
	if len(cfgAddrs) < 3 {
		fmt.Println("too low nodes in current config cluster")
		w.Write([]byte("too low nodes in current config cluster"))
		return
	}

	sharedAddrs := strings.Split(shared_addrs, ",")
	if len(sharedAddrs) < 3 {
		fmt.Println("too low nodes in current shared cluster")
		w.Write([]byte("too low nodes in current shared cluster"))
		return
	}

	//启动sharedserver
	fmt.Println("START: [./../../output/sharedserver " + sid + " " + gid + " " + cfg_addrs + " " + shared_addrs + "]")
	cmd := exec.Command("./../../output/sharedserver", sid, gid, cfg_addrs, shared_addrs)

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

	err := cmd.Run()
	if err != nil {
		fmt.Println("failed to begin the sharedserver!")
	}
}

func (as *ApiLogServer) StopServer(w http.ResponseWriter, r *http.Request) {

	stype := GetstypeFromHeader(r.Header)
	sid, _ := strconv.Atoi(GetstypeFromHeader(r.Header))

	addrs, err := as.stm.Query(-1)
	if err != nil {
		fmt.Println("get the lastest log failed; err: " + err.Error())
		return
	}

	var addr = ""
	if stype == "configserver" {
		addr = addrs.Cfg_server_addr[sid]
	} else if stype == "sharedserver" {
		addr = addrs.Shared_server_addr[0][sid]
	}

	port := strings.Split(addr, ":")[1]

	c := cmd.NewCmd("bash", "-c", "netstat -nltp |grep "+port)
	<-c.Start()
	fmt.Printf("nestat result: %v", c.Status().Stdout)

	fields := strings.Fields(c.Status().Stdout[0])
	if len(fields) == 0 {
		fmt.Printf("can't find proc with port %v", port)
		return
	}

	//取出进程pid
	pid := strings.Split(fields[6], "/")[0]

	//log.Printf("get the process %v %v with pid %v",)
	//杀死进程
	prc := exec.Command("kill", "-9", pid)
	out, err := prc.Output()
	if err != nil {
		fmt.Printf("kill proc with port %v failed", port)
		panic(err)
	}
	fmt.Printf("kill proc with port %v success! %v", port, string(out))

}

func GetServerIdFromHeader(h http.Header) string {
	kvs_id := h.Get("s_id")
	if id, _ := strconv.Atoi(kvs_id); id < 0 {
		log.Println("a illegal serverId! it should be more than 0")
		return ""
	}
	return kvs_id
}
func GetstypeFromHeader(h http.Header) string {
	stype := h.Get("stype")
	return stype
}

func GetGroupIdFromHeader(h http.Header) string {
	gid := h.Get("gid")
	gidstring, _ := strconv.Atoi(gid)
	if gidstring < 0 {
		fmt.Println("a illage gid! it should be more than 0")
		return ""
	}
	return gid
}

func GetGroupAddrsFromHeader(h http.Header, addrType string) string {
	var addrs = ""
	switch addrType {
	case "cfg_addrs":
		addrs = h.Get("cfg_addrs")
	case "shared_addrs":
		addrs = h.Get("shared_addrs")
	}
	if len(addrs) <= 0 {
		fmt.Println("a illagal addrs!")
	}
	return addrs
}
