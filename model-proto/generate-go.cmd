# Usage: generate-go.cmd  <servicename> 
SVC_NAME=$1
SVC_PATH=../go/$SVC_NAME

mkdir -p $SVC_PATH/internal/model
protoc ./*.proto --go_out=plugins=grpc:$SVC_PATH/internal/model/ --proto_path=$SVC_PATH:.
protoc $SVC_PATH/$SVC_NAME.proto --go_out=plugins=grpc:$SVC_PATH/internal/model/ --proto_path=$SVC_PATH:.
