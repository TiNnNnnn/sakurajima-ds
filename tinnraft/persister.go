package tinnraft

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraftpb"
)

// 持久化日志
type Persister struct {
	CurrentTerm int64
	VotedFor    int64
}

// 初始化一个持久化日志
func MakePersister(dbEngine storage_engine.KvStorage) *Log {
	empEnt := &tinnraftpb.Entry{}
	empEntEncode := EncodeEntry((empEnt))
	dbEngine.PutBytesKv(EncodeLogKey(INIT_LOG_INDEX), empEntEncode)
	return &Log{
		engine: dbEngine,
	}
}

// 持久化raft状态
func (log *Log) PersistRaftState(curTerm int64, votedFor int64) {
	rfState := &Persister{
		CurrentTerm: curTerm,
		VotedFor:    votedFor,
	}
	log.engine.PutBytesKv(RAFT_STATE_KEY, EncodeRaftState(rfState))
}

// 读取持久化的raft状态
func (log *Log) ReadRaftState() (curTerm int64, votedFor int64) {
	rfBytes, err := log.engine.GetBytesValue(RAFT_STATE_KEY)
	if err != nil {
		return 0, -1
	}
	rfState := DecodeRaftState(rfBytes)
	return rfState.CurrentTerm, rfState.VotedFor
}

// 从存储引擎中读取第一条日志
func (log *Log) GetPersistFirstEntry() *tinnraftpb.Entry {
	log.mu.Lock()
	defer log.mu.Unlock()
	k, v, err := log.engine.GetPerixFirst(string(RAFTLOG_PREFIX))
	if err != nil {
		panic(err)
	}
	DLog("get the first log with id: ", DecodeLogKey(k))
	return DecodeEntry(v)

}

// 获取最新的日志条目
func (log *Log) GetPersistLastEntry() *tinnraftpb.Entry {
	log.mu.Lock()
	defer log.mu.Unlock()
	lastLogId, err := log.engine.GetPrefixKeyIdMax(RAFTLOG_PREFIX)
	if err != nil {
		panic(err)
	}
	firstIdx := log.FirstLogIdx()
	return log.GetEntryWithoutLock(int64(lastLogId) - int64(firstIdx))

}

// 获取日志长度
func (log *Log) LogPersistLen() int {
	log.mu.Lock()
	defer log.mu.Unlock()
	kBytes, _, err := log.engine.GetPerixFirst(string(RAFTLOG_PREFIX))
	if err != nil {
		panic(err)
	}
	logIdFirst := DecodeLogKey(kBytes)
	logIdLast, err := log.engine.GetPrefixKeyIdMax(RAFTLOG_PREFIX)
	if err != nil {
		panic(err)
	}
	return int(logIdLast) - int(logIdFirst) + 1
}

// 切片，获取idx以及之后的所有日志
func (log *Log) TruncatePersistLog(idx int64) []*tinnraftpb.Entry {
	log.mu.Lock()
	defer log.mu.Unlock()
	entries := []*tinnraftpb.Entry{}
	firstLogIdx := log.FirstLogIdx()
	lastLogIdx := log.LastLogIdx()
	for i := int64(firstLogIdx) + idx; i <= int64(lastLogIdx); i++ {
		entries = append(entries, log.GetEntryWithoutLock(i-int64(firstLogIdx)))
	}
	return entries
}

// 切片，获取idx以及之前的所有日志
func (log *Log) SlicePersistLog(idx int64, withDel bool) []*tinnraftpb.Entry {
	log.mu.Lock()
	defer log.mu.Lock()
	firstLogId := log.FirstLogIdx()
	if withDel {
		for i := int64(firstLogId) + idx; i <= int64(log.LastLogIdx()); i++ {
			err := log.engine.DeleteBytesKey(EncodeLogKey(uint64(i)))
			if err != nil {
				panic(err)
			}
		}
	}
	entries := []*tinnraftpb.Entry{}
	for i := firstLogId; i < firstLogId+uint64(idx); i++ {
		entries = append(entries, log.GetEntryWithoutLock(int64(i)-int64(firstLogId)))
	}
	return entries
}

// 获取特定区间的日志
func (log *Log) GetPersistInterLog(lIdx int64, rIdx int64) []*tinnraftpb.Entry {
	log.mu.Lock()
	defer log.mu.Unlock()
	interEntries := []*tinnraftpb.Entry{}
	for i := lIdx; i < rIdx; i++ {
		interEntries = append(interEntries, log.GetEntryWithoutLock(i))
	}
	return interEntries
}

// 添加一条新的日志
func (log *Log) PersistAppend(entry *tinnraftpb.Entry) {
	logIdLast, err := log.engine.GetPrefixKeyIdMax(RAFTLOG_PREFIX)
	if err != nil {
		panic(err)
	}
	log.engine.PutBytesKv(EncodeLogKey(uint64(logIdLast)+1), EncodeEntry(entry))
}

