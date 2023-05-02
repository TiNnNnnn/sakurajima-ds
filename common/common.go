package common

import "hash/crc32"

const BucketsNum = 15

const (
	ErrCodeNoErr int64 = iota
	ErrCodeWrongGroup
	ErrCodeWrongLeader
	ErrCodeExecTimeout
	ErrCodeNotReady
	ErrCodeCopyBuckets
)

func KeyToBucketId(key string) int {
	return CRC32KeyHash(key, BucketsNum)
}

func CRC32KeyHash(k string, base int) int {
	bucketId := 0
	crc32q := crc32.MakeTable(0xD5828281)
	sum := crc32.Checksum([]byte(k), crc32q)
	bucketId = int(sum) % BucketsNum
	return bucketId
}

func Int64ArrToIntArr(in []int64) []int {
	out := []int{}
	for _, item := range in {
		out = append(out, int(item))
	}
	return out
}
