package ordermanager

import (
	"fmt"
	"github.com/coronationstreet/open-trading-platform/execution-venue/internal/ordercache"
	"github.com/coronationstreet/open-trading-platform/execution-venue/internal/ordergateway"
	"github.com/coronationstreet/open-trading-platform/execution-venue/pb"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"log"
)

var zero decimal.Decimal

func init() {
	zero = decimal.New(0,0)
}



type OrderManager interface {
	CancelOrder(id *pb.OrderId) error
	CreateAndRouteOrder(params *pb.CreateAndRouteOrderParams) (*pb.OrderId, error)
	SetOrderStatus(orderId string, status pb.OrderStatus) error
	UpdateTradedQuantity(orderId string, lastPrice pb.Decimal64, lastQty pb.Decimal64) error
	Close()
}

type orderManagerImpl struct {
	//TODO Prometheus stat these queues
	createOrderChan     chan createAndRouteOrderCmd
	cancelOrderChan     chan cancelOrderCmd
	setOrderStatusChan  chan setOrderStatusCmd
	updateTradedQntChan chan updateTradedQntCmd

	closeChan chan struct{}

	orderStore *ordercache.OrderCache
	gateway    ordergateway.OrderGateway
}

func NewOrderManager(cache *ordercache.OrderCache, gateway ordergateway.OrderGateway) OrderManager {

	om := orderManagerImpl{}

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


func (om *orderManagerImpl) SetOrderStatus(orderId string, status pb.OrderStatus) error {
	log.Printf("updating order %v status to %v", orderId, status)

	resultChan := make(chan errorCmdResult)

	om.setOrderStatusChan <- setOrderStatusCmd{
		orderId:    orderId,
		status:     status,
		ResultChan: resultChan,
	}

	result := <-resultChan

	log.Printf("update order %v status result:%v", orderId, result)

	return result.Error
}



func (om *orderManagerImpl) UpdateTradedQuantity(orderId string, lastPrice pb.Decimal64, lastQty pb.Decimal64) error {
	log.Printf( orderId +":adding execution for price %v and quantity %v", lastPrice, lastQty)

	resultChan := make(chan errorCmdResult)

	om.updateTradedQntChan <- updateTradedQntCmd{
		orderId:    orderId,
		lastPrice:  lastPrice,
		lastQty:    lastQty,
		ResultChan: resultChan,
	}

	result := <-resultChan

	log.Printf(orderId + ":update traded quantity result:%v", result)

	return result.Error
}

func (om *orderManagerImpl) CreateAndRouteOrder(params *pb.CreateAndRouteOrderParams) (*pb.OrderId, error) {

	log.Printf("creating order with params:%v", params)

	resultChan := make(chan createAndRouteOrderCmdResult)

	om.createOrderChan <- createAndRouteOrderCmd{
		Params:     params,
		ResultChan: resultChan,
	}

	result := <-resultChan

	log.Printf("create and route order result:%v", result)
	return result.OrderId, result.Error
}

func (om *orderManagerImpl) CancelOrder(id *pb.OrderId) error {

	log.Print(id.OrderId + ":cancelling order")

	resultChan := make(chan errorCmdResult)

	om.cancelOrderChan <- cancelOrderCmd{
		Params:     id,
		ResultChan: resultChan,
	}

	result := <-resultChan

	log.Printf(id.OrderId + ":cancel order result:%v", result)

	return result.Error
}



func (om *orderManagerImpl) executeUpdateTradedQntCmd(id string, lastPrice pb.Decimal64, lastQty pb.Decimal64, resultChan chan errorCmdResult) {

	order, exists := om.orderStore.GetOrder(id)
	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("update traded quantity failed, no order found for id %v", id)}
	}

	if order.TargetStatus == pb.OrderStatus_LIVE {
		order.SetStatus(pb.OrderStatus_LIVE)
	}

	order.AvgTradePrice = calculateAveragePrice(order.AvgTradePrice, order.TradedQuantity, &lastPrice, &lastQty)

	order.RemainingQuantity = pb.ToDecimal64(order.RemainingQuantity.AsDecimal().Sub(lastQty.AsDecimal()))
	order.TradedQuantity = pb.ToDecimal64(order.TradedQuantity.AsDecimal().Add(lastQty.AsDecimal()))

	if order.RemainingQuantity.AsDecimal().LessThanOrEqual(zero) {
		err := order.SetStatus(pb.OrderStatus_FILLED)
		if err != nil {
			resultChan <- errorCmdResult{Error: err}
			return
		}
	}

	err := om.orderStore.Store(order)
	resultChan <- errorCmdResult{Error:err}

}

