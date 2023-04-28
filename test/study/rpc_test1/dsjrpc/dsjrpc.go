package dsjrpc

import (
	"encoding/json"
	"log"
	"reflect"
)

/*
	Emulation Rpc Library For DsJourney
*/

var Debug = true

type MsgArgs struct {
	clientname interface{} // clientend name
	methodName string      //rpc method name
	argsType   reflect.Type
	msg        []byte
	replyCh    chan MsgReply
}

type MsgReply struct {
	success bool
	msg     []byte
}

type ClientEnd struct {
	clientname interface{}
	ch         chan MsgArgs
	done       chan struct{}
}

/*
远程调用
*/
func (ce *ClientEnd) Invoke(MethodName string, args interface{}, reply interface{}) bool {

	argsBytes, err := json.Marshal(args)
	if err != nil {
		DsjLog("encoding args failed,err: ", err.Error())
		return false
	}

	req := MsgArgs{
		clientname: ce.clientname,
		methodName: MethodName,
		argsType:   reflect.TypeOf(args),
		replyCh:    make(chan MsgReply),
		msg:        argsBytes,
	}

	select {
	case ce.ch <- req: //发送args成功

	case <-ce.done: //发送args失败,destoryed
		return false
	}

	res := <-req.replyCh
	if res.success {
		err = json.Unmarshal(res.msg, reply)
		if err != nil {
			DsjLog("decoding args failed,err: ", err.Error())
		}
		return true
	} else {
		return false
	}

}

func MakeClientEnd(endaname interface{}) {
	
}

/*
Format Log
*/
func DsjLog(format string, a ...interface{}) (n int, err error) {
	if Debug {
		if format[len(format)-1] != '\n' {
			format += "\n"
		}
		log.Printf(format, a...)
	}
	return
}
