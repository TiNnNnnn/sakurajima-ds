package storage_engine

/*
存储引擎，接口类
*/

type KvStorage interface {
	Put(string, string) error
	Get(string) (string, error)
	Delete(string) error
}

func EngineFactory(name string, dbPath string) KvStorage {
	switch name {
	case "leveldb":
		levelDb, err := MakeLevelDBKvStorage(dbPath)
		if err != nil {
			panic(err)
		}
		return levelDb
	}
}
