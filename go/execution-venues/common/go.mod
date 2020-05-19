module github.com/ettec/open-trading-platform/go/execution-venues/common

go 1.13

require (
	github.com/ettec/open-trading-platform/go/common v0.0.0
	github.com/ettec/open-trading-platform/go/model v0.0.0
	github.com/golang/protobuf v1.4.0
	github.com/google/uuid v1.1.1
	github.com/segmentio/kafka-go v0.3.4
	github.com/shopspring/decimal v0.0.0-20191009025716-f1972eb1d1f5

)

replace github.com/ettec/open-trading-platform/go/common v0.0.0 => ../../common

replace github.com/ettec/open-trading-platform/go/model v0.0.0 => ../../model
