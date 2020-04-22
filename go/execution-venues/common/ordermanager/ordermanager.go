package ordermanager

import (
	"fmt"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/google/uuid"

	"github.com/ettec/open-trading-platform/go/execution-venues/common/ordercache"
	"github.com/ettec/open-trading-platform/go/execution-venues/common/ordergateway"

	"github.com/ettec/open-trading-platform/go/model"
	"log"
	"os"
)

type OrderManager interface {
	CancelOrder(id *api.CancelOrderParams) error
	CreateAndRouteOrder(params *api.CreateAndRouteOrderParams) (*api.OrderId, error)
	SetOrderStatus(orderId string, status model.OrderStatus) error
	AddExecution(orderId string, lastPrice model.Decimal64, lastQty model.Decimal64, execId string) error
	Close()
}

type orderManagerImpl struct {
	//TODO Prometheus stat these queues
	createOrderChan    chan createAndRouteOrderCmd
	cancelOrderChan    chan cancelOrderCmd
	setOrderStatusChan chan setOrderStatusCmd
	addExecChan        chan addExecutionCmd

	closeChan chan struct{}

	execVenueId string
	orderStore  *ordercache.OrderCache
	gateway     ordergateway.OrderGateway
	log         *log.Logger
	errLog      *log.Logger
}

func NewOrderManager(cache *ordercache.OrderCache, gateway ordergateway.OrderGateway, execVenueId string) OrderManager {

	om := orderManagerImpl{
		log:         log.New(os.Stdout, "", log.Lshortfile|log.Ltime),
		errLog:      log.New(os.Stderr, "", log.Lshortfile|log.Ltime),
		execVenueId: execVenueId,
	}

	om.createOrderChan = make(chan createAndRouteOrderCmd, 100)
	om.cancelOrderChan = make(chan cancelOrderCmd, 100)
	om.setOrderStatusChan = make(chan setOrderStatusCmd, 100)
	om.addExecChan = make(chan addExecutionCmd, 100)

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
			case tu := <-om.addExecChan:
				om.executeUpdateTradedQntCmd(tu.orderId, tu.lastPrice, tu.lastQty, tu.execId, tu.ResultChan)
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

func (om *orderManagerImpl) AddExecution(orderId string, lastPrice model.Decimal64, lastQty model.Decimal64,
	execId string) error {
	om.log.Printf(orderId+":adding execution for price %v and quantity %v", lastPrice, lastQty)

	resultChan := make(chan errorCmdResult)

	om.addExecChan <- addExecutionCmd{
		orderId:    orderId,
		lastPrice:  lastPrice,
		lastQty:    lastQty,
		execId:     execId,
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

func (om *orderManagerImpl) CancelOrder(params *api.CancelOrderParams) error {

	om.log.Print(params.OrderId + ":cancelling order")

	resultChan := make(chan errorCmdResult)

	om.cancelOrderChan <- cancelOrderCmd{
		Params:     params,
		ResultChan: resultChan,
	}

	result := <-resultChan

	om.log.Printf(params.OrderId+":cancel order result: %v", result)

	return result.Error
}

func (om *orderManagerImpl) executeUpdateTradedQntCmd(id string, lastPrice model.Decimal64, lastQty model.Decimal64,
	execId string, resultChan chan errorCmdResult) {

	order, exists := om.orderStore.GetOrder(id)
	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("update traded quantity failed, no order found for id %v", id)}
		return
	}

	err := order.AddExecution(lastPrice, lastQty, execId)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	err = om.orderStore.Store(order)
	resultChan <- errorCmdResult{Error: err}

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

func (om *orderManagerImpl) executeCancelOrderCmd(params *api.CancelOrderParams, resultChan chan errorCmdResult) {

	order, exists := om.orderStore.GetOrder(params.OrderId)
	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("cancel order failed, no order found for id %v", params.OrderId)}
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

	order := model.NewOrder(id, params.OrderSide, params.Quantity,
		params.Price, params.Listing.Id, params.OriginatorId, params.OriginatorRef)

	order.SetTargetStatus(model.OrderStatus_LIVE)

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
			OrderId: order.Id,
		},
		Error: err,
	}

}

type addExecutionCmd struct {
	orderId    string
	lastPrice  model.Decimal64
	lastQty    model.Decimal64
	execId     string
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
	Params     *api.CancelOrderParams
	ResultChan chan errorCmdResult
}

type errorCmdResult struct {
	Error error
}
