package common

const BucketsNum = 15

const (
	ErrCodeNoErr int64 = iota
	ErrCodeWrongGroup
	ErrCodeWrongLeader
	ErrCodeExecTimeout
)
