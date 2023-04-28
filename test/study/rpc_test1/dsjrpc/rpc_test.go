package dsjrpc

import (
	"fmt"
	"testing"
)

type Raft struct {
	client *ClientEnd
}

func MakeRaft() *Raft {

	newclient := &ClientEnd{
		clientname: "test",
		ch:         make(chan MsgArgs),
		done:       make(chan struct{}),
	}

	newraft := &Raft{
		client: newclient,
	}
	return newraft
}

func TestBasic(t *testing.T) {

	r := MakeRaft()
	args := ""
	reply := ""

	ok := r.client.Invoke("Raftr.RequestVote", args, reply)
	fmt.Println(ok)
}
