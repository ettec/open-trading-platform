package ordergateway

import "github.com/ettec/open-trading-platform/execution-venue/pb"

type OrderGateway interface {
	Send(order *pb.Order) error
}
