package tinnraft

// A service wants to switch to snapshot.  Only do so if Raft hasn't
// have more recent info since it communicate the snapshot on applyCh.
func (rf *Raft) CondInstallSnapshot(lastIncludedTerm int, lastIncludedIndex int, snapshot []byte) bool {

	rf.mu.Lock()
	defer rf.mu.Unlock()

	return true
}

// 进行日志快照
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (2D).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.isSnapshoting = true
	snapshotIndex := rf.log.FirstLogIdx()
	if index < int(snapshotIndex) {
		rf.isSnapshoting = false
		return
	}
	rf.log.TruncatePersistLogWithDel(int64(index) - int64(snapshotIndex))
	rf.log.SetPersistFirstData([]byte{})
	rf.isSnapshoting = false
	rf.log.PersistSnapshot(snapshot)
}
