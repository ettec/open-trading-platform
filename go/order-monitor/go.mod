module github.com/ettech/open-trading-platform/go/order-monitor

go 1.13

require (
	github.com/ettec/otp-common v0.0.0-20200627105436-e95192c9a900
	github.com/ettec/otp-evcommon v0.0.0-20200627105543-deacb041b39e
	github.com/ettec/otp-model v0.0.2-0.20200627105317-ed67f7c52141
	github.com/golang/protobuf v1.4.0
	github.com/google/uuid v1.1.1
	github.com/prometheus/client_golang v1.6.0
	github.com/segmentio/kafka-go v0.3.4
	google.golang.org/grpc v1.25.1
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4
)

replace github.com/ettec/otp-evcommon v0.0.0-20200627105543-deacb041b39e => ../execution-venues/common
