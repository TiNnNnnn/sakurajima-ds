package api_gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sakurajima-ds/config_server"
)

func query(w http.ResponseWriter, r *http.Request, configaddrs []string) {
	fmt.Printf("[configserver addr: %v]\n", configaddrs)

	cfgCli := config_server.MakeCfgSvrClient(99, configaddrs)
	lastConf := cfgCli.Query(-1)
	var outBytes = []byte{}
	if lastConf != nil {
		outBytes, _ = json.Marshal(lastConf)
		log.Println("Last Config: " + string(outBytes))
		w.Write([]byte("get last config sucess! config: " + string(outBytes)))
	}

	w.Write([]byte("get last config failed,all configservers have down"))

}
