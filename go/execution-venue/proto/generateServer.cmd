mkdir -p ../pb
protoc $1 --go_out=plugins=grpc:../pb/
