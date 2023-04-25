package api_gateway

import (
	"log"
	"reflect"
)

type AddrConfig struct {
	Cfg_server_addr    map[int]string
	Shared_server_addr map[int][]string
	CurVerison         int64
}

func MakeDefaultConfig() AddrConfig {
	return AddrConfig{
		Cfg_server_addr:    make(map[int]string, 0),
		Shared_server_addr: make(map[int][]string, 0),
		CurVerison:         0,
	}
}

func CopySharedCOnfig(groups map[int][]string) map[int][]string {
	newGroup := make(map[int][]string)
	for groupId, addrs := range groups {
		newAddrs := make([]string, len(addrs))
		copy(newAddrs, addrs)
		newGroup[groupId] = newAddrs
	}
	return newGroup
}

func IsEqual(precfg, curcfg *AddrConfig) bool {

	if !reflect.DeepEqual(precfg.Cfg_server_addr, curcfg.Cfg_server_addr) {
		return false
	}

	if !reflect.DeepEqual(precfg.Shared_server_addr, precfg.Shared_server_addr) {
		return false
	}

	return true
}

func ShowCurConfig(cfg *AddrConfig) {
	log.Println("-----------cur addrconfig ----------------")
	log.Printf("curaddrConfig: %v", cfg)
	log.Println("------------------------------------------")
}

// HeartBeat
type HBLog struct {
	Logtype  string
	Content  string
	From     int
	To       int
	PreState string
	CurState string
	SvrType  string
	GroupId  int
	Time     int64
}
