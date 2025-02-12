package api_gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sakurajima-ds/config_server"
)

func move(w http.ResponseWriter, r *http.Request, configaddrs []string) {
	fmt.Printf("[configserver addr: %v]\n", configaddrs)

	bLists := GetBucketsFromHeader(r.Header)
	gid := GetGroupIdFromHeader(r.Header)

	if len(bLists) != 2 || gid < 0 {
		return
	}

	log.Printf("begin move buckets [%v-%v] to group %v\n", bLists[0], bLists[1], gid)

	cfgCli := config_server.MakeCfgSvrClient(99, configaddrs)
	for i := bLists[0]; i <= bLists[1]; i++ {
		cfgCli.Move(i, gid)
	}
	
	lastConf := cfgCli.Query(-1)
	if lastConf != nil {
		outBytes, _ := json.Marshal(lastConf)
		log.Printf("last config has change to: %v\n", string(outBytes))
		w.Write([]byte("move sucess, last config has change to: " + string(outBytes)))
		w.Header().Add("res", "success")
		w.Header().Add("config", string(outBytes))
		return
	}

	w.Header().Add("res", "failed")
	w.Write([]byte("move group failed"))

}
