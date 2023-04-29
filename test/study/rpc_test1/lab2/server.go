package main

import (
	"log"
	"runtime"
	rpc "sakurajima-ds/test/study/labrpc"
	"strconv"
	"sync"
)

type KvServer struct {
	stm map[string]string
	mu  sync.Mutex
}

type PutArgs struct {
	k string
	v string
}

type PutReply struct {
	errcode int
	errMsg  string
	Success bool
}

func (ks *KvServer) Handler(args string, reply *int) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	//ks.stm["hello"] = "world"
	*reply, _ = strconv.Atoi(args)
}

func test1() {
	runtime.GOMAXPROCS(4)

	rn := rpc.MakeNetwork()
	defer rn.Cleanup()

	//构建clientend
	e := rn.MakeEnd("end1-99")

	//注册方法
	kvs := &KvServer{}
	svc := rpc.MakeService(kvs)

	//构建server
	rs := rpc.MakeServer()
	rs.AddService(svc)
	rn.AddServer("server-99", rs)

	rn.Connect("end1-99", "server-99")
	rn.Enable("end1-99", true)

	{
		reply := 0
		//调用远程rpc
		e.Call("KvServer.Handler", "9099", &reply)
		if reply != 9099 {
			log.Println("wrony reply from Handler")
		} else {
			log.Println("right reply from Handler")
		}
	}

}

func main() {
	test1()
}
