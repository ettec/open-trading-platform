# Usage: generate-react.cmd  <servicename> 
SVC_PATH=../go/$1
protoc $SVC_PATH/*.proto --js_out=import_style=commonjs,binary:../react/opentp-client/src/serverapi --grpc-web_out=import_style=typescript,mode=grpcwebtext:../react/opentp-client/src/serverapi --proto_path=$SVC_PATH:.
protoc ./*.proto --js_out=import_style=commonjs,binary:../react/opentp-client/src/serverapi --grpc-web_out=import_style=typescript,mode=grpcwebtext:../react/opentp-client/src/serverapi --proto_path=$SVC_PATH:.
