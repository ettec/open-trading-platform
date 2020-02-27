module github.com/ettec/open-trading-platform/go/view-service

go 1.13

require (
	github.com/ettec/open-trading-platform/go/model v0.0.0
	github.com/golang/protobuf v1.3.2
	github.com/segmentio/kafka-go v0.3.4
	google.golang.org/grpc v1.25.1
)

replace github.com/ettec/open-trading-platform/go/model v0.0.0 => ../model
