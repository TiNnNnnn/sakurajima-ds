package storage_engine

/*
存储引擎，接口类
*/
type KvStorage interface {
	//put
	Put(string, string) error
	PutBytesKv(k []byte, v []byte) error
	//get
	Get(string) (string, error)
	GetBytesValue(k []byte) ([]byte, error)
	//delete
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
	return nil
}
