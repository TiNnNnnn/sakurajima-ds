package api_gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	api_gateway "sakurajima-ds/api_gateway"
	shared_server "sakurajima-ds/shared_server_raft"
	"strconv"
	"strings"
)

type CmdReply struct {
	Success bool
	Msg     []byte
}

func Handler(w http.ResponseWriter, r *http.Request, as *api_gateway.ApiLogServer) {
	if len(strings.Split(r.URL.EscapedPath(), "/")) < 2 {
		log.Println("wrong args in urls")
		return
	}

	m := r.Method
	op := strings.Split(r.URL.EscapedPath(), "/")[2]
	fmt.Printf("apis op: %v\n", op)

	if m == http.MethodPut {
		switch op {
		case "put":
			put(w, r, as)
			return
		}
	}
	if m == http.MethodGet {
		switch op {
		case "get":
			get(w, r, as)
			return
		case "query":
			query(w, r, as)
			return
		}

		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func put(w http.ResponseWriter, r *http.Request, as *api_gateway.ApiLogServer) {

	key := GetKeyFromHeader(r.Header)
	value := GetValueFromHeader(r.Header)

	fmt.Printf("key: %v value: %v\n", key, value)

	if key == "" || value == "" {
		log.Println("key or value is empty,failed")
		w.Write([]byte("key or value is empty,failed"))
		return
	}

	//find configServer addrs
	curAddrsCfg, err := as.Stm.Query(-1)
	if err != nil {
		log.Println("query lastest log failed")
		return
	}
	ConfigPeersMap := curAddrsCfg.Cfg_server_addr
	cfgstring := addrsList2str(ConfigPeersMap)

	fmt.Printf("config addrs: %v\n", cfgstring)

	//向sharedServer 存储数据
	skvclient := shared_server.MakeSharedKvClient(cfgstring)
	puterr := skvclient.Put(key, value)
	if puterr != nil {
		fmt.Println("err: " + puterr.Error())
		return
	}
	fmt.Printf("put data sucess\n")

	res := &CmdReply{
		Success: true,
		Msg:     []byte("put kv [" + key + ":" + value + "] success"),
	}
	resBytes, _ := json.Marshal(res)
	w.Write(resBytes)
}

func get(w http.ResponseWriter, r *http.Request, as *api_gateway.ApiLogServer) {

	key := GetKeyFromHeader(r.Header)
	if key == "" {
		log.Println("key is empty,failed")
		w.Write([]byte("key is empty,failed"))
		return
	}

	//find configServer addrs
	curAddrsCfg, err := as.Stm.Query(-1)
	if err != nil {
		log.Println("query lastest log failed")
		return
	}
	ConfigPeersMap := curAddrsCfg.Cfg_server_addr
	cfgstring := addrsList2str(ConfigPeersMap)

	skvclient := shared_server.MakeSharedKvClient(cfgstring)
	v, err := skvclient.Get(key)
	if err != nil {
		fmt.Println("err: " + err.Error())
	}
	log.Printf("get the value %v success", v)
	w.Write([]byte("get value " + v + " success"))
}

func query(w http.ResponseWriter, r *http.Request, as *api_gateway.ApiLogServer) {
	bliststr := GetBucketListFromHeader(r.Header)
	gid, _ := strconv.Atoi(GetGroupIdFromHeader(r.Header))

	//find configServer addrs
	curAddrsCfg, err := as.Stm.Query(-1)
	if err != nil {
		log.Println("query lastest log failed")
		return
	}

	ConfigPeersMap := curAddrsCfg.Cfg_server_addr
	cfgstring := addrsList2str(ConfigPeersMap)

	bids := []int64{}
	blist := strings.Split(bliststr, ",")
	for _, b := range blist {
		bid, _ := strconv.Atoi(b)
		bids = append(bids, int64(bid))
	}

	skvclient := shared_server.MakeSharedKvClient(cfgstring)
	datas := skvclient.GetBucketDatas(gid, bids)

	log.Printf("get buckets data: %v", datas)

	w.Write([]byte("get buckets data success,data:" + datas))

	res := &CmdReply{
		Success: true,
		Msg:     []byte(datas),
	}
	resBytes, _ := json.Marshal(res)
	w.Write(resBytes)

}

func addrsList2str(ConfigPeersMap []string) string {

	cfgstring := ""
	for _, addr := range ConfigPeersMap {
		cfgstring += addr
		cfgstring += ","
	}

	if cfgstring[len(cfgstring)-1] == ',' {
		cfgstring = cfgstring[0 : len(cfgstring)-1]
	}

	return cfgstring
}

func GetKeyFromHeader(h http.Header) string {
	key := h.Get("key")
	return key
}

func GetValueFromHeader(h http.Header) string {
	value := h.Get("value")
	return value
}

func GetBucketListFromHeader(h http.Header) string {
	blist := h.Get("blist")
	b := strings.Split(blist, ",")
	if len(b) <= 0 {
		log.Println("no bid in blist")
		return ""
	}

	return blist
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
