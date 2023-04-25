package api_gateway

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

func CopyGroup(groups map[int][]string) map[int][]string {
	newGroup := make(map[int][]string)
	for groupId, addrs := range groups {
		newAddrs := make([]string, len(addrs))
		copy(newAddrs, addrs)
		newGroup[groupId] = newAddrs
	}
	return newGroup
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
	Time     int64
}
