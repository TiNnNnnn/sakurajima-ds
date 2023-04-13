package main

import (
	"log"
	"net/http"
	"os"
	"sakurajima-ds/api_gateway/objects"
)

func main() {
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))

}

/*
网关服务器
*/
