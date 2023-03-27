// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: tinnraft.proto

package tinnraftpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EntryType int32

const (
	EntryType_EntryData       EntryType = 0
	EntryType_EntryConfUpdate EntryType = 1
)

// Enum value maps for EntryType.
var (
	EntryType_name = map[int32]string{
		0: "EntryData",
		1: "EntryConfUpdate",
	}
	EntryType_value = map[string]int32{
		"EntryData":       0,
		"EntryConfUpdate": 1,
	}
)

func (x EntryType) Enum() *EntryType {
	p := new(EntryType)
	*p = x
	return p
}

func (x EntryType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EntryType) Descriptor() protoreflect.EnumDescriptor {
	return file_tinnraft_proto_enumTypes[0].Descriptor()
}

func (EntryType) Type() protoreflect.EnumType {
	return &file_tinnraft_proto_enumTypes[0]
}

func (x EntryType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EntryType.Descriptor instead.
func (EntryType) EnumDescriptor() ([]byte, []int) {
	return file_tinnraft_proto_rawDescGZIP(), []int{0}
}

type OpType int32

const (
	OpType_Put           OpType = 0
	OpType_Append        OpType = 1
	OpType_Get           OpType = 2
	OpType_ConfigChange  OpType = 3
	OpType_DeleteBuckets OpType = 4
	OpType_InsertBuckets OpType = 5
)

// Enum value maps for OpType.
var (
	OpType_name = map[int32]string{
		0: "Put",
		1: "Append",
		2: "Get",
		3: "ConfigChange",
		4: "DeleteBuckets",
		5: "InsertBuckets",
	}
	OpType_value = map[string]int32{
		"Put":           0,
		"Append":        1,
		"Get":           2,
		"ConfigChange":  3,
		"DeleteBuckets": 4,
		"InsertBuckets": 5,
	}
)

func (x OpType) Enum() *OpType {
	p := new(OpType)
	*p = x
	return p
}

func (x OpType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OpType) Descriptor() protoreflect.EnumDescriptor {
	return file_tinnraft_proto_enumTypes[1].Descriptor()
}

func (OpType) Type() protoreflect.EnumType {
	return &file_tinnraft_proto_enumTypes[1]
}

func (x OpType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OpType.Descriptor instead.
func (OpType) EnumDescriptor() ([]byte, []int) {
	return file_tinnraft_proto_rawDescGZIP(), []int{1}
}

type RequestVoteArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         int64 `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	CandidateId  int64 `protobuf:"varint,2,opt,name=CandidateId,proto3" json:"CandidateId,omitempty"`
	LastLogIndex int64 `protobuf:"varint,3,opt,name=LastLogIndex,proto3" json:"LastLogIndex,omitempty"`
	LastLogTerm  int64 `protobuf:"varint,4,opt,name=LastLogTerm,proto3" json:"LastLogTerm,omitempty"`
}

func (x *RequestVoteArgs) Reset() {
	*x = RequestVoteArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tinnraft_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestVoteArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVoteArgs) ProtoMessage() {}

func (x *RequestVoteArgs) ProtoReflect() protoreflect.Message {
	mi := &file_tinnraft_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVoteArgs.ProtoReflect.Descriptor instead.
func (*RequestVoteArgs) Descriptor() ([]byte, []int) {
	return file_tinnraft_proto_rawDescGZIP(), []int{0}
}

func (x *RequestVoteArgs) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *RequestVoteArgs) GetCandidateId() int64 {
	if x != nil {
		return x.CandidateId
	}
	return 0
}

func (x *RequestVoteArgs) GetLastLogIndex() int64 {
	if x != nil {
		return x.LastLogIndex
	}
	return 0
}

func (x *RequestVoteArgs) GetLastLogTerm() int64 {
	if x != nil {
		return x.LastLogTerm
	}
	return 0
}

type RequestVoteReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term        int64 `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	VoteGranted bool  `protobuf:"varint,2,opt,name=VoteGranted,proto3" json:"VoteGranted,omitempty"`
}

func (x *RequestVoteReply) Reset() {
	*x = RequestVoteReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tinnraft_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestVoteReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVoteReply) ProtoMessage() {}

func (x *RequestVoteReply) ProtoReflect() protoreflect.Message {
	mi := &file_tinnraft_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVoteReply.ProtoReflect.Descriptor instead.
func (*RequestVoteReply) Descriptor() ([]byte, []int) {
	return file_tinnraft_proto_rawDescGZIP(), []int{1}
}

func (x *RequestVoteReply) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *RequestVoteReply) GetVoteGranted() bool {
	if x != nil {
		return x.VoteGranted
	}
	return false
}

type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  EntryType `protobuf:"varint,1,opt,name=Type,proto3,enum=pbs.EntryType" json:"Type,omitempty"`
	Term  uint64    `protobuf:"varint,2,opt,name=Term,proto3" json:"Term,omitempty"`
	Index int64     `protobuf:"varint,3,opt,name=Index,proto3" json:"Index,omitempty"`
	Data  []byte    `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tinnraft_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_tinnraft_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_tinnraft_proto_rawDescGZIP(), []int{2}
}

func (x *Entry) GetType() EntryType {
	if x != nil {
		return x.Type
	}
	return EntryType_EntryData
}

func (x *Entry) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *Entry) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *Entry) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type AppendEntriesArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         int64    `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	LeaderId     int64    `protobuf:"varint,2,opt,name=LeaderId,proto3" json:"LeaderId,omitempty"`
	PrevLogIndex int64    `protobuf:"varint,3,opt,name=PrevLogIndex,proto3" json:"PrevLogIndex,omitempty"`
	PrevLogTerm  int64    `protobuf:"varint,4,opt,name=PrevLogTerm,proto3" json:"PrevLogTerm,omitempty"`
	Entries      []*Entry `protobuf:"bytes,5,rep,name=Entries,proto3" json:"Entries,omitempty"`
	LeaderCommit int64    `protobuf:"varint,6,opt,name=LeaderCommit,proto3" json:"LeaderCommit,omitempty"`
}

