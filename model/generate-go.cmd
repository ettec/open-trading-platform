mkdir -p $1/model
mkdir -p $1/api
protoc $1/*.proto --go_out=plugins=grpc:$1/api/ --proto_path=$1:.
protoc ./*.proto --go_out=plugins=grpc:$1/model/ --proto_path=$1:.
