package api_gateway

import (
	"fmt"
	"log"
	"net/http"

	//cfg_server "sakurajima-ds/config_server"
	shared_server "sakurajima-ds/shared_server_raft"
	"strings"
)

var ConfigPeersMap = string("127.0.0.1:8088,127.0.0.1:8089,127.0.0.1:8090")

func ObjectHandler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func put(w http.ResponseWriter, r *http.Request) {
	if len(strings.Split(r.URL.EscapedPath(), "/")) < 3 {
		log.Println("wrong args in urls")
		return
	}

	//向sharedServer 存储数据
	skvclient := shared_server.MakeSharedKvClient(ConfigPeersMap)
	err := skvclient.Put("232343882", "hahahaha")
	if err != nil {
		fmt.Println("err: " + err.Error())
	}

}

func get(w http.ResponseWriter, r *http.Request) {
	if len(strings.Split(r.URL.EscapedPath(), "/")) < 3 {
		log.Println("wrong args in urls")
		return
	}
	skvclient := shared_server.MakeSharedKvClient(ConfigPeersMap)
	v, err := skvclient.Get("232343882")
	if err != nil {
		fmt.Println("err: " + err.Error())
	}
	log.Printf("get the value %v success", v)
}