# Usage: generate-go.cmd <path> 
SVC_PATH=$1

protoc ./*.proto --go_out=plugins=grpc:$SVC_PATH/ --proto_path=$SVC_PATH:.
