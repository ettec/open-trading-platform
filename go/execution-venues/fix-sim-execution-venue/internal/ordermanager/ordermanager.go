package ordermanager

import (
	"fmt"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/common/executionvenue"
	"github.com/ettec/open-trading-platform/go/execution-venues/fix-sim-execution-venue/internal/ordercache"
	"github.com/ettec/open-trading-platform/go/execution-venues/fix-sim-execution-venue/internal/ordergateway"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"log"
	"os"
	"time"
)

var zero decimal.Decimal

func init() {
	zero = decimal.New(0, 0)
}

type orderManagerImpl struct {
	//TODO Prometheus stat these queues
	createOrderChan     chan createAndRouteOrderCmd
	cancelOrderChan     chan cancelOrderCmd
	setOrderStatusChan  chan setOrderStatusCmd
	updateTradedQntChan chan updateTradedQntCmd

	closeChan chan struct{}

	execVenueId string
	orderStore  *ordercache.OrderCache
	gateway     ordergateway.OrderGateway
	log         *log.Logger
	errLog      *log.Logger
}

func NewOrderManager(cache *ordercache.OrderCache, gateway ordergateway.OrderGateway, execVenueId string) executionvenue.OrderManager {

	om := orderManagerImpl{
		log:         log.New(os.Stdout, "", log.Lshortfile|log.Ltime),
		errLog:      log.New(os.Stderr, "", log.Lshortfile|log.Ltime),
		execVenueId: execVenueId,
	}

	om.createOrderChan = make(chan createAndRouteOrderCmd, 100)
	om.cancelOrderChan = make(chan cancelOrderCmd, 100)
	om.setOrderStatusChan = make(chan setOrderStatusCmd, 100)
	om.updateTradedQntChan = make(chan updateTradedQntCmd, 100)

	om.closeChan = make(chan struct{}, 1)

	om.orderStore = cache
	om.gateway = gateway

	go om.executeOrderCommands()

	return &om
}

func (om *orderManagerImpl) executeOrderCommands() {

	for {
		select {
		case <-om.closeChan:
			return
		// Cancel Requests take priority over all other message types
		case oc := <-om.cancelOrderChan:
			om.executeCancelOrderCmd(oc.Params, oc.ResultChan)
		default:
			select {
			case oc := <-om.cancelOrderChan:
				om.executeCancelOrderCmd(oc.Params, oc.ResultChan)
			case cro := <-om.createOrderChan:
				om.executeCreateAndRouteOrderCmd(cro.Params, cro.ResultChan)
			case su := <-om.setOrderStatusChan:
				om.executeSetOrderStatusCmd(su.orderId, su.status, su.ResultChan)
			case tu := <-om.updateTradedQntChan:
				om.executeUpdateTradedQntCmd(tu.orderId, tu.lastPrice, tu.lastQty, tu.ResultChan)
			}
		}
	}

}

func (om *orderManagerImpl) Close() {
	om.closeChan <- struct{}{}
}

func (om *orderManagerImpl) SetOrderStatus(orderId string, status model.OrderStatus) error {
	om.log.Printf("updating order %v status to %v", orderId, status)

	resultChan := make(chan errorCmdResult)

	om.setOrderStatusChan <- setOrderStatusCmd{
		orderId:    orderId,
		status:     status,
		ResultChan: resultChan,
	}

	result := <-resultChan

	return result.Error
}

func (om *orderManagerImpl) UpdateTradedQuantity(orderId string, lastPrice model.Decimal64, lastQty model.Decimal64) error {
	om.log.Printf(orderId+":adding execution for price %v and quantity %v", lastPrice, lastQty)

	resultChan := make(chan errorCmdResult)

	om.updateTradedQntChan <- updateTradedQntCmd{
		orderId:    orderId,
		lastPrice:  lastPrice,
		lastQty:    lastQty,
		ResultChan: resultChan,
	}

	result := <-resultChan

	om.log.Printf(orderId+":update traded quantity result:%v", result)

	return result.Error
}

func (om *orderManagerImpl) CreateAndRouteOrder(params *api.CreateAndRouteOrderParams) (*api.OrderId, error) {

	resultChan := make(chan createAndRouteOrderCmdResult)

	om.createOrderChan <- createAndRouteOrderCmd{
		Params:     params,
		ResultChan: resultChan,
	}

	result := <-resultChan

	if result.Error != nil {
		return nil, result.Error
	}

	return result.OrderId, nil
}

func (om *orderManagerImpl) CancelOrder(id *api.OrderId) error {

	om.log.Print(id.OrderId + ":cancelling order")

	resultChan := make(chan errorCmdResult)

	om.cancelOrderChan <- cancelOrderCmd{
		Params:     id,
		ResultChan: resultChan,
	}

	result := <-resultChan

	om.log.Printf(id.OrderId+":cancel order result: %v", result)

	return result.Error
}

