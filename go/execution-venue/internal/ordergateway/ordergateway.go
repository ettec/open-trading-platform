package ordergateway

import (
	"github.com/ettec/open-trading-platform/execution-venue/internal/model"
)

type OrderGateway interface {
	Send(order *model.Order) error
}
