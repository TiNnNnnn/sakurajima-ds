package main

import (
	"log"
	"net/http"
	"sakurajima-ds/api_gateway/config"
	"sakurajima-ds/api_gateway/objects"
)

func main() {
	http.HandleFunc("/apis/", objects.Handler)
	http.HandleFunc("/config/", config.Handler)
	log.Fatal(http.ListenAndServe(":10055", nil))

}

/*
网关服务器
*/
