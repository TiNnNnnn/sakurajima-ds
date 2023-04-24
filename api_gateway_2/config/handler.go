package api_gateway

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

var configPeersMap = string("127.0.0.1:8088,127.0.0.1:8089,127.0.0.1:8090")

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Method

	if m == http.MethodPut {
		if len(strings.Split(r.URL.EscapedPath(), "/")) < 3 {
			log.Println("wrong args in urls")
			return
		}
		configOp := strings.Split(r.URL.EscapedPath(), "/")[2]
		log.Printf("configOp: %v\n", configOp)

		if configOp == "join" {
			join(w, r)
		} else if configOp == "leave" {
			leave(w, r)
		} else if configOp == "move" {
			move(w, r)
		}
		return
	}
	if m == http.MethodGet {
		if len(strings.Split(r.URL.EscapedPath(), "/")) != 3 {
			log.Println("wrong args in urls")
			return
		}
		configOp := strings.Split(r.URL.EscapedPath(), "/")[2]
		log.Printf("configOp: %v\n", configOp)

		if configOp == "query" {
			query(w, r)
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
	groupId := h.Get("groupId")
	log.Printf("buckerId: %v\n", groupId)

	gid, _ := strconv.Atoi(groupId)
	if gid < 0 {
		return 0
	}
	return gid
}

func GetSvraddrsFromHeader(h http.Header) string {
	svraddrs := h.Get("svraddrs")
	log.Printf("sharedAddrs: %v\n", svraddrs)

	addrs := strings.Split(svraddrs, ",")

	if len(addrs) < 3 {
		log.Printf("a group of sharded cluster should more than 3 servers")
		return ""
	}
	// ret := []string{}
	// ret = append(ret, addrs...)
	// return ret

	return svraddrs
}
