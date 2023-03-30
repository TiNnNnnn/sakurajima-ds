kvserver:
	go build -o output/kvserver test/kvserver.go
clean:
	rm -rf output/*