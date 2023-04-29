package main

import (
	"log"
	"net/http"
	"sakurajima-ds/api_gateway_test/config"
	"sakurajima-ds/api_gateway_test/objects"
	"sakurajima-ds/api_gateway_test/start"

)

func main() {
	http.HandleFunc("/apis/", objects.Handler)
	http.HandleFunc("/config/", config.Handler)
	http.HandleFunc("/start/", start.Handler)
	log.Fatal(http.ListenAndServe(":10055", nil))
}


