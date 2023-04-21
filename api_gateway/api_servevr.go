package main

import (
	"log"
	"net/http"
	"sakurajima-ds/api_gateway/config"
	"sakurajima-ds/api_gateway/objects"
	"sakurajima-ds/api_gateway/start"

	//_ "github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/apis/", objects.Handler)
	http.HandleFunc("/config/", config.Handler)
	http.HandleFunc("/start/", start.Handler)
	log.Fatal(http.ListenAndServe(":10055", nil))
}


