kv_server:
	go build -o output/kv_server cmd/kvraft/kvserver.go
clean:
	rm -rf output/*