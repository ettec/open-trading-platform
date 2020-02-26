module github.com/ettec/open-trading-platform/go/execution-venue

go 1.13

require (
	github.com/ettec/open-trading-platform/go/model v0.0.0
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/quickfixgo/quickfix v0.6.0
	github.com/segmentio/kafka-go v0.3.4
	github.com/shopspring/decimal v0.0.0-20191009025716-f1972eb1d1f5
	google.golang.org/grpc v1.25.1
)

replace github.com/ettec/open-trading-platform/go/model v0.0.0 => ../model
