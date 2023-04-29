package api_gateway

import (
	"fmt"
	"reflect"
	"testing"
)

// func TestTmp(t *testing.T) {
// 	newEngine, err := storage_engine.MakeLevelDBKvStorage("./conf_data/" + "/test")
// 	if err != nil {
// 		log.Printf("build storage engine failer: %v", err.Error())
// 		return
// 	}

// 	stm := MakeAddrConfigStm(newEngine)

// 	conf, _ := stm.Query(-1)
// 	t.Logf("%v", conf)

// 	//configAddrs := []string{"127.0.0.1:10088", "127.0.0.1:10089", "127.0.0.1:10090"}
// 	configAddrs := make(map[int]string)
// 	configAddrs[0] = "127.0.0.1:10088"
// 	configAddrs[1] = "127.0.0.1:10089"
// 	configAddrs[2] = "127.0.0.1:10090"

// 	sharderAddrs := make(map[int][]string)

// 	sharderAddrs[0] = []string{"127.0.0.1:10020", "127.0.0.1:10021", "127.0.0.1:10022"}
// 	sharderAddrs[1] = []string{"127.0.0.1:20020", "127.0.0.1:20021", "127.0.0.1:20022"}
// 	sharderAddrs[2] = []string{"127.0.0.1:30020", "127.0.0.1:30021", "127.0.0.1:30022"}

// 	newConf := AddrConfig{
// 		configAddrs,
// 		sharderAddrs,
// 		2,
// 	}

// 	stm.Update(newConf)

// 	conf, _ = stm.Query(-1)
// 	t.Logf("%v", conf)
// }

func TestMapEqual(t *testing.T) {
	configAddrs1 := make(map[int]string)
	configAddrs1[0] = "127.0.0.1:10088"
	configAddrs1[1] = "127.0.0.1:10089"
	configAddrs1[2] = "127.0.0.1:10090"

	configAddrs2 := make(map[int]string)
	configAddrs2[0] = "127.0.0.1:10088"
	configAddrs2[1] = "127.0.0.1:10089"
	configAddrs2[2] = "127.0.0.1:10090"

	ok := reflect.DeepEqual(configAddrs1, configAddrs2)
	fmt.Printf("ansL: %v", ok)
}
