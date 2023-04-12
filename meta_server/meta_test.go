package meta_server

import (
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"
	"strconv"
	"testing"
)

func TestMetaPutName(t *testing.T) {
	newEngine, err := storage_engine.MakeLevelDBKvStorage("./meta_data/" + "/test")
	if err != nil {
		tinnraft.DLog("build storage engine failer: %v", err.Error())
		return
	}

	stm := MakeMetaStm(newEngine)

	t.Log("----测试 putName -------")

	stm.PutName("root", "test", 1)
	stm.PutName("root", "yyk", 2)
	//t.Logf("%v", bid)

	k, _ := stm.engine.GetAllPrefixKey(META + "root")
	t.Logf("以meta_root为前缀的内容: %v", k)

	t.Log("----测试 putObject-----------")

	lists := [11]*tinnraftpb.Block{}

	for i := 1; i <= 10; i++ {

		var tmp = tinnraftpb.Block{
			DataTableName: "block_" + strconv.Itoa(i),
			BlockSize:     10,
			BlockId:       int64(i),
		}
		lists[i] = &tmp
	}

	//t.Logf("lists: %v", lists)
	blocks := tinnraftpb.DataBlocks{
		BlockList: lists[1:10],
	}

	stm.PutObject(1, &blocks)
	k, _ = stm.engine.GetAllPrefixKey("ob_1")
	t.Logf("以ob_1为前缀的内容: %v", k)


	t.Log("----测试 GetObjectList-----------")
	
	l,_ := stm.GetBlockList("roottest", 4)
	t.Logf("查找roottest的4个block块 %v", l)

}