// 通过日志索引获取日志条目
func (log *Log) GetPersistEntryByidx(idx int64) *tinnraftpb.Entry {
	log.mu.Lock()
	defer log.mu.Unlock()
	return log.GetEntryWithoutLock(idx)
}

// 获取日志,工具函数
func (log *Log) GetEntryWithoutLock(offset int64) *tinnraftpb.Entry {
	firstLogidx := log.FirstLogIdx()
	encodeValue, err := log.engine.GetBytesValue(EncodeLogKey(firstLogidx + uint64(offset)))
	if err != nil {
		DLog("get log entry with id %d failed", offset)
	}
	return DecodeEntry(encodeValue)
}

// 获取第一条日志的idx
func (log *Log) FirstLogIdx() uint64 {
	kBytes, _, err := log.engine.GetPerixFirst(string(RAFTLOG_PREFIX))
	if err != nil {
		panic(err)
	}
	return DecodeLogKey(kBytes)
}

// 获取最后一条日志的idx
func (log *Log) LastLogIdx() uint64 {
	maxId, err := log.engine.GetPrefixKeyIdMax(RAFTLOG_PREFIX)
	if err != nil {
		panic(err)
	}
	return maxId
}

// 序列化 持久化raft状态日志的key(RAFTLOG+idx)
func EncodeLogKey(idx uint64) []byte {
	var outBuf bytes.Buffer
	outBuf.Write(RAFTLOG_PREFIX)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(idx))
	outBuf.Write(b)
	return outBuf.Bytes()
}

// 获取raft状态日志的idx (4个bytes之后是idx)
func DecodeLogKey(bts []byte) uint64 {
	return binary.LittleEndian.Uint64(bts[4:])
}

// 日志条目序列化
func EncodeEntry(ent *tinnraftpb.Entry) []byte {
	var bytesEnt bytes.Buffer
	enc := gob.NewEncoder(&bytesEnt)
	enc.Encode(ent)
	return bytesEnt.Bytes()
}

// 日志条目反序列化
func DecodeEntry(in []byte) *tinnraftpb.Entry {
	dec := gob.NewDecoder(bytes.NewBuffer(in))
	ent := tinnraftpb.Entry{}
	dec.Decode(&ent)
	return &ent
}

// 持久化raft结构体 序列化
func EncodeRaftState(rfState *Persister) []byte {
	var bytesState bytes.Buffer
	enc := gob.NewEncoder(&bytesState)
	enc.Encode(rfState)
	return bytesState.Bytes()
}

// 持久化raft状态结构体 反序列化
func DecodeRaftState(in []byte) *Persister {
	dec := gob.NewDecoder(bytes.NewBuffer(in))
	rfState := Persister{}
	dec.Decode(&rfState)
	return &rfState
}

// import "sync"

// type Persister struct {
// 	mu        sync.Mutex
// 	raftstate []byte
// 	snapshot  []byte
// }

// func MakePersister() *Persister {
// 	return &Persister{}
// }

// func clone(orig []byte) []byte {
// 	x := make([]byte, len(orig))
// 	copy(x, orig)
// 	return x
// }

// func (ps *Persister) Copy() *Persister {
// 	ps.mu.Lock()
// 	defer ps.mu.Unlock()
// 	np := MakePersister()
// 	np.raftstate = ps.raftstate
// 	np.snapshot = ps.snapshot
// 	return np
// }

// func (ps *Persister) SaveRaftState(state []byte) {
// 	ps.mu.Lock()
// 	defer ps.mu.Unlock()
// 	ps.raftstate = clone(state)
// }

// func (ps *Persister) ReadRaftState() []byte {
// 	ps.mu.Lock()
// 	defer ps.mu.Unlock()
// 	return clone(ps.raftstate)
// }

// func (ps *Persister) RaftStateSize() int {
// 	ps.mu.Lock()
// 	defer ps.mu.Unlock()
// 	return len(ps.raftstate)
// }

// // Save both Raft state and K/V snapshot as a single atomic action,
// // to help avoid them getting out of sync.
// func (ps *Persister) SaveStateAndSnapshot(state []byte, snapshot []byte) {
// 	ps.mu.Lock()
// 	defer ps.mu.Unlock()
// 	ps.raftstate = clone(state)
// 	ps.snapshot = clone(snapshot)
// }

// func (ps *Persister) ReadSnapshot() []byte {
// 	ps.mu.Lock()
// 	defer ps.mu.Unlock()
// 	return clone(ps.snapshot)
// }

// func (ps *Persister) SnapshotSize() int {
// 	ps.mu.Lock()
// 	defer ps.mu.Unlock()
// 	return len(ps.snapshot)
// }
