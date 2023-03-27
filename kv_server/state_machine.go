package kv_server

import "errors"

type StateMachine interface {
	Get(key string) (string, error)
	Put(key, value string) error
	Append(key, value string) error
}

type MemKv struct {
	KV map[string]string
}

func NewMenKv() *MemKv {
	return &MemKv{
		make(map[string]string),
	}
}

func (kv *MemKv) Get(k string) (string, error) {
	if v, ok := kv.KV[k]; ok {
		return v, nil
	}
	return "", errors.New("can't find the value by this key")
}

func (kv *MemKv) Put(key string, value string) error {
	kv.KV[key] = value
	return nil
}

func (kv *MemKv) Append(key string, value string) error {
	kv.KV[key] += value
	return nil
}
