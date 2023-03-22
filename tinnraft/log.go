package tinnraft

import (
	"fmt"
	"sakurajima-ds/tinnraftpb"
	"strings"
)

type Log struct {
	Entries   []tinnraftpb.Entry
	IndexHead int
}

// 初始化日志
func makeEmptyLog() Log {
	log := Log{
		Entries:   make([]tinnraftpb.Entry, 0),
		IndexHead: 0,
	}
	return log
}

// 添加日志
func (l *Log) append(entries ...tinnraftpb.Entry) {
	l.Entries = append(l.Entries, entries...)
}

// 查询日志字段
func (l *Log) at(idx int) *tinnraftpb.Entry {
	return &l.Entries[idx]
}

// 清空日志
func (l *Log) truncate(idx int) {
	l.Entries = l.Entries[:idx]
}

// 获取日志切片
func (l *Log) slice(idx int) []tinnraftpb.Entry {
	return l.Entries[idx:]
}

// 获取日志长度
func (l *Log) len() int {
	return len(l.Entries)
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
	for _, entry := range l.Entries {
		//将每条日志的所有term转化为字符串写入nums数组
		nums = append(nums, fmt.Sprintf("%4d", entry.Term))
	}
	//将num每个字段用 | 连接并返回为字符串
	return fmt.Sprint(strings.Join(nums, "|"))
}
