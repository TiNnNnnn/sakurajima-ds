package tinnraft

import "sakurajima-ds/tinnraftpb"

/*
message InstallSnapshotArgs{
	int64 Term = 1;
	int64 LeaderId = 2;
	int64 LastIncludedIndex = 3; //快照中包含的最后日志条目的索引
	int64 LastIncludeTerm = 4;
	bytes Data = 5;
}

message InstallSnapshotReply{
	int64 Term = 1;
}
*/
// A service wants to switch to snapshot.  Only do so if Raft hasn't
// have more recent info since it communicate the snapshot on applyCh.

// 从Leader下载日志快照
func (rf *Raft) HandleInstallSnapshot(args *tinnraftpb.InstallSnapshotArgs, reply *tinnraftpb.InstallSnapshotReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	reply.Term = int64(rf.currentTerm)

	if args.Term < int64(rf.currentTerm) {
		return
	}

	if args.Term > int64(rf.currentTerm) {
		rf.setNewTerm(int(args.Term))
	}

	rf.state = Follower
	rf.resetElectionTimer()

	if args.LastIncludedIndex <= int64(rf.commitIndex) {
		return
	}

	go func() {
		rf.applyCh <- &tinnraftpb.ApplyMsg{
			SnapshotValid: true,
			Snapshot:      args.Data,
			SnapshotTerm:  args.LastIncludeTerm,
			SnapshotIndex: args.LastIncludedIndex,
		}
	}()

}

// 下载一份日志快照
func (rf *Raft) CondInstallSnapshot(lastIncludedTerm int, lastIncludedIndex int, snapshot []byte) bool {

	rf.mu.Lock()
	defer rf.mu.Unlock()
	if lastIncludedIndex <= rf.commitIndex {
		return false
	}

	if lastIncludedIndex > int(rf.log.GetlastLog().Index) {
		rf.log.ReInitPersistLog()
	} else {
		//删除日志
		rf.log.TruncatePersistLogWithDel(int64(lastIncludedIndex) - rf.log.GetfirstLog().Index)
		//将第一个操作日志设置为空
		rf.log.SetPersistFirstData([]byte{})
	}
	rf.log.SetPersistFirstEntryTermAndIndex(int64(lastIncludedTerm), int64(lastIncludedIndex))

	rf.lastApplied = lastIncludedIndex
	rf.commitIndex = lastIncludedIndex

	return true
}

// 进行日志快照
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.isSnapshoting = true
	snapshotIndex := rf.log.FirstLogIdx()
	if index < int(snapshotIndex) {
		rf.isSnapshoting = false
		return
	}
	//删除日志
	rf.log.TruncatePersistLogWithDel(int64(index) - int64(snapshotIndex))
	//将第一个操作日志设置为空
	rf.log.SetPersistFirstData([]byte{})
	rf.isSnapshoting = false
	rf.log.PersistSnapshot(snapshot)
}

// 读取持久化的快照信息
func (rf *Raft) ReadSnapshot() []byte {
	bytes, err := rf.log.ReadSnapshot()
	if err != nil {
		DLog("read snapshot error:%v", err.Error())
	}
	return bytes
}
