package api_gateway

import (
	"log"
	"net/http"
	api_gateway "sakurajima-ds/api_gateway_2"
	"strconv"
	"strings"
)

//var configPeersMap = string("127.0.0.1:8088,127.0.0.1:8089,127.0.0.1:8090")

func Handler(w http.ResponseWriter, r *http.Request, as *api_gateway.ApiLogServer) {
	if len(strings.Split(r.URL.EscapedPath(), "/")) < 3 {
		log.Println("wrong args in urls")
		return
	}
	m := r.Method
	configOp := strings.Split(r.URL.EscapedPath(), "/")[2]
	log.Printf("configOp: %v\n", configOp)

	//find configServer addrs
	curAddrsCfg, err := as.Stm.Query(-1)
	if err != nil {
		log.Println("query lastest log failed")
		return
	}

	//查询ConfigServer
	ConfigPeersMap := curAddrsCfg.Cfg_server_addr
	ConfigList := []string{}
	for _, addr := range ConfigPeersMap {
		if len(addr) > 0 {
			ConfigList = append(ConfigList, addr)
		}
	}

	if m == http.MethodPut {

		if configOp == "join" { //加入分组
			join(w, r, ConfigList)
		} else if configOp == "leave" { //删除分组
			leave(w, r, ConfigList)
		} else if configOp == "move" { //流量分配
			move(w, r, ConfigList)
		}
		return
	}
	if m == http.MethodGet {

		if configOp == "query" { //查询配置
			query(w, r, ConfigList)
		}
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetBucketsFromHeader(h http.Header) []int {
	buckets := h.Get("buckets")
	log.Printf("buckets: %v\n", buckets)
	if len(buckets) < 3 {
		log.Printf("a wrong fromt og bucket, a proper one is like :[0-4]")
		return []int{}
	}
	items := strings.Split(buckets, "-")
	if len(items) != 2 {
		log.Printf("a wrong fromt og bucket, a proper one is like :[0-4]")
		return []int{}
	}

	start, _ := strconv.Atoi(items[0])
	end, _ := strconv.Atoi(items[1])

	return []int{start, end}
}

func GetGroupIdFromHeader(h http.Header) int {
	groupId := h.Get("gid")
	//log.Printf("buckerId: %v\n", groupId)

	gid, _ := strconv.Atoi(groupId)
	if gid < 0 {
		return 0
	}
	return gid
}

func GetSvraddrsFromHeader(h http.Header) string {
	svraddrs := h.Get("shared_addrs")

	addrs := strings.Split(svraddrs, ",")

	if len(addrs) < 3 {
		log.Printf("a group of sharded cluster should more than 3 servers")
		return ""
	}
	return svraddrs
}
