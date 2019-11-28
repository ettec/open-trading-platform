mkdir -p ../client-ts
protoc $1 --js_out=import_style=commonjs,binary:../client-ts --grpc-web_out=import_style=typescript,mode=grpcwebtext:../client-ts
