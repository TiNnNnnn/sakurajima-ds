package api_gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sakurajima-ds/config_server"
)

func leave(w http.ResponseWriter, r *http.Request, configaddrs []string) {
	fmt.Printf("[configserver addr: %v]\n", configaddrs)
	gid := GetGroupIdFromHeader(r.Header)
	if gid == 0 {
		fmt.Println("illage gid")
		return
	}
	cfgCli := config_server.MakeCfgSvrClient(99, configaddrs)
	cfgCli.Leave([]int64{int64(gid)})

	lastConf := cfgCli.Query(-1)
	if lastConf != nil {
		outBytes, _ := json.Marshal(lastConf)
		log.Printf("Leave group sucess ,last config has change to: %v\n", string(outBytes))
		w.Write([]byte("leave group sucess, last config has change to: " + string(outBytes)))
		w.Header().Add("res", "success")
		w.Header().Add("config", string(outBytes))
		return
	}
	w.Header().Add("res", "failed")
	w.Write([]byte("leave group failed"))

}