func calculateAveragePrice( avgPrice *pb.Decimal64, tradeQnt *pb.Decimal64, lastPx *pb.Decimal64, lastQty *pb.Decimal64) *pb.Decimal64 {
	totalTradeValue := avgPrice.AsDecimal().Mul(tradeQnt.AsDecimal()).Add(lastPx.AsDecimal().Mul(lastQty.AsDecimal()))
	totalTradedQnt := tradeQnt.AsDecimal().Add(lastQty.AsDecimal())
	newAvgPrice := totalTradeValue.Div(totalTradedQnt)

	return pb.ToDecimal64(newAvgPrice)
}


func (om *orderManagerImpl) executeSetOrderStatusCmd(id string, status pb.OrderStatus, resultChan chan errorCmdResult) {

	order, exists := om.orderStore.GetOrder(id)
	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("set order status failed, no order found for id %v", id)}
	}

	err := order.SetStatus(status)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}


	err = om.orderStore.Store(order)
	resultChan <- errorCmdResult{Error:err}

}

func (om *orderManagerImpl) executeCancelOrderCmd(id *pb.OrderId, resultChan chan errorCmdResult) {

	order, exists := om.orderStore.GetOrder(id.OrderId)
	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("cancel order failed, no order found for id %v", id.OrderId)}
	}

	err := order.SetTargetStatus(pb.OrderStatus_CANCELLED)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	err = om.orderStore.Store(order)
	resultChan <- errorCmdResult{Error:err}
}

func (om *orderManagerImpl) executeCreateAndRouteOrderCmd(params *pb.CreateAndRouteOrderParams,
	resultChan chan createAndRouteOrderCmdResult) {

	uniqueId, err := uuid.NewUUID()
	if err != nil {
		resultChan <- createAndRouteOrderCmdResult{
			OrderId: nil,
			Error:   fmt.Errorf("failed to create new order id: %w", err),
		}
	}

	id := uniqueId.String()

	order := &pb.Order{
		Id:                id,
		Side:              params.Side,
		Quantity:          params.Quantity,
		Price:             params.Price,
		ListingId:         params.ListingId,
		RemainingQuantity: params.Quantity,
		Status:            pb.OrderStatus_NONE,
		TargetStatus:      pb.OrderStatus_LIVE,
	}

	err = om.orderStore.Store(order)
	if err != nil {
		resultChan <- createAndRouteOrderCmdResult{
			OrderId: nil,
			Error:   err,
		}
	}

	err = om.gateway.Send(order)

	resultChan <- createAndRouteOrderCmdResult{
		OrderId: &pb.OrderId{
			OrderId: id,
		},
		Error: err,
	}

}

type updateTradedQntCmd struct {
	orderId    string
	lastPrice  pb.Decimal64
	lastQty    pb.Decimal64
	ResultChan chan errorCmdResult
}

type setOrderStatusCmd struct {
	orderId    string
	status     pb.OrderStatus
	ResultChan chan errorCmdResult
}

type createAndRouteOrderCmd struct {
	Params     *pb.CreateAndRouteOrderParams
	ResultChan chan createAndRouteOrderCmdResult
}

type createAndRouteOrderCmdResult struct {
	OrderId *pb.OrderId
	Error   error
}

type cancelOrderCmd struct {
	Params     *pb.OrderId
	ResultChan chan errorCmdResult
}

type errorCmdResult struct {
	Error error
}
