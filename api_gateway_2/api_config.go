package api_gateway

type AddrConfig struct {
	Cfg_server_addr    []string
	Shared_server_addr map[int][]string
	CurVerison         int64
}

func MakeDefaultConfig() AddrConfig {
	return AddrConfig{
		Cfg_server_addr:    make([]string, 0),
		Shared_server_addr: make(map[int][]string, 0),
		CurVerison:         0,
	}
}
