# Usage: generate-react.cmd  <servicename> 
SVC_PATH=../go/$1
OUT=../react/opentp-client/src/serverapi
protoc $SVC_PATH/*.proto --js_out=import_style=commonjs,binary:$OUT --grpc-web_out=import_style=typescript,mode=grpcwebtext:$OUT --proto_path=$SVC_PATH:.
protoc ./*.proto --js_out=import_style=commonjs,binary:$OUT --grpc-web_out=import_style=typescript,mode=grpcwebtext:$OUT --proto_path=$SVC_PATH:.

# A workaround in the typescript plugin for grpcweb 
for f in $OUT/*.js 
do
 if !(grep -q "eslint-disable" $f) 
 then
    sed -i '1s/^/\/* eslint-disable *\/\n/' $f
 fi
done



