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
	"sakurajima-ds/common"

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
	CfgMu    sync.Mutex
	LogChan  chan string              //unprduce log channel
	MutiChan chan *tinnraftpb.LogArgs // format json log channel
	StopChan chan bool
	cfgCond  *sync.Cond
	Stm      *AddrStateMachine //restore configsever,sharedserver addrs
	tinnraftpb.UnimplementedRaftServiceServer
}

func MakeApiGatwayServer(saddr string) *ApiLogServer {

	addrEngine := storage_engine.EngineFactory("leveldb", "./saddr_data/"+"api_server/")
	apiServer := &ApiLogServer{
		LogChan:  make(chan string),
		MutiChan: make(chan *tinnraftpb.LogArgs, 1024),
		StopChan: make(chan bool),
		Stm:      MakeAddrConfigStm(addrEngine),
	}

	apiServer.cfgCond = sync.NewCond(&apiServer.CfgMu)

	return apiServer
}

// 将不同类别的日志发送给客户端
func (as *ApiLogServer) SendMutiLog(c *websocket.Conn) {
	for {
		mutiLog := <-as.MutiChan
		sType := common.GetNameBypId(int(mutiLog.Pid))
		groupId := 0

		if sType == "cfgserver" {
			sType = "configserver"
		} else if sType == "sharedserver" {
			curAddrCfg, _ := as.Stm.Query(-1)
			saddrs := curAddrCfg.Shared_server_addr

			addr := ""
			for i := 0; i < 15; i++ {
				addr = common.GetGroupIdBypId2(int(mutiLog.Pid), "./../../outp")
				if addr != "" {
					break
				}
			}
			if addr == "" {
				return
			}

			for gid, addrs := range saddrs {
				for _, adr := range addrs {
					if adr == addr {
						groupId = gid
						break
					}
				}
			}
		}

		log.Printf("mutilogbytes: {%v %v %v %v %v %v %v %v}\n", mutiLog.Op, mutiLog.Contents, mutiLog.FromId, mutiLog.ToId, mutiLog.PreState, mutiLog.CurState, sType, groupId)

		hblog := HBLog{
			Logtype:  mutiLog.Op.String(),
			Content:  mutiLog.Contents,
			From:     int(mutiLog.FromId),
			To:       int(mutiLog.ToId),
			PreState: mutiLog.PreState,
			CurState: mutiLog.CurState,
			SvrType:  sType,
			GroupId:  groupId,
			Time:     mutiLog.Time,
		}

		logBytes, _ := json.Marshal(hblog)

		err := c.WriteMessage(1, logBytes)
		if err != nil {
			log.Printf("writemessgae err: %v", err.Error())
			break
		}
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

// 启动KvServer test only
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

	curConfig, err := as.Stm.Query(-1)
	if err != nil {
		fmt.Println("read lastest config failed")
	}

	//更新 configaddrs数据 并持久化
	newConfig := &AddrConfig{
		Shared_server_addr: curConfig.Shared_server_addr,
		Cfg_server_addr:    cfgAddrMap,
		CurVerison:         curConfig.CurVerison + 1,
	}
	if !IsEqual(&curConfig, newConfig) {
		as.Stm.Update(*newConfig)
	}

	Config, _ := as.Stm.Query(-1)
	ShowCurConfig(&Config)

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
			//out.Grow(1024)
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

	curConfig, err := as.Stm.Query(-1)
	if err != nil {
		fmt.Println("read lastest config failed")
	}

	newGroupId, _ := strconv.Atoi(gid)
	newsharedConfig := CopySharedCOnfig(curConfig.Shared_server_addr)
	newsharedConfig[newGroupId] = sharedAddrs

	//更新 sharedaddrs数据 并持久化
	newConfig := &AddrConfig{
		Shared_server_addr: newsharedConfig,
		Cfg_server_addr:    curConfig.Cfg_server_addr,
		CurVerison:         curConfig.CurVerison + 1,
	}

	if !IsEqual(&curConfig, newConfig) {
		as.Stm.Update(*newConfig)
	}

	Config, _ := as.Stm.Query(-1)
	ShowCurConfig(&Config)

	//启动sharedserver
	fmt.Println("START: [./../../output/sharedserver " + sid + " " + gid + " [" + cfg_addrs + "] [" + shared_addrs + "]]")
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

	runErr := cmd.Run()
	if runErr != nil {
		fmt.Println("failed to begin the sharedserver!")
	}
}

func (as *ApiLogServer) StopServer(w http.ResponseWriter, r *http.Request) {

	stype := GetstypeFromHeader(r.Header)
	sid, _ := strconv.Atoi(GetServerIdFromHeader(r.Header))
	gid, _ := strconv.Atoi(GetGroupIdFromHeader(r.Header))

	addrs, err := as.Stm.Query(-1)
	if err != nil {
		fmt.Println("get the lastest log failed; err: " + err.Error())
		return
	}

	var addr = ""
	if stype == "configserver" {
		addr = addrs.Cfg_server_addr[sid]
	} else if stype == "sharedserver" {
		addr = addrs.Shared_server_addr[gid][sid]
	}

	port := strings.Split(addr, ":")[1]

	c := cmd.NewCmd("bash", "-c", "netstat -nltp |grep "+port)
	<-c.Start()
	fmt.Printf("nestat result: %v", c.Status().Stdout)

	if len(c.Status().Stdout) == 0 {
		fmt.Printf("can't find proc with port %v", port)
		w.Write([]byte("can't find proc with port " + port))
		return
	}

	fields := strings.Fields(c.Status().Stdout[0])
	if len(fields) == 0 {
		fmt.Printf("can't find proc with port %v", port)
		w.Write([]byte("can't find proc with port " + port))
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
		w.Write([]byte("kill proc " + stype + ":" + strconv.Itoa(sid) + "with port " + port + " failed"))
		panic(err)
	}

	fmt.Printf("kill proc with port %v success! %v", port, string(out))
	w.Write([]byte("kill proc " + stype + ":" + strconv.Itoa(sid) + " with port " + port + " success"))

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
