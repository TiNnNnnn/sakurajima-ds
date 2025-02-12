package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	api_gateway "sakurajima-ds/api_gateway"
	config "sakurajima-ds/api_gateway/config"
	objects "sakurajima-ds/api_gateway/objects"
	"sakurajima-ds/common"
	"sakurajima-ds/tinnraftpb"
	"strings"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

var logrpcAddr = "127.0.0.1:10030"
var apiSvr = api_gateway.MakeApiGatwayServer(logrpcAddr)

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
	setupCORS(&w, r)
	apiSvr.StopServer(w, r)
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
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

	apiSvr.IsConnect = true

	//apiSvr.MutiChan = make(chan *tinnraftpb.LogArgs, 1024)

	fmt.Println("create a connect from api_server to client success")

	go apiSvr.SendMutiLog(c)

	//go apiSvr.SendCommonLog(c)

	for {
		if <-apiSvr.StopChan {
			fmt.Println("close the connect from api_server to client")
			apiSvr.IsConnect = false
			break
		}
	}

}

func handleClose(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	apiSvr.StopChan <- true
}

// 启动日志接受服务
func LogRpcServer() {
	lis, err := net.Listen("tcp", logrpcAddr)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	log.Printf("server mutiLog listen on:  %v", logrpcAddr)
	s := grpc.NewServer()

	tinnraftpb.RegisterRaftServiceServer(s, apiSvr)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func ObjectHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	objects.Handler(w, r, apiSvr)
}

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	config.Handler(w, r, apiSvr)
}

func handleClear(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w, r)
	apiSvr.ClearDates(w, r)
}

func setupCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusNoContent)
		return
	}
}

func main() {
	common.Getstartlogo()

	fmt.Println("========================================================================================")
	fmt.Println("usage: apigateway [api_gatway addr] [log_server addr] [cfgserver_path,sharedserver_path]")
	fmt.Println("----default args---------------")
	fmt.Println("api_gatway addr: 0.0.0.0:10055")
	fmt.Println("log_server:    : 127.0.0.1:10030")
	fmt.Println("-------------------------------")
	fmt.Println("========================================================================================")

	gateAddr := "0.0.0.0:10055"

	if len(os.Args) == 2 {
		gateAddr = os.Args[1]
	}

	if len(os.Args) == 3 {
		gateAddr = os.Args[1]
		logrpcAddr = os.Args[2]
	}

	if len(os.Args) == 4 {
		gateAddr = os.Args[1]
		logrpcAddr = os.Args[2]
		pathList := strings.Split(os.Args[3], ",")
		if len(pathList) != 2 {
			fmt.Printf("wrong path with cfgserver and sharedserver!")
			return
		}
		apiSvr.ServerPath = append(apiSvr.ServerPath, pathList...)
	}

	var addr = flag.String("addr", gateAddr, "http service address")

	go LogRpcServer()

	log.Println("api gateway begining working...")
	log.Printf("api_gatway listen on: %v", gateAddr)
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/stop", handleStop)

	http.HandleFunc("/close", handleClose)
	http.HandleFunc("/log", handleLog)
	http.HandleFunc("/clear", handleClear)

	http.HandleFunc("/apis/", ObjectHandler)
	http.HandleFunc("/config/", ConfigHandler)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
