package shared_server

import (
	"sakurajima-ds/storage_engine"
	"strconv"
)

// bucket状态
type bucketStatus uint8

const (
	Runing bucketStatus = iota
	Stopped
	Migrating
	Importing
)

type Bucket struct {
	ID     int
	KvDb   storage_engine.KvStorage
	Status bucketStatus
}

// 创建一个新的桶
func MakeNewBucket(engine storage_engine.KvStorage, id int) *Bucket {
	return &Bucket{id, engine, Runing}
}

func (b *Bucket) Put(key string, value string) error {
	ekey := strconv.Itoa(b.ID) + "_" + key
	return b.KvDb.Put(ekey, value)
}

func (b *Bucket) Get(key string) (string, error) {
	ekey := strconv.Itoa(b.ID) + "_" + key
	value, err := b.KvDb.Get(ekey)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (b *Bucket) Append(key string, value string) error {
	oldV, err := b.Get(key)
	if err != nil {
		return err
	}
	return b.Put(key, oldV+value)
}

func (b *Bucket) DeepCopy() (map[string]string, error) {
	ekeyprefix := strconv.Itoa(b.ID) + "_"
	kvs, err := b.KvDb.GetAllPrefixKey(ekeyprefix)
	if err != nil {
		return nil, err
	}
	return kvs, nil
}

func (b *Bucket) DelBucketData() error {
	ekeyprefex := strconv.Itoa(b.ID) + "_"
	return b.KvDb.DeltePrefixKeys(ekeyprefex)
}