func (x *AppendEntriesArgs) Reset() {
	*x = AppendEntriesArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tinnraft_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntriesArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntriesArgs) ProtoMessage() {}

func (x *AppendEntriesArgs) ProtoReflect() protoreflect.Message {
	mi := &file_tinnraft_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntriesArgs.ProtoReflect.Descriptor instead.
func (*AppendEntriesArgs) Descriptor() ([]byte, []int) {
	return file_tinnraft_proto_rawDescGZIP(), []int{3}
}

func (x *AppendEntriesArgs) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntriesArgs) GetLeaderId() int64 {
	if x != nil {
		return x.LeaderId
	}
	return 0
}

func (x *AppendEntriesArgs) GetPrevLogIndex() int64 {
	if x != nil {
		return x.PrevLogIndex
	}
	return 0
}

func (x *AppendEntriesArgs) GetPrevLogTerm() int64 {
	if x != nil {
		return x.PrevLogTerm
	}
	return 0
}

func (x *AppendEntriesArgs) GetEntries() []*Entry {
	if x != nil {
		return x.Entries
	}
	return nil
}

func (x *AppendEntriesArgs) GetLeaderCommit() int64 {
	if x != nil {
		return x.LeaderCommit
	}
	return 0
}

type AppendEntriesReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term     int64 `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	Success  bool  `protobuf:"varint,2,opt,name=Success,proto3" json:"Success,omitempty"`
	Conflict bool  `protobuf:"varint,3,opt,name=Conflict,proto3" json:"Conflict,omitempty"`
	XTerm    int64 `protobuf:"varint,4,opt,name=XTerm,proto3" json:"XTerm,omitempty"`
	XIndex   int64 `protobuf:"varint,5,opt,name=XIndex,proto3" json:"XIndex,omitempty"`
	XLen     int64 `protobuf:"varint,6,opt,name=XLen,proto3" json:"XLen,omitempty"`
}

func (x *AppendEntriesReply) Reset() {
	*x = AppendEntriesReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tinnraft_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntriesReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntriesReply) ProtoMessage() {}

func (x *AppendEntriesReply) ProtoReflect() protoreflect.Message {
	mi := &file_tinnraft_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntriesReply.ProtoReflect.Descriptor instead.
func (*AppendEntriesReply) Descriptor() ([]byte, []int) {
	return file_tinnraft_proto_rawDescGZIP(), []int{4}
}

func (x *AppendEntriesReply) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntriesReply) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *AppendEntriesReply) GetConflict() bool {
	if x != nil {
		return x.Conflict
	}
	return false
}

func (x *AppendEntriesReply) GetXTerm() int64 {
	if x != nil {
		return x.XTerm
	}
	return 0
}

func (x *AppendEntriesReply) GetXIndex() int64 {
	if x != nil {
		return x.XIndex
	}
	return 0
}

func (x *AppendEntriesReply) GetXLen() int64 {
	if x != nil {
		return x.XLen
	}
	return 0
}

type ApplyMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommandValid  bool   `protobuf:"varint,1,opt,name=CommandValid,proto3" json:"CommandValid,omitempty"`
	Command       []byte `protobuf:"bytes,2,opt,name=Command,proto3" json:"Command,omitempty"`
	CommandIndex  int64  `protobuf:"varint,3,opt,name=CommandIndex,proto3" json:"CommandIndex,omitempty"`
	CommandTerm   int64  `protobuf:"varint,4,opt,name=CommandTerm,proto3" json:"CommandTerm,omitempty"`
	SnapshotValid bool   `protobuf:"varint,5,opt,name=SnapshotValid,proto3" json:"SnapshotValid,omitempty"`
	Snapshot      []byte `protobuf:"bytes,6,opt,name=Snapshot,proto3" json:"Snapshot,omitempty"`
	SnapshotTerm  int64  `protobuf:"varint,7,opt,name=SnapshotTerm,proto3" json:"SnapshotTerm,omitempty"`
	SnapshotIndex int64  `protobuf:"varint,8,opt,name=SnapshotIndex,proto3" json:"SnapshotIndex,omitempty"`
}

func (x *ApplyMsg) Reset() {
	*x = ApplyMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tinnraft_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplyMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplyMsg) ProtoMessage() {}

func (x *ApplyMsg) ProtoReflect() protoreflect.Message {
	mi := &file_tinnraft_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplyMsg.ProtoReflect.Descriptor instead.
func (*ApplyMsg) Descriptor() ([]byte, []int) {
	return file_tinnraft_proto_rawDescGZIP(), []int{5}
}

func (x *ApplyMsg) GetCommandValid() bool {
	if x != nil {
		return x.CommandValid
	}
	return false
}

func (x *ApplyMsg) GetCommand() []byte {
	if x != nil {
		return x.Command
	}
	return nil
}

func (x *ApplyMsg) GetCommandIndex() int64 {
	if x != nil {
		return x.CommandIndex
	}
	return 0
}

func (x *ApplyMsg) GetCommandTerm() int64 {
	if x != nil {
		return x.CommandTerm
	}
	return 0
}

func (x *ApplyMsg) GetSnapshotValid() bool {
	if x != nil {
		return x.SnapshotValid
	}
	return false
}

func (x *ApplyMsg) GetSnapshot() []byte {
	if x != nil {
		return x.Snapshot
	}
	return nil
}

func (x *ApplyMsg) GetSnapshotTerm() int64 {
	if x != nil {
		return x.SnapshotTerm
	}
	return 0
}

func (x *ApplyMsg) GetSnapshotIndex() int64 {
	if x != nil {
		return x.SnapshotIndex
	}
	return 0
}

type InstallSnapshotArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term              int64  `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	LeaderId          int64  `protobuf:"varint,2,opt,name=LeaderId,proto3" json:"LeaderId,omitempty"`
	LastIncludedIndex int64  `protobuf:"varint,3,opt,name=LastIncludedIndex,proto3" json:"LastIncludedIndex,omitempty"`
	LastIncludeTerm   int64  `protobuf:"varint,4,opt,name=LastIncludeTerm,proto3" json:"LastIncludeTerm,omitempty"`
	Data              []byte `protobuf:"bytes,5,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *InstallSnapshotArgs) Reset() {
	*x = InstallSnapshotArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tinnraft_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstallSnapshotArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstallSnapshotArgs) ProtoMessage() {}

func (x *InstallSnapshotArgs) ProtoReflect() protoreflect.Message {
	mi := &file_tinnraft_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstallSnapshotArgs.ProtoReflect.Descriptor instead.
func (*InstallSnapshotArgs) Descriptor() ([]byte, []int) {
	return file_tinnraft_proto_rawDescGZIP(), []int{6}
}

func (x *InstallSnapshotArgs) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *InstallSnapshotArgs) GetLeaderId() int64 {
	if x != nil {
		return x.LeaderId
	}
	return 0
}

func (x *InstallSnapshotArgs) GetLastIncludedIndex() int64 {
	if x != nil {
		return x.LastIncludedIndex
	}
	return 0
}

func (x *InstallSnapshotArgs) GetLastIncludeTerm() int64 {
	if x != nil {
		return x.LastIncludeTerm
	}
	return 0
}

func (x *InstallSnapshotArgs) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type InstallSnapshotReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term int64 `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
}

