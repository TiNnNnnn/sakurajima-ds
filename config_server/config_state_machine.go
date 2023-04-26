package config_server

import (
	"encoding/json"
	"sakurajima-ds/common"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraft"
	"strconv"
)

const CUR_VERSION_KEY = "CUR_VERISON_KEY"

type ConfigOp interface {
	Join(groups map[int][]string) error
	Leave(groupId []int) error
	Move(bucketId int, groupId int) error
	Query(version int) (Config, error)
}

type ConfigStateMachine struct {
	//mu         sync.Mutex
	engine     storage_engine.KvStorage
	curVerison int
}

func MakeConfigStm(engine storage_engine.KvStorage) *ConfigStateMachine {
	newStm := &ConfigStateMachine{
		engine:     engine,
		curVerison: 0,
	}
	_, err := engine.Get("cf_" + strconv.Itoa(0))
	if err != nil {
		newconfig := MakeDefaultConfig()
		newconfigBytes, err := json.Marshal(newconfig)
		if err != nil {
			panic(err)
		}
		err = newStm.engine.Put("cf_"+strconv.Itoa(0), string(newconfigBytes))
		if err != nil {
			panic(err)
		}

		err = newStm.engine.Put(CUR_VERSION_KEY, strconv.Itoa(newStm.curVerison))
		if err != nil {
			panic(err)
		}
		return newStm
	}
	verStr, _ := engine.Get(CUR_VERSION_KEY)
	version, _ := strconv.Atoi(verStr)
	newStm.curVerison = version
	return newStm
}

func (stm *ConfigStateMachine) Join(groups map[int][]string) error {
	confBytes, err := stm.engine.Get("cf_" + strconv.Itoa(stm.curVerison))
	if err != nil {
		return err
	}
	//读取最新config到lastConfig
	lastConf := &Config{}
	json.Unmarshal([]byte(confBytes), lastConf)

	//构建最新的Config
	newConfig := Config{
		stm.curVerison + 1,
		lastConf.Buckets,
		CopyGroup(lastConf.Groups),
		lastConf.LeaderId,
	}

	//更新Groups,向newConfig添加新的组别
	for groupId, addrs := range groups {
		_, ok := newConfig.Groups[groupId]
		if !ok {
			newAddrs := make([]string, len(addrs))
			copy(newAddrs, addrs)
			newConfig.Groups[groupId] = newAddrs
		}
	}

	//更新Buckets
	tmp := make(map[int][]int)
	for groupId := range newConfig.Groups {
		tmp[groupId] = make([]int, 0)
	}

	for bucketId, groupId := range lastConf.Buckets {
		tmp[groupId] = append(tmp[groupId], bucketId)
	}

	var newBuckets [common.BucketsNum]int
	for groupId, buckets := range tmp {
		for _, bucketId := range buckets {
			newBuckets[bucketId] = groupId
		}
	}
	//持久化最新的config
	newConfig.Buckets = newBuckets
	newConfigBytes, _ := json.Marshal(newConfig)
	stm.engine.Put(CUR_VERSION_KEY, strconv.Itoa(stm.curVerison+1))
	stm.engine.Put("cf_"+strconv.Itoa(stm.curVerison+1), string(newConfigBytes))
	stm.curVerison++
	return nil
}

func (stm *ConfigStateMachine) Leave(groupIds []int) error {
	configBytes, err := stm.engine.Get("cf_" + strconv.Itoa(stm.curVerison))
	if err != nil {
		return nil
	}
	lastConf := &Config{}
	json.Unmarshal([]byte(configBytes), lastConf)
	newConf := Config{
		stm.curVerison + 1,
		lastConf.Buckets,
		CopyGroup(lastConf.Groups),
		lastConf.LeaderId,
	}
	//删除指定的分组
	for _, groupId := range groupIds {
		delete(newConf.Groups, groupId)
	}

	var newBuckets [common.BucketsNum]int
	newConf.Buckets = newBuckets
	newConfigBytes, _ := json.Marshal(newConf)
	//持久化最新Config
	stm.engine.Put(CUR_VERSION_KEY, strconv.Itoa(stm.curVerison))
	stm.engine.Put("cf_"+strconv.Itoa(stm.curVerison+1), string(newConfigBytes))
	stm.curVerison += 1
	return nil
}

// 将bucketid号桶挂到groupId分组服务下
func (stm *ConfigStateMachine) Move(bucketId int, groupId int) error {
	confBytes, err := stm.engine.Get("cf_" + strconv.Itoa(stm.curVerison))
	if err != nil {
		return err
	}
	lastConf := &Config{}
	json.Unmarshal([]byte(confBytes), lastConf)
	newConf := Config{
		stm.curVerison + 1,
		lastConf.Buckets,
		CopyGroup(lastConf.Groups),
		lastConf.LeaderId,
	}
	//bucketId号桶 对于到 groupId号存储集群
	newConf.Buckets[bucketId] = groupId
	//持久化新版本Config
	newConfigBytes, _ := json.Marshal(newConf)
	stm.engine.Put(CUR_VERSION_KEY, strconv.Itoa(stm.curVerison+1))
	stm.engine.Put("cf_"+strconv.Itoa(stm.curVerison+1), string(newConfigBytes))
	stm.curVerison++
	return nil
}

// 查询集群配置
func (stm *ConfigStateMachine) Query(version int) (Config, error) {
	if version < 0 || version >= stm.curVerison {
		//返回最新版本的配置
		lastConf := &Config{}
		tinnraft.DLog("query latest config version: %v", strconv.Itoa(stm.curVerison))
		confBytes, err := stm.engine.Get("cf_" + strconv.Itoa(stm.curVerison))
		if err != nil {
			return MakeDefaultConfig(), err
		}
		json.Unmarshal([]byte(confBytes), lastConf)
		return *lastConf, nil
	}
	tinnraft.DLog("query former config version: %v", strconv.Itoa(version))
	confBytes, err := stm.engine.Get("cf_" + strconv.Itoa(version))
	if err != nil {
		return MakeDefaultConfig(), err
	}
	ans := &Config{}
	json.Unmarshal([]byte(confBytes), ans)
	return *ans, nil
}
