package orderstore

import (
	"github.com/ettec/open-trading-platform/execution-venue/pb"
)


type OrderStore interface {
	Write(order *pb.Order) error
	Close()
}
