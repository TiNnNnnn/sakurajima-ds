default: kvserver kvclient

kvserver:
	go build -o output/kvserver test/kvserver.go
kvclient:
	go build -o output/kvclient test/kvc/kvclient.go 
clean:
	rm -rf output/*