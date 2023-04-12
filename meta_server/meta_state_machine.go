package meta_server

import (
	"encoding/json"
	"log"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraftpb"
	"strconv"
)

type MetaOp interface {
	PutName(bucketName string, objectName string, objectId int64) (int64, error)
	PutObject(objectId int64, blocks tinnraftpb.DataBlocks) (int64, error)
	GetBlockList(key string, begin int, end int) (string, error)
	Del(key string) bool
	List(key string) []int64
}

const META = "meta_"

type MetaStateMachine struct {
	engine storage_engine.KvStorage
}

func MakeMetaStm(engine storage_engine.KvStorage) *MetaStateMachine {
	newStm := &MetaStateMachine{
		engine: engine,
	}
	return newStm
}

func (stm *MetaStateMachine) PutName(
	bucketName string,
	objectName string,
	objectId int64,
) (int64, error) {
	key := bucketName + objectName
	newName := &Name{
		BucketName: bucketName,
		Key:        key,
		ObjectId:   objectId,
	}
	newNameBytes, err := json.Marshal(newName)
	if err != nil {
		log.Fatalln("newName Marshal failed!")
		return -1, err
	}
	err = stm.engine.Put(META+key, string(newNameBytes))
	if err != nil {
		log.Panicln("put newName to meta_stm failed")
		return -1, err
	}
	return objectId, err
}

func (stm *MetaStateMachine) PutObject(objectId int64, blocks *tinnraftpb.DataBlocks) (int64, error) {
	newObject := &Object{
		ObjectId: objectId,
		Blocks:   blocks,
	}
	//log.Printf("%v", newObject)
	newObjectBytes, err := json.Marshal(newObject)
	if err != nil {
		log.Fatalln("newObject Marshal failed!")
		return -1, err
	}
	err = stm.engine.Put("ob_"+strconv.Itoa(int(objectId)), string(newObjectBytes))
	if err != nil {
		log.Panicln("put newObject to meta_stm failed")
		return -1, err
	}
	return objectId, nil
}

func (stm *MetaStateMachine) GetBlockList(key string, num int) ([]int64, error) {
	//查找name表
	//log.Fatalf("key: %v\n", META+key)
	nameBytes, err := stm.engine.Get(META + key)
	if err != nil {
		log.Println("find objectId from NameTable failed")
		return nil, err
	}
	name := &Name{}
	json.Unmarshal([]byte(nameBytes), name)
	objectId := name.ObjectId
	//log.Fatalf("bid: %v\n", objectId)
	//根据objecId查找object表获取blocklist
	objectBytes, err := stm.engine.Get("ob_" + strconv.Itoa(int(objectId)))
	if err != nil {
		log.Println("find blockLIst from ObjectTable failed")
		return nil, err
	}
	object := &Object{}
	json.Unmarshal([]byte(objectBytes), object)
	//log.Fatalf("blockList: %v\n", *object)
	//读取前num个block
	replylist := []int64{}
	for i, block := range object.Blocks.BlockList {
		if i == num || i == len(object.Blocks.BlockList) {
			break
		}
		replylist = append(replylist, block.BlockId)
	}
	return replylist, nil
}
