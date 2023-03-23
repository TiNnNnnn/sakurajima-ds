package tinnraft

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"sakurajima-ds/storage_engine"
	"sakurajima-ds/tinnraftpb"
)

type Persister struct {
	CurrentTerm int64
	VotedFor    int64
}

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

func EncodeLogKey(idx uint64) []byte {
	var outBuf bytes.Buffer
	outBuf.Write(RAFTLOG_PREFIX)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(idx))
	outBuf.Write(b)
	return outBuf.Bytes()
}

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
