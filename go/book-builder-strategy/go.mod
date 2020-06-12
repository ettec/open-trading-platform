module github.com/ettec/open-trading-platform/go/book-builder-strategy

go 1.13

require (
	github.com/ettec/open-trading-platform/go/common v0.0.0
	github.com/ettec/open-trading-platform/go/model v0.0.0
	github.com/ettech/open-trading-platform/go/market-data/market-data-common v0.0.0
	github.com/golang/protobuf v1.4.0
	github.com/google/uuid v1.1.1
	google.golang.org/grpc v1.25.1
)

replace github.com/ettec/open-trading-platform/go/common v0.0.0 => ../common

replace github.com/ettec/open-trading-platform/go/model v0.0.0 => ../model

replace github.com/ettech/open-trading-platform/go/market-data/market-data-common v0.0.0 => ../market-data/market-data-common
