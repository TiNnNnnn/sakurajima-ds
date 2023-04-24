package api_gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sakurajima-ds/config_server"
	"strings"
)


func query(w http.ResponseWriter, r *http.Request) {
	addrs := strings.Split(configPeersMap, ",")
	fmt.Printf("[configserver addr: %v]\n", addrs)

	cfgCli := config_server.MakeCfgSvrClient(99, addrs)
	lastConf := cfgCli.Query(-1)
	outBytes, _ := json.Marshal(lastConf)
	log.Println("last config: " + string(outBytes))

}
