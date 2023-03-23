package tinnraft

import (
	"fmt"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraftpb"
	"strings"
	"sync"
)

type Log struct {
	mu        sync.RWMutex
	entries   []tinnraftpb.Entry //日志条目列表
	indexHead int
	indexLast int
	engine    storage_engine.KvStorage //存储引擎
}

/*
持久化日志 接口
*/
type LogOperation interface {
	Append(entry *tinnraftpb.Entry)
	GetEntryByidx(idx int64) *tinnraftpb.Entry
	GetLastEntry() *tinnraftpb.Entry
}

/*
内存存储日志
*/

// 初始化日志
func makeEmptyLog() Log {
	log := Log{
		entries:   make([]tinnraftpb.Entry, 0),
		indexHead: 0,
		indexLast: 0,
	}
	return log
}

// 添加日志
func (l *Log) append(entries ...tinnraftpb.Entry) {
	l.entries = append(l.entries, entries...)
}

func (l *Log) append2(entries []*tinnraftpb.Entry) {

	for i := 0; i < len(entries); i++ {
		l.entries = append(l.entries, *entries[i])
	}

}

// 查询日志字段
func (l *Log) at(idx int) *tinnraftpb.Entry {
	return &l.entries[idx]
}

// 清空日志
func (l *Log) truncate(idx int) {
	l.entries = l.entries[:idx]
}

// 获取日志切片
func (l *Log) slice(idx int) []tinnraftpb.Entry {
	return l.entries[idx:]
}

func (l *Log) slice2(idx int) []*tinnraftpb.Entry {
	entries := l.entries[idx:]
	ents := []*tinnraftpb.Entry{}
	for i := 0; i < len(entries); i++ {
		ents = append(ents, &entries[i])
	}
	return ents
}

// 获取日志长度
func (l *Log) len() int {
	return len(l.entries)
}

// 获取最后一条日志
func (l *Log) lastLog() *tinnraftpb.Entry {
	return l.at(l.len() - 1)
}

// 获取单条日志的term
// func ()String() string {
// 	e := make(tinnraftpb.Entry,0)
// 	return fmt.Sprint(e.Term)
// }

// 返回所有term
func (l *Log) String() string {
	nums := []string{}
	for _, entry := range l.entries {
		//将每条日志的所有term转化为字符串写入nums数组
		nums = append(nums, fmt.Sprintf("%4d", entry.Term))
	}
	//将num每个字段用 | 连接并返回为字符串
	return fmt.Sprint(strings.Join(nums, "|"))
}
