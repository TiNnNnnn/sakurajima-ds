rm -rf ../tinnraftpb/*
protoc -I ../pbs ../pbs/test.proto --go_out=../pbs/ --go-grpc_out=../pbs/