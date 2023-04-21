package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gorilla/websocket"
)

var logChan = make(chan string)

var addr = flag.String("addr", "0.0.0.0:10055", "http service address")

var upgrader = websocket.Upgrader{}

func StartKvServer(w http.ResponseWriter, r *http.Request) {

	sid := GetServerIdFromHeader(r.Header)
	if sid == "" {
		return
	}
	cmd := exec.Command("./../output/kvserver", sid)

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
			logChan <- l

		}
	}()

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

func start(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	go StartKvServer(w, r)

	for {
		line := <-logChan
		if len(line) > 0 {
			fmt.Print("***** " + line)
			c.WriteMessage(1, []byte(line))

		}
	}
}



func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/start", start)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
