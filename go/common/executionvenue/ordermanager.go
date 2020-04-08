package executionvenue

import (
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/model"
)

type OrderManager interface {
	CancelOrder(id *api.OrderId) error
	CreateAndRouteOrder(params *api.CreateAndRouteOrderParams) (*api.OrderId, error)
	SetOrderStatus(orderId string, status model.OrderStatus) error
	UpdateTradedQuantity(orderId string, lastPrice model.Decimal64, lastQty model.Decimal64) error
	Close()
}