default: kvserver kvclient configserver sharedserver

kvserver:
	go build -o output/kvserver test/kvs/kvserver.go
kvclient:
	go build -o output/kvclient test/kvc/kvclient.go 
configserver:
	go build -o output/cfgserver test/configserver/configServer.go 
sharedserver:
	go build -o output/sharedserver test/sharedserver/shareserver.go 
clean:
	rm -rf output/*

	