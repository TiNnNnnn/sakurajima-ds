package meta_server

import (
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraft"
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
	t.Logf("-----------------test putMetadata----------------")
	stm.PutMetadata(0, "test", "yyk", 1, 1024, "sxsaxs&*saajisoo")
	stm.PutMetadata(0, "test", "yyk", 2, 1099, "sx788xs&*sajlisoo")
	stm.PutMetadata(0, "test", "lgh", 2, 10000, "09080798&*sajlisoo")

	tkey := META + strconv.Itoa(int(2345)) + "_" + "yyk" + "_" + "test"
	ret, _ := stm.engine.GetAllPrefixKey(tkey)

	t.Logf("ret: %v", ret)
	t.Log("\n")

	t.Logf("-----------------test getMetadata----------------")
	ret2, _ := stm.GetMetadata(0, "test", "yyk", 1)
	t.Logf("get result: %v", ret2)
	t.Log("\n")

	t.Logf("-----------------test SearchLastestVersion----------------")
	ret3,_ := stm.SearchLastestVersion(0, "test", "yyk")
	t.Logf("Lastest verison : %v", ret3)
	t.Log("\n")

	t.Logf("-----------------test SearcAllVersion----------------")
	ret4,_ := stm.SearvhAllVersion(0, "test", "lgh",0,100)
	t.Logf("Lastest verison : %v", ret4)
	t.Log("\n")

	t.Logf("-----------------test DelMeataData----------------")
	stm.DelMetadata(0, "test", "lgh",2)
	ret5,_ := stm.SearvhAllVersion(0, "test", "lgh",0,100)
	t.Logf("Lastest verison : %v", ret5)
	t.Log("\n")
}
