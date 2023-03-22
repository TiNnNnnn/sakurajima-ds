package storage_engine

import (
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
func (ldb *LevelDBKvStorage) DeleteByteKey(k []byte) error {
	return ldb.db.Delete(k, nil)
}
