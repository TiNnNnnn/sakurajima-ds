package api_gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sakurajima-ds/config_server"
)

func join(w http.ResponseWriter, r *http.Request, configaddrs []string) {

	fmt.Printf("[configserver addr: %v]\n", configaddrs)

	gid := GetGroupIdFromHeader(r.Header)
	svradrs := GetSvraddrsFromHeader(r.Header)

	if len(svradrs) == 0 || gid <= 0 {
		return
	}

	cfgCli := config_server.MakeCfgSvrClient(99, configaddrs)

	addrMap := make(map[int64]string)
	addrMap[int64(gid)] = svradrs

	ret := cfgCli.Join(addrMap)

	lastConf := cfgCli.Query(-1)
	var outBytes = []byte{}

	if ret {
		if lastConf != nil {
			outBytes, _ = json.Marshal(lastConf)
			log.Println("Join Success,cur Config: " + string(outBytes))
			w.Write([]byte("join last config sucess! cur config: " + string(outBytes)))
		}
	} else {
		if lastConf != nil {
			outBytes, _ = json.Marshal(lastConf)
			log.Println("Join Failed,cur Config: " + string(outBytes))
			w.Write([]byte("join last config failed! cur config: " + string(outBytes)))
		} else {
			w.Write([]byte("join failed by all configservers down"))
		}
	}

}
