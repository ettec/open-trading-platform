package ordergateway

import "github.com/coronationstreet/open-trading-platform/execution-venue/pb"

type OrderGateway interface {
	Send(order *pb.Order) error
}
