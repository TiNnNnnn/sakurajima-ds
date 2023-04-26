package config_server

/*
	记录数据服务器的分组情况
*/

import (
	"sakurajima-ds/common"
)

type Config struct {
	//mu      sync.Mutex
	Version  int
	Buckets  [common.BucketsNum]int
	Groups   map[int][]string
	LeaderId int //for log ,no another meanings
}

func MakeDefaultConfig() Config {
	return Config{
		Groups: make(map[int][]string),
	}
}

func CopyGroup(groups map[int][]string) map[int][]string {
	newGroup := make(map[int][]string)
	for groupId, addrs := range groups {
		newAddrs := make([]string, len(addrs))
		copy(newAddrs, addrs)
		newGroup[groupId] = newAddrs
	}
	return newGroup
}
