package storage_engine

import (
	"encoding/binary"
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDBKvStorage struct {
	Path string
	db   *leveldb.DB
}

func MakeLevelDBKvStorage(path string) (*LevelDBKvStorage, error) {
	//创建leveldb连接，并且在当前目录下创建一个文件夹
	newdb, err := leveldb.OpenFile(path, &opt.Options{})
	if err != nil {
		return nil, err
	}
	return &LevelDBKvStorage{
		Path: path,
		db:   newdb,
	}, nil
}

// 写入kv
func (ldb *LevelDBKvStorage) Put(k string, v string) error {
	return ldb.db.Put([]byte(k), []byte(v), nil)
}

// 以字节写入kv
func (ldb *LevelDBKvStorage) PutBytesKv(k []byte, v []byte) error {
	return ldb.db.Put(k, v, nil)
}

// 按照key查找返回value
func (ldb *LevelDBKvStorage) Get(k string) (string, error) {
	//按k查找v,没有则返回nil
	v, err := ldb.db.Get([]byte(k), nil)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

// 以字节读取value
func (ldb *LevelDBKvStorage) GetBytesValue(k []byte) ([]byte, error) {
	return ldb.db.Get(k, nil)
}

// 根据key删除kv
func (ldb *LevelDBKvStorage) Delete(k string) error {
	return ldb.db.Delete([]byte(k), nil)
}

// 以字节根据key删除kv
func (ldb *LevelDBKvStorage) DeleteBytesKey(k []byte) error {
	return ldb.db.Delete(k, nil)
}

// 查找所有以pre为前缀的key值所对应的kv
func (ldb *LevelDBKvStorage) GetAllPrefixKey(pre string) (map[string]string, error) {
	kvs := make(map[string]string)
	iter := ldb.db.NewIterator(util.BytesPrefix([]byte(pre)), nil)
	defer iter.Release()
	for iter.Next() {
		k := string(iter.Key())
		v := string(iter.Value())
		kvs[k] = v
	}
	return kvs, iter.Error()
}

// 查询最新版本的以pre为前缀的key值所对应的kv
func (ldb *LevelDBKvStorage) GetPrefixLast(pre []byte) ([]byte, []byte, error) {
	iter := ldb.db.NewIterator(util.BytesPrefix(pre), nil)
	defer iter.Release()
	ok := iter.Last()
	if ok {
		return iter.Key(), iter.Value(), nil
	}
	return []byte{}, []byte{}, nil
}

// 查询最旧版本的以pre为前缀的key值所对应的kv
func (ldb *LevelDBKvStorage) GetPerixFirst(pre string) ([]byte, []byte, error) {
	iter := ldb.db.NewIterator(util.BytesPrefix([]byte(pre)), nil)
	defer iter.Release()
	if iter.Next() {
		return iter.Key(), iter.Value(), nil
	}
	return []byte{}, []byte{}, errors.New("no key with a prefix " + string(pre))
}

// 获取以pre为前缀的日志中最大的key
func (ldb *LevelDBKvStorage) GetPrefixKeyIdMax(pre []byte) (uint64, error) {
	iter := ldb.db.NewIterator(util.BytesPrefix([]byte(pre)), nil)
	defer iter.Release()
	var maxKeyId uint64 = 0
	for iter.Next() {
		if iter.Error() != nil {
			return maxKeyId, iter.Error()
		}
		kBytes := iter.Key()
		KeyId := binary.LittleEndian.Uint64(kBytes[len(pre):])
		if KeyId > maxKeyId {
			maxKeyId = KeyId
		}
	}
	return maxKeyId, nil
}

// 删除所有以pre为前缀的key值所对应的kv
func (ldb *LevelDBKvStorage) DeltePrefixKeys(pre string) error {
	iter := ldb.db.NewIterator(util.BytesPrefix([]byte(pre)), nil)
	defer iter.Release()
	for iter.Next() {
		err := ldb.db.Delete(iter.Key(), nil)
		if err != nil {
			return err
		}
	}
	return nil
}
