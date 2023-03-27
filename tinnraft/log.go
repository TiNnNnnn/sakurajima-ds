package tinnraft

import (
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraftpb"
	"sync"
)

type Log struct {
	mu         sync.RWMutex
	entries    []*tinnraftpb.Entry //日志条目列表
	indexHead  int64
	indexLast  int64
	appliedIdx int64
	engine     storage_engine.KvStorage //存储引擎
}

/*
持久化日志 接口
*/
type LogOperation interface {
	GetPersistFirstEntry() *tinnraftpb.Entry
	GetPersistLastEntry() *tinnraftpb.Entry

	LogPersistLen() int

	TruncatePersistLog(idx int64) []*tinnraftpb.Entry
	SlicePersistLog(idx int64, withDel bool) []*tinnraftpb.Entry

	PersistAppend(entry *tinnraftpb.Entry)

	GetPersistInterLog(lIdx int64, rIdx int64) []*tinnraftpb.Entry
	GetPersistEntryByidx(idx int64) *tinnraftpb.Entry
}

/*
内存存储日志,已经弃用
*/

/*
// 初始化一份日志,并添加一条空日志
func makeEmptyLog() *Log {
	empEnt := &tinnraftpb.Entry{}
	newEntries := []*tinnraftpb.Entry{}
	newEntries = append(newEntries, empEnt)
	return &Log{
		entries:   newEntries,
		indexHead: 0,
		indexLast: 1,
	}
}

// 添加一条新的日志
func (log *Log) AppendLog(entry *tinnraftpb.Entry) {
	log.entries = append(log.entries, entry)
}

// 根据条目小标查找日志
func (l *Log) at(idx int) *tinnraftpb.Entry {
	return l.entries[idx]
}

// 获取最后一条日志
func (l *Log) GetlastLog() *tinnraftpb.Entry {
	return l.at(l.LenLog() - 1)
}

// 获取第一条日志
func (l *Log) GetfirstLog() *tinnraftpb.Entry {
	return l.at(0)
}

// 获取特定区间日志
func (l *Log) GetInterLog(lIdx int, rIdx int) []*tinnraftpb.Entry {
	return l.entries[lIdx:rIdx]
}

// 清空idx之前的所有日志
func (l *Log) TruncateLog(idx int) {
	l.entries = l.entries[:idx]
}

// 获取从idx开始的所有日志
func (l *Log) SliceLog(idx int) {
	l.entries = l.entries[idx:]
}

// 获取日志长度
func (l *Log) LenLog() int {
	return len(l.entries)
}

//old_verion------------------------------------

func (l *Log) append2(entries []*tinnraftpb.Entry) {

	for i := 0; i < len(entries); i++ {
		l.entries = append(l.entries, entries[i])
	}
}

// 获取日志切片
func (l *Log) slice2(idx int) []*tinnraftpb.Entry {
	entries := l.entries[idx:]
	ents := []*tinnraftpb.Entry{}
	for i := 0; i < len(entries); i++ {
		ents = append(ents, entries[i])
	}
	return ents
}

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
*/
