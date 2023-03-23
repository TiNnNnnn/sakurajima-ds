package tinnraft

const (
	VOTE_FOR_NO_ONE = -1
)

var RAFTLOG_PREFIX = []byte{0x11, 0x11, 0x19, 0x96}

var FIRST_IDX_KEY = []byte{0x88, 0x88}

var LAST_IDX_KEY = []byte{0x99, 0x99}

var RAFT_STATE_KEY = []byte{0x19, 0x49}

const INIT_LOG_INDEX = 0

var SNAPSHOT_STATE_KEY = []byte{0x19, 0x97}
