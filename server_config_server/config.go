package serverconfigserver

type Config struct {
	cfg_server_addr  []string
	meta_server_addr []string
	curVerison       int64
}


func MakeDefaultConfig() Config {
	return Config{
		cfg_server_addr:  make([]string, 0),
		meta_server_addr: make([]string, 0),
		curVerison:       0,
	}
}

