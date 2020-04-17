module github.com/ettec/open-trading-platform/go/order-router

go 1.13

require (
	github.com/ettec/open-trading-platform/go/common v0.0.0
	github.com/ettec/open-trading-platform/go/model v0.0.0
	github.com/golang/protobuf v1.3.2
	google.golang.org/grpc v1.25.1
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
)

replace github.com/ettec/open-trading-platform/go/model v0.0.0 => ../../model

replace github.com/ettec/open-trading-platform/go/common v0.0.0 => ../../common
