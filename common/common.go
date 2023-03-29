package common

const (
	ErrCodeNoErr int64 = iota
	ErrCodeWrongGroup
	ErrCodeWrongLeader
	ErrCodeExecTimeout
)
