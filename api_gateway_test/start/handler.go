package start

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var logChan = make(chan string)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method

	if m == http.MethodPut {
		configOp := strings.Split(r.URL.EscapedPath(), "/")[2]
		log.Printf("configOp: %v\n", configOp)
		if configOp == "kvserver" {
			go StartKvServer(w, r)
			go readLog()
			return
		}
	}
	if m == http.MethodGet {
		//get(w, r)

	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func readLog() {
	for {
		c := <-logChan
		if len(c) > 0 {
			fmt.Print("***** " + c)
			
		}
	}
}

func StartKvServer(w http.ResponseWriter, r *http.Request) {

	sid := GetServerIdFromHeader(r.Header)
	if sid == "" {
		return
	}
	cmd := exec.Command("./../output/kvserver", sid)

	//stdin,stdout,stderr重定向
	cmd.Stdin = os.Stdin
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	go func() {
		for {
			l, err := out.ReadString('\n')
			if err != nil && err.Error() != "EOF" {
				log.Print(err)
				continue
			}
			//fmt.Print(l)
			logChan <- l

		}
	}()

	//go readLog()

	//创建子进程，并进行程序替换(父进程之后的代码不会再运行)
	cmd.Run()

}

func GetServerIdFromHeader(h http.Header) string {
	kvs_id := h.Get("kvs_id")
	if id, _ := strconv.Atoi(kvs_id); id < 0 {
		log.Println("a illegal serverId! it should be more than 0")
		return ""
	}
	return kvs_id
}
