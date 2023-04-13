package meta_server

import (
	"encoding/json"
	"errors"
	"log"
	"sakurajima-ds/storage_engine"
	"strconv"
)

type MetaOp interface {
	SearchLastestVersion(userid int64, bucketName string, name string) (MetaData, error)
	//按照名称和版本作为索引返回元数据
	GetMetadata(userid int64, bucketName string, name string, version int) (MetaData, error)
	PutMetadata(userid int64, bucketName string, name string, version int, size int64, hash string) error
	AddVersion(userid int64, bucketName string, name, hash string, size int64) error
	SearvhAllVersion(userid int64, bucketName string, name string, from int, size int) ([]MetaData, error)
	DelMetadata(userid int64, bucketName string, name string, verison int) error
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

// 上传新的元数据信息
func (stm *MetaStateMachine) PutMetadata(userid int64, bucketName string, name string, version int, size int64, hash string) error {
	newMeta := &MetaData{
		UserId:     userid,
		BucketName: bucketName,
		ObjectName: name,
		Version:    int64(version),
		Size:       size,
		Hash:       hash,
	}
	newMetaBytes, err := json.Marshal(newMeta)
	if err != nil {
		log.Fatalln("newMeta Marshal failed")
		return err
	}
	key := META + strconv.Itoa(int(userid)) + "_" + bucketName + "_" + name + "_" + strconv.Itoa(version)
	err = stm.engine.Put(key, string(newMetaBytes))
	if err != nil {
		log.Fatalln("put new metadata to stm failed")
		return err
	}
	return nil
}

func (stm *MetaStateMachine) GetMetadata(userid int64, bucketName string, name string, version int) (MetaData, error) {
	if version == 0 {
		return stm.SearchLastestVersion(userid, bucketName, name)
	}
	return stm.getMetadata(userid, bucketName, name, version)
}

func (stm *MetaStateMachine) getMetadata(userid int64, bucketName string, name string, version int) (meta MetaData, e error) {
	key := META + strconv.Itoa(int(userid)) + "_" + bucketName + "_" + name + "_" + strconv.Itoa(version)
	metaBytes, err := stm.engine.Get(key)
	if err != nil {
		log.Fatalln("find meta by useid&bucketName&objectName/version faield")
		return MetaData{}, err
	}
	retmeta := &MetaData{}
	json.Unmarshal([]byte(metaBytes), retmeta)
	return *retmeta, nil
}

func (stm *MetaStateMachine) SearchLastestVersion(userid int64, bucketName string, name string) (MetaData, error) {
	key := META + strconv.Itoa(int(userid)) + "_" + bucketName + "_" + name
	metaBytesList, err := stm.engine.GetAllPrefixKey(key)
	if err != nil {
		log.Fatalln("find lastest version from stm failed")
		return MetaData{}, err
	}

	retmeta := &MetaData{}
	var maxVersion int = 0
	for _, v := range metaBytesList {
		tmpmeta := &MetaData{}
		json.Unmarshal([]byte(v), tmpmeta)
		if tmpmeta.Version > int64(maxVersion) {
			maxVersion = int(tmpmeta.Version)
			*retmeta = *tmpmeta
		}
	}
	//log.Fatalf("lasted versionid: %v", maxVersion)
	if maxVersion < 1 {
		err := errors.New("find no meta by any version")
		return MetaData{}, err
	}

	return *retmeta, nil
}

func (stm *MetaStateMachine) AddVersion(userid int64, bucketName, name, hash string, size int64) error {
	metadata, err := stm.SearchLastestVersion(userid, bucketName, name)
	if err != nil {
		return nil
	}
	return stm.PutMetadata(userid, bucketName, name, int(metadata.Version+1), size, hash)
}

func (stm *MetaStateMachine) SearvhAllVersion(userid int64, bucketName, name string, from int, size int) ([]MetaData, error) {
	key := META + strconv.Itoa(int(userid)) + "_" + bucketName + "_" + name
	metaBytesList, err := stm.engine.GetAllPrefixKey(key)
	if err != nil {
		log.Fatalln("find all version from stm failed")
		return nil, err
	}
	retmetalist := make([]MetaData, 0)
	for _, v := range metaBytesList {
		tmpmeta := &MetaData{}
		err = json.Unmarshal([]byte(v), tmpmeta)
		if err != nil {
			log.Fatalln("newMeta Marshal failed")
			return nil, err
		}
		retmetalist = append(retmetalist, *tmpmeta)
	}
	return retmetalist, err
}

func (stm *MetaStateMachine) DelMetadata(userid int64, bucketName, name string, version int) error {
	_, err := stm.GetMetadata(userid, bucketName, name, version)
	if err != nil {
		log.Fatalln("deltet failed,can't find this meta")
		return err
	}
	delkey := META + strconv.Itoa(int(userid)) + "_" + bucketName + "_" + name + "_" + strconv.Itoa(version)
	err = stm.engine.Delete(delkey)
	if err != nil {
		return err
	}
	return nil
}
