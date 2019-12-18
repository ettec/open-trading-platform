package ordergateway

import (
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/model"
)

type OrderGateway interface {
	Send(order *model.Order, listing *model.Listing) error
	Cancel(order *model.Order) error
}
