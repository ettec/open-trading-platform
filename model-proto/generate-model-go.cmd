# Usage: generate-go.cmd  <servicename> 
SVC_NAME=$1
SVC_PATH=../go/$SVC_NAME

protoc ./*.proto --go_out=plugins=grpc:$SVC_PATH/ --proto_path=$SVC_PATH:.
