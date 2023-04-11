package config_server

import (
	"sakurajima-ds/common"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraft"
	"testing"
)

type testConfig struct {
	Version int
	Buckets [common.BucketsNum]int
	Groups  map[int][]string
}

// 测试添加分组更新Buckets逻辑
func TestTmp(t *testing.T) {
	tmp := make(map[int][]int)

	config := testConfig{Groups: make(map[int][]string)}

	for groupId := range config.Groups {
		tmp[groupId] = make([]int, 0)
	}
	// [0:[],1:[],2:[],3:[]...]
	for bucketId, groupId := range config.Buckets {
		tmp[groupId] = append(tmp[groupId], bucketId)
	}

	t.Log("tmp:", tmp)

	var newBuckets [common.BucketsNum]int
	for gid, buckets := range tmp {
		t.Logf("gid: %v, buckets: %v", gid, buckets)
		for _, bid := range buckets {
			newBuckets[bid] = gid
		}
	}

	t.Log("newBuckets: ", newBuckets)
}

func testJoin(id int, addrs []string) *map[int][]string {
	newGroup := make(map[int][]string)
	newAddrs := make([]string, 3)
	copy(newAddrs, addrs)
	newGroup[2] = newAddrs
	return &newGroup
}

//测试添加新分组功能
func TestAddGroups(t *testing.T) {
	newEngine, err := storage_engine.MakeLevelDBKvStorage("./conf_data/" + "/test")
	if err != nil {
		tinnraft.DLog("build storage engine failer: %v", err.Error())
		return
	}

	stm := MakeConfigStm(newEngine)

	conf, _ := stm.Query(-1)
	t.Logf("%v", conf)

	//添加新分组
	addr := []string{"9088", "9089", "9090"}
	newGroup := testJoin(2, addr)
	stm.Join(*newGroup)

	stm.Move(0, 2)

	conf2, _ := stm.Query(-1)

	t.Logf("%v", conf2)

}
