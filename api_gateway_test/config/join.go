package config

import (
	"fmt"
	"net/http"
	"sakurajima-ds/config_server"
	"strings"
)

func join(w http.ResponseWriter, r *http.Request) {
	configaddrs := strings.Split(configPeersMap, ",")
	fmt.Printf("[configserver addr: %v]\n", configaddrs)

	gid := GetGroupIdFromHeader(r.Header)
	svradrs := GetSvraddrsFromHeader(r.Header)

	if len(svradrs) == 0 || gid <= 0 {
		return
	}

	cfgCli := config_server.MakeCfgSvrClient(99, configaddrs)

	addrMap := make(map[int64]string)
	addrMap[int64(gid)] = svradrs
	cfgCli.Join(addrMap)
}
