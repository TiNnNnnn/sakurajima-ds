default: kvserver kvclient configserver sharedserver apigateway
   
kvserver:
	go build -o output/kvserver test/kvs/kvserver.go
kvclient:
	go build -o output/kvclient test/kvc/kvclient.go 
configserver:
	go build -o output/cfgserver test/configserver/configServer.go 
sharedserver:
	go build -o output/sharedserver test/sharedserver/shareserver.go 
apigateway:
	go build -o output/apigateway test/apiServer/api_server.go
clean:
	rm -rf output/*
	