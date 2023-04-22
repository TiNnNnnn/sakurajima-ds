package api_gateway

import (
	"encoding/json"
	"log"
	"sakurajima-ds/storage_engine"
	"strconv"
)

const CUR_VERSION_KEY = "CUR_VERISON_KEY"

type AddrConfigOp interface {
	Update(AddrConfig) error
	Query() (AddrConfig, error)
}

type AddrStateMachine struct {
	engine     storage_engine.KvStorage
	curVersion int
}

func MakeAddrConfigStm(engine storage_engine.KvStorage) *AddrStateMachine {
	newStm := &AddrStateMachine{
		engine:     engine,
		curVersion: 0,
	}

	_, err := engine.Get("acf_" + strconv.Itoa(0))
	if err != nil {
		newconfig := MakeDefaultConfig()
		newconfigBytes, err := json.Marshal(newconfig)
		if err != nil {
			panic(err)
		}
		err = newStm.engine.Put("acf_"+strconv.Itoa(0), string(newconfigBytes))
		if err != nil {
			panic(err)
		}

		err = newStm.engine.Put(CUR_VERSION_KEY, strconv.Itoa(newStm.curVersion))
		if err != nil {
			panic(err)
		}
		return newStm
	}
	verStr, _ := engine.Get(CUR_VERSION_KEY)
	version, _ := strconv.Atoi(verStr)
	newStm.curVersion = version
	return newStm
}

func (stm *AddrStateMachine) Query(version int) (AddrConfig, error) {
	if version < 0 || version >= stm.curVersion {
		//返回最新版本的配置
		lastConf := &AddrConfig{}
		log.Printf("query latest addrconfig version: %v", strconv.Itoa(stm.curVersion))
		confBytes, err := stm.engine.Get("acf_" + strconv.Itoa(stm.curVersion))
		if err != nil {
			return MakeDefaultConfig(), err
		}
		json.Unmarshal([]byte(confBytes), lastConf)
		return *lastConf, nil
	}
	//查找历史版本配置
	log.Printf("query former config version: %v", strconv.Itoa(version))
	confBytes, err := stm.engine.Get("acf_" + strconv.Itoa(version))
	if err != nil {
		return MakeDefaultConfig(), err
	}
	ans := &AddrConfig{}
	json.Unmarshal([]byte(confBytes), ans)
	return *ans, nil
}

// 更新配置信息
func (stm *AddrStateMachine) Update(conf AddrConfig) error {
	confBytes, err := stm.engine.Get("acf_" + strconv.Itoa(stm.curVersion))
	if err != nil {
		return err
	}
	//读取最新config到lastConfig
	lastConf := &AddrConfig{}
	json.Unmarshal([]byte(confBytes), lastConf)

	//构建最新的Config
	newConfig := AddrConfig{
		conf.Cfg_server_addr,
		conf.Shared_server_addr,
		int64(stm.curVersion + 1),
	}
	newConfigBytes, _ := json.Marshal(newConfig)
	stm.engine.Put(CUR_VERSION_KEY, strconv.Itoa(stm.curVersion+1))
	stm.engine.Put("acf_"+strconv.Itoa(stm.curVersion+1), string(newConfigBytes))
	stm.curVersion++
	return nil

}
