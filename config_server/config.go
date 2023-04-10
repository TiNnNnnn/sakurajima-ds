package config_server

/*
	记录数据服务器的分组情况
*/

import (
	"sakurajima-ds/common"
)

type Config struct {
	Version int
	Buckets [common.BucketsNum]int
	Groups  map[int][]string
}