func (om *orderManagerImpl) executeUpdateTradedQntCmd(id string, lastPrice model.Decimal64, lastQty model.Decimal64, resultChan chan errorCmdResult) {

	order, exists := om.orderStore.GetOrder(id)
	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("update traded quantity failed, no order found for id %v", id)}
		return
	}

	if order.TargetStatus == model.OrderStatus_LIVE {
		err := order.SetStatus(model.OrderStatus_LIVE)
		if err != nil {
			resultChan <- errorCmdResult{Error: err}
			return
		}
	}

	order.AvgTradePrice = calculateAveragePrice(order.AvgTradePrice, order.TradedQuantity, &lastPrice, &lastQty)

	order.RemainingQuantity = model.ToDecimal64(order.RemainingQuantity.AsDecimal().Sub(lastQty.AsDecimal()))
	order.TradedQuantity = model.ToDecimal64(order.TradedQuantity.AsDecimal().Add(lastQty.AsDecimal()))

	if order.RemainingQuantity.AsDecimal().LessThanOrEqual(zero) {
		err := order.SetStatus(model.OrderStatus_FILLED)
		if err != nil {
			resultChan <- errorCmdResult{Error: err}
			return
		}
	}

	err := om.orderStore.Store(order)
	resultChan <- errorCmdResult{Error: err}

}

func calculateAveragePrice(avgPrice *model.Decimal64, tradeQnt *model.Decimal64, lastPx *model.Decimal64, lastQty *model.Decimal64) *model.Decimal64 {
	totalTradeValue := avgPrice.AsDecimal().Mul(tradeQnt.AsDecimal()).Add(lastPx.AsDecimal().Mul(lastQty.AsDecimal()))
	totalTradedQnt := tradeQnt.AsDecimal().Add(lastQty.AsDecimal())
	newAvgPrice := totalTradeValue.Div(totalTradedQnt)

	return model.ToDecimal64(newAvgPrice)
}

func (om *orderManagerImpl) executeSetOrderStatusCmd(id string, status model.OrderStatus, resultChan chan errorCmdResult) {

	order, exists := om.orderStore.GetOrder(id)
	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("set order status failed, no order found for id %v", id)}
		return
	}

	err := order.SetStatus(status)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	err = om.orderStore.Store(order)
	resultChan <- errorCmdResult{Error: err}

}

func (om *orderManagerImpl) executeCancelOrderCmd(id *api.OrderId, resultChan chan errorCmdResult) {

	order, exists := om.orderStore.GetOrder(id.OrderId)
	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("cancel order failed, no order found for id %v", id.OrderId)}
		return
	}

	err := order.SetTargetStatus(model.OrderStatus_CANCELLED)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	err = om.orderStore.Store(order)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	err = om.gateway.Cancel(order)

	resultChan <- errorCmdResult{Error: err}
}

func (om *orderManagerImpl) executeCreateAndRouteOrderCmd(params *api.CreateAndRouteOrderParams,
	resultChan chan createAndRouteOrderCmdResult) {

	uniqueId, err := uuid.NewUUID()
	if err != nil {
		resultChan <- createAndRouteOrderCmdResult{
			OrderId: nil,
			Error:   fmt.Errorf("failed to create new order id: %w", err),
		}
	}

	id := uniqueId.String()

	now := time.Now()

	order := &model.Order{
		Id:                id,
		Side:              params.OrderSide,
		Quantity:          params.Quantity,
		Price:             params.Price,
		ListingId:         params.Listing.GetId(),
		RemainingQuantity: params.Quantity,
		Status:            model.OrderStatus_NONE,
		TargetStatus:      model.OrderStatus_LIVE,
		Created: &model.Timestamp{
			Seconds:     now.Unix(),
			Nanoseconds: int32(now.Nanosecond()),
		},
		PlacedWithExecVenueId: om.execVenueId,
	}

	err = om.orderStore.Store(order)
	if err != nil {
		resultChan <- createAndRouteOrderCmdResult{
			OrderId: nil,
			Error:   err,
		}

		return
	}

	err = om.gateway.Send(order, params.Listing)

	resultChan <- createAndRouteOrderCmdResult{
		OrderId: &api.OrderId{
			OrderId: id,
		},
		Error: err,
	}

}

type updateTradedQntCmd struct {
	orderId    string
	lastPrice  model.Decimal64
	lastQty    model.Decimal64
	ResultChan chan errorCmdResult
}

type setOrderStatusCmd struct {
	orderId    string
	status     model.OrderStatus
	ResultChan chan errorCmdResult
}

type createAndRouteOrderCmd struct {
	Params     *api.CreateAndRouteOrderParams
	ResultChan chan createAndRouteOrderCmdResult
}

type createAndRouteOrderCmdResult struct {
	OrderId *api.OrderId
	Error   error
}

type cancelOrderCmd struct {
	Params     *api.OrderId
	ResultChan chan errorCmdResult
}

type errorCmdResult struct {
	Error error
}
