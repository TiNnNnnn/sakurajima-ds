default: kvserver kvclient configserver

kvserver:
	go build -o output/kvserver test/kvserver.go
kvclient:
	go build -o output/kvclient test/kvc/kvclient.go 
configserver:
	go build -o output/cfgserver test/configServer/configServer.go 
clean:
	rm -rf output/*

	