func (x *InstallSnapshotReply) Reset() {
	*x = InstallSnapshotReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tinnraft_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstallSnapshotReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstallSnapshotReply) ProtoMessage() {}

func (x *InstallSnapshotReply) ProtoReflect() protoreflect.Message {
	mi := &file_tinnraft_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstallSnapshotReply.ProtoReflect.Descriptor instead.
func (*InstallSnapshotReply) Descriptor() ([]byte, []int) {
	return file_tinnraft_proto_rawDescGZIP(), []int{7}
}

func (x *InstallSnapshotReply) GetTerm() int64 {
	if x != nil {
		return x.Term
	}
	return 0
}

type CommandArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key       string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value     string `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
	OpType    OpType `protobuf:"varint,3,opt,name=Op_type,json=OpType,proto3,enum=pbs.OpType" json:"Op_type,omitempty"`
	ClientId  int64  `protobuf:"varint,4,opt,name=Client_id,json=ClientId,proto3" json:"Client_id,omitempty"`
	CommandId int64  `protobuf:"varint,5,opt,name=command_id,json=commandId,proto3" json:"command_id,omitempty"`
	Context   []byte `protobuf:"bytes,6,opt,name=context,proto3" json:"context,omitempty"`
}

func (x *CommandArgs) Reset() {
	*x = CommandArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tinnraft_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommandArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommandArgs) ProtoMessage() {}

func (x *CommandArgs) ProtoReflect() protoreflect.Message {
	mi := &file_tinnraft_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommandArgs.ProtoReflect.Descriptor instead.
func (*CommandArgs) Descriptor() ([]byte, []int) {
	return file_tinnraft_proto_rawDescGZIP(), []int{8}
}

func (x *CommandArgs) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *CommandArgs) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *CommandArgs) GetOpType() OpType {
	if x != nil {
		return x.OpType
	}
	return OpType_Put
}

func (x *CommandArgs) GetClientId() int64 {
	if x != nil {
		return x.ClientId
	}
	return 0
}

func (x *CommandArgs) GetCommandId() int64 {
	if x != nil {
		return x.CommandId
	}
	return 0
}

func (x *CommandArgs) GetContext() []byte {
	if x != nil {
		return x.Context
	}
	return nil
}

type CommandReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value    string `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
	LeaderId int64  `protobuf:"varint,2,opt,name=LeaderId,proto3" json:"LeaderId,omitempty"`
	ErrCode  int64  `protobuf:"varint,3,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
}

func (x *CommandReply) Reset() {
	*x = CommandReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tinnraft_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommandReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommandReply) ProtoMessage() {}

func (x *CommandReply) ProtoReflect() protoreflect.Message {
	mi := &file_tinnraft_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommandReply.ProtoReflect.Descriptor instead.
func (*CommandReply) Descriptor() ([]byte, []int) {
	return file_tinnraft_proto_rawDescGZIP(), []int{9}
}

func (x *CommandReply) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *CommandReply) GetLeaderId() int64 {
	if x != nil {
		return x.LeaderId
	}
	return 0
}

func (x *CommandReply) GetErrCode() int64 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

var File_tinnraft_proto protoreflect.FileDescriptor

var file_tinnraft_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x74, 0x69, 0x6e, 0x6e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x03, 0x70, 0x62, 0x73, 0x22, 0x8d, 0x01, 0x0a, 0x0f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x56, 0x6f, 0x74, 0x65, 0x41, 0x72, 0x67, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72,
	0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x20, 0x0a,
	0x0b, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0b, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49, 0x64, 0x12,
	0x22, 0x0a, 0x0c, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65,
	0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f,
	0x67, 0x54, 0x65, 0x72, 0x6d, 0x22, 0x48, 0x0a, 0x10, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72,
	0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x20, 0x0a,
	0x0b, 0x56, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x22,
	0x69, 0x0a, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x22, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x54, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d,
	0x12, 0x14, 0x0a, 0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x22, 0xd3, 0x01, 0x0a, 0x11, 0x41,
	0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x41, 0x72, 0x67, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04,
	0x54, 0x65, 0x72, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x22, 0x0a, 0x0c, 0x50, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x50, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x72, 0x65, 0x76, 0x4c, 0x6f, 0x67, 0x54,
	0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x50, 0x72, 0x65, 0x76, 0x4c,
	0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x24, 0x0a, 0x07, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65,
	0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x07, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x22, 0x0a, 0x0c,
	0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0c, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x22, 0xa0, 0x01, 0x0a, 0x12, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69,
	0x65, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x53,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x53, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6f, 0x6e, 0x66, 0x6c, 0x69, 0x63,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x43, 0x6f, 0x6e, 0x66, 0x6c, 0x69, 0x63,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x58, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x58, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x58, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x58, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12,
	0x12, 0x0a, 0x04, 0x58, 0x4c, 0x65, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x58,
	0x4c, 0x65, 0x6e, 0x22, 0x9a, 0x02, 0x0a, 0x08, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x4d, 0x73, 0x67,
	0x12, 0x22, 0x0a, 0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x56, 0x61, 0x6c, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x22,
	0x0a, 0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x54, 0x65, 0x72,
	0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x54, 0x65, 0x72, 0x6d, 0x12, 0x24, 0x0a, 0x0d, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x53, 0x6e, 0x61,
	0x70, 0x73, 0x68, 0x6f, 0x74, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x6e,
	0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x53, 0x6e,
	0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68,
	0x6f, 0x74, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x53, 0x6e,
	0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x24, 0x0a, 0x0d, 0x53, 0x6e,
	0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0d, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x22, 0xb1, 0x01, 0x0a, 0x13, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x53, 0x6e, 0x61, 0x70,
	0x73, 0x68, 0x6f, 0x74, 0x41, 0x72, 0x67, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x1a, 0x0a, 0x08,
	0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x11, 0x4c, 0x61, 0x73, 0x74,
	0x49, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x64, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x11, 0x4c, 0x61, 0x73, 0x74, 0x49, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65,
	0x64, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x28, 0x0a, 0x0f, 0x4c, 0x61, 0x73, 0x74, 0x49, 0x6e,
	0x63, 0x6c, 0x75, 0x64, 0x65, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0f, 0x4c, 0x61, 0x73, 0x74, 0x49, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x54, 0x65, 0x72, 0x6d,
	0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x44, 0x61, 0x74, 0x61, 0x22, 0x2a, 0x0a, 0x14, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x53,
	0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04,
	0x54, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d,
	0x22, 0xb1, 0x01, 0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x41, 0x72, 0x67, 0x73,
	0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x24, 0x0a, 0x07, 0x4f, 0x70, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x70, 0x62, 0x73, 0x2e,
	0x4f, 0x70, 0x54, 0x79, 0x70, 0x65, 0x52, 0x06, 0x4f, 0x70, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x78, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x78, 0x74, 0x22, 0x5a, 0x0a, 0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x4c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x4c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65,
	0x2a, 0x2f, 0x0a, 0x09, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0d, 0x0a,
	0x09, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x44, 0x61, 0x74, 0x61, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x10,
	0x01, 0x2a, 0x5e, 0x0a, 0x06, 0x4f, 0x70, 0x54, 0x79, 0x70, 0x65, 0x12, 0x07, 0x0a, 0x03, 0x50,
	0x75, 0x74, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x10, 0x01,
	0x12, 0x07, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x10, 0x03, 0x12, 0x11, 0x0a, 0x0d, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x10, 0x04, 0x12, 0x11,
	0x0a, 0x0d, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x10,
	0x05, 0x32, 0xd2, 0x01, 0x0a, 0x0b, 0x52, 0x61, 0x66, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x3c, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65,
	0x12, 0x14, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f,
	0x74, 0x65, 0x41, 0x72, 0x67, 0x73, 0x1a, 0x15, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12,
	0x42, 0x0a, 0x0d, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x12, 0x16, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74,
	0x72, 0x69, 0x65, 0x73, 0x41, 0x72, 0x67, 0x73, 0x1a, 0x17, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x41,
	0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x08, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x12,
	0x18, 0x2e, 0x70, 0x62, 0x73, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x53, 0x6e, 0x61,
	0x70, 0x73, 0x68, 0x6f, 0x74, 0x41, 0x72, 0x67, 0x73, 0x1a, 0x19, 0x2e, 0x70, 0x62, 0x73, 0x2e,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2e, 0x2f, 0x74, 0x69, 0x6e,
	0x6e, 0x72, 0x61, 0x66, 0x74, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tinnraft_proto_rawDescOnce sync.Once
	file_tinnraft_proto_rawDescData = file_tinnraft_proto_rawDesc
)

func file_tinnraft_proto_rawDescGZIP() []byte {
	file_tinnraft_proto_rawDescOnce.Do(func() {
		file_tinnraft_proto_rawDescData = protoimpl.X.CompressGZIP(file_tinnraft_proto_rawDescData)
	})
	return file_tinnraft_proto_rawDescData
}

var file_tinnraft_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_tinnraft_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_tinnraft_proto_goTypes = []interface{}{
	(EntryType)(0),               // 0: pbs.EntryType
	(OpType)(0),                  // 1: pbs.OpType
	(*RequestVoteArgs)(nil),      // 2: pbs.RequestVoteArgs
	(*RequestVoteReply)(nil),     // 3: pbs.RequestVoteReply
	(*Entry)(nil),                // 4: pbs.Entry
	(*AppendEntriesArgs)(nil),    // 5: pbs.AppendEntriesArgs
	(*AppendEntriesReply)(nil),   // 6: pbs.AppendEntriesReply
	(*ApplyMsg)(nil),             // 7: pbs.ApplyMsg
	(*InstallSnapshotArgs)(nil),  // 8: pbs.InstallSnapshotArgs
	(*InstallSnapshotReply)(nil), // 9: pbs.InstallSnapshotReply
	(*CommandArgs)(nil),          // 10: pbs.CommandArgs
	(*CommandReply)(nil),         // 11: pbs.CommandReply
}
var file_tinnraft_proto_depIdxs = []int32{
	0, // 0: pbs.Entry.Type:type_name -> pbs.EntryType
	4, // 1: pbs.AppendEntriesArgs.Entries:type_name -> pbs.Entry
	1, // 2: pbs.CommandArgs.Op_type:type_name -> pbs.OpType
	2, // 3: pbs.RaftService.RequestVote:input_type -> pbs.RequestVoteArgs
	5, // 4: pbs.RaftService.AppendEntries:input_type -> pbs.AppendEntriesArgs
	8, // 5: pbs.RaftService.Snapshot:input_type -> pbs.InstallSnapshotArgs
	3, // 6: pbs.RaftService.RequestVote:output_type -> pbs.RequestVoteReply
	6, // 7: pbs.RaftService.AppendEntries:output_type -> pbs.AppendEntriesReply
	9, // 8: pbs.RaftService.Snapshot:output_type -> pbs.InstallSnapshotReply
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_tinnraft_proto_init() }
func file_tinnraft_proto_init() {
	if File_tinnraft_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tinnraft_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestVoteArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tinnraft_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestVoteReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tinnraft_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tinnraft_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntriesArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tinnraft_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppendEntriesReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tinnraft_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplyMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tinnraft_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstallSnapshotArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tinnraft_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstallSnapshotReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tinnraft_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommandArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tinnraft_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommandReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tinnraft_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tinnraft_proto_goTypes,
		DependencyIndexes: file_tinnraft_proto_depIdxs,
		EnumInfos:         file_tinnraft_proto_enumTypes,
		MessageInfos:      file_tinnraft_proto_msgTypes,
	}.Build()
	File_tinnraft_proto = out.File
	file_tinnraft_proto_rawDesc = nil
	file_tinnraft_proto_goTypes = nil
	file_tinnraft_proto_depIdxs = nil
}
