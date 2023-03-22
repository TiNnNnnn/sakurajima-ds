rm -rf ../tinnraftpb/*
protoc -I ../pbs ../pbs/tinnraft.proto --go_out=../pbs/ --go-grpc_out=../pbs/