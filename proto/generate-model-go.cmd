# Usage: generate-go-model.cmd   
SVC_PATH=../go/model

protoc ./*.proto --go_out=plugins=grpc:$SVC_PATH/ --proto_path=$SVC_PATH:.
