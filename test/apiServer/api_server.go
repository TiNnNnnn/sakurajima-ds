package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	api_gateway "sakurajima-ds/api_gateway_2"
	config "sakurajima-ds/api_gateway_2/config"
	objects "sakurajima-ds/api_gateway_2/objects"
	"sakurajima-ds/tinnraft"
	"sakurajima-ds/tinnraftpb"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

var logrpcAddr = "127.0.0.1:10030"
var apiSvr = api_gateway.MakeApiGatwayServer(logrpcAddr)
var addr = flag.String("addr", "0.0.0.0:10055", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func startkv(w http.ResponseWriter, r *http.Request) {
	go apiSvr.StartKvServer(w, r)
}

func startconfig(w http.ResponseWriter, r *http.Request) {
	go apiSvr.StartConfigServer(w, r)
}

func startshared(w http.ResponseWriter, r *http.Request) {
	go apiSvr.StartSharedServer(w, r)
}

func handleStop(w http.ResponseWriter, r *http.Request) {
	apiSvr.StopServer(w, r)
}

func handleStart(w http.ResponseWriter, r *http.Request) {

	m := r.Method
	if m == http.MethodPut {
		serverType := api_gateway.GetstypeFromHeader(r.Header)
		log.Printf("serverType: %v\n", serverType)

		if serverType == "kvserver" {
			go startkv(w, r)
		} else if serverType == "configserver" {
			go startconfig(w, r)
		} else if serverType == "sharedserver" {
			go startshared(w, r)
		}
	}
}

// 建立长连接，接受并转发日志到client
func handleLog(w http.ResponseWriter, r *http.Request) {

	//升级http为websocket 长连接
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade err:", err)
		return
	}
	defer c.Close()

	fmt.Println("create a connect from api_server to client success")

	go apiSvr.SendMutiLog(c)

	//go apiSvr.SendCommonLog(c)

	for {
		if <-apiSvr.StopChan {
			fmt.Println("close the connect from api_server to client")
			break
		}
	}

}

func handleClose(w http.ResponseWriter, r *http.Request) {
	apiSvr.StopChan <- true
}

// 启动日志接受服务
func LogRpcServer() {
	lis, err := net.Listen("tcp", logrpcAddr)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	tinnraft.DLog("server mutiLog listen on:  " + logrpcAddr)
	s := grpc.NewServer()

	tinnraftpb.RegisterRaftServiceServer(s, apiSvr)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func ObjectHandler(w http.ResponseWriter, r *http.Request) {
	objects.Handler(w, r, apiSvr)
}

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	config.Handler(w, r, apiSvr)
}

func main() {

	go LogRpcServer()

	log.Println("api gateway begining working...")
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/stop", handleStop)

	http.HandleFunc("/close", handleClose)
	http.HandleFunc("/log", handleLog)

	http.HandleFunc("/apis/", ObjectHandler)
	http.HandleFunc("/config/", ConfigHandler)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
