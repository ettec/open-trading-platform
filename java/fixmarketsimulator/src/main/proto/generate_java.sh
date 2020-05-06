protoc *.proto  --java_out=../java
protoc --plugin=protoc-gen-grpc-java=../../../protoc-gen-grpc-java.exe --grpc-java_out=../java --proto_path=. marketdataserver.proto
protoc --plugin=protoc-gen-grpc-java=../../../protoc-gen-grpc-java.exe --grpc-java_out=../java --proto_path=. orderentryapi.proto
