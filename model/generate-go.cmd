mkdir -p $1/pb
protoc $1/*.proto --go_out=plugins=grpc:$1/pb/ --proto_path=$1:.
protoc ./*.proto --go_out=plugins=grpc:$1/pb/ --proto_path=$1:.
