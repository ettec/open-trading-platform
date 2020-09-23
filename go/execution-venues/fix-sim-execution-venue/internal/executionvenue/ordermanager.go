package executionvenue

import (
	"fmt"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/executionvenue"
	"github.com/ettec/otp-common/model"
	"github.com/google/uuid"
	"log"
	"os"
)



type orderGateway interface {
	Send(order *model.Order, listing *model.Listing) error
	Cancel(order *model.Order) error
	Modify(order *model.Order, listing *model.Listing, Quantity *model.Decimal64, Price *model.Decimal64) error
}

type orderManagerImpl struct {
	createOrderChan    chan createAndRouteOrderCmd
	cancelOrderChan    chan cancelOrderCmd
	modifyOrderChan    chan modifyOrderCmd
	setOrderStatusChan chan setOrderStatusCmd
	setOrderErrMsgChan chan setOrderErrorMsgCmd
	addExecChan        chan addExecutionCmd

	closeChan chan struct{}

	orderStore  *executionvenue.OrderCache
	gateway     orderGateway
	getListing  func(listingId int32, result chan<- *model.Listing)
	log         *log.Logger
	errLog      *log.Logger
}

func NewOrderManager(cache *executionvenue.OrderCache, gateway orderGateway,
	getListing func(listingId int32, result chan<- *model.Listing)) *orderManagerImpl {

	om := &orderManagerImpl{
		log:         log.New(os.Stdout, "", log.Lshortfile|log.Ltime),
		errLog:      log.New(os.Stderr, "", log.Lshortfile|log.Ltime),
		getListing:  getListing,
	}

	om.createOrderChan = make(chan createAndRouteOrderCmd, 100)
	om.cancelOrderChan = make(chan cancelOrderCmd, 100)
	om.modifyOrderChan = make(chan modifyOrderCmd, 100)
	om.setOrderStatusChan = make(chan setOrderStatusCmd, 100)
	om.setOrderErrMsgChan = make(chan setOrderErrorMsgCmd, 100)
	om.addExecChan = make(chan addExecutionCmd, 100)

	om.closeChan = make(chan struct{}, 1)

	om.orderStore = cache
	om.gateway = gateway

	go om.executeOrderCommands()

	return om
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
			case mp := <-om.modifyOrderChan:
				om.executeModifyOrderCmd(mp.Params, mp.ResultChan)
			case cro := <-om.createOrderChan:
				om.executeCreateAndRouteOrderCmd(cro.Params, cro.ResultChan)
			case su := <-om.setOrderStatusChan:
				om.executeSetOrderStatusCmd(su.orderId, su.status, su.ResultChan)
			case em := <-om.setOrderErrMsgChan:
				om.executeSetErrorMsg(em.orderId, em.msg, em.ResultChan)
			case tu := <-om.addExecChan:
				om.executeUpdateTradedQntCmd(tu.orderId, tu.lastPrice, tu.lastQty, tu.execId, tu.ResultChan)
			}
		}
	}

}

func (om *orderManagerImpl) Close() {
	om.closeChan <- struct{}{}
}

func (om *orderManagerImpl) SetErrorMsg(orderId string, msg string) error {
	om.log.Printf("updating order %v error message to %v", orderId, msg)

	resultChan := make(chan errorCmdResult)

	om.setOrderErrMsgChan <- setOrderErrorMsgCmd{
		orderId:    orderId,
		msg:        msg,
		ResultChan: resultChan,
	}

	result := <-resultChan

	return result.Error
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

func (om *orderManagerImpl) ModifyOrder(params *api.ModifyOrderParams) error {
	om.log.Printf("modifying order %v, price %v, quantity %v", params.OrderId, params.Price, params.Quantity)
	resultChan := make(chan errorCmdResult)

	om.modifyOrderChan <- modifyOrderCmd{
		Params:     params,
		ResultChan: resultChan,
	}

	result := <-resultChan

	om.log.Printf(params.OrderId+":modify order result: %v", result)

	return result.Error
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

	err := order.AddExecution(model.Execution{
		Id:    execId,
		Price: lastPrice,
		Qty:   lastQty,
	})

	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	err = om.orderStore.Store(order)
	resultChan <- errorCmdResult{Error: err}

}

func (om *orderManagerImpl) executeSetErrorMsg(id string, msg string, resultChan chan errorCmdResult) {

	order, exists := om.orderStore.GetOrder(id)
	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("set order error message failed, no order found for id %v", id)}
		return
	}

	order.ErrorMessage = msg

	err := om.orderStore.Store(order)
	resultChan <- errorCmdResult{Error: err}

}

func (om *orderManagerImpl) executeSetOrderStatusCmd(id string, status model.OrderStatus, resultChan chan errorCmdResult) {

	order, exists := om.orderStore.GetOrder(id)
	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("set order status failed, no order found for id %v", id)}
		return
	}


	oldStatus := order.GetStatus()
	err := order.SetStatus(status)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	if order.Status != oldStatus {
		err = om.orderStore.Store(order)
	}

	resultChan <- errorCmdResult{Error: err}

}

func (om *orderManagerImpl) executeModifyOrderCmd(params *api.ModifyOrderParams, resultChan chan errorCmdResult) {

	order, exists := om.orderStore.GetOrder(params.OrderId)
	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("modify order failed, no order found for id %v", params.OrderId)}
		return
	}

	err := order.SetTargetStatus(model.OrderStatus_LIVE)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	order.Price = params.Price
	order.Quantity = params.Quantity

	err = om.orderStore.Store(order)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	listingChan := make(chan *model.Listing, 1)
	om.getListing(params.ListingId, listingChan)
	listing := <-listingChan

	err = om.gateway.Modify(order, listing, params.Quantity, params.Price)

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
		params.Price, params.ListingId, params.OriginatorId, params.OriginatorRef,
		params.RootOriginatorId, params.RootOriginatorRef, params.Destination)


	err = order.SetTargetStatus(model.OrderStatus_LIVE)
	if err != nil {
		om.errLog.Printf("failed to set target status;%v", err)
	}

	err = om.orderStore.Store(order)
	if err != nil {
		resultChan <- createAndRouteOrderCmdResult{
			OrderId: nil,
			Error:   err,
		}

		return
	}

	listingChan := make(chan *model.Listing, 1)
	om.getListing(params.ListingId, listingChan)

	listing := <-listingChan
	err = om.gateway.Send(order, listing)

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

type setOrderErrorMsgCmd struct {
	orderId    string
	msg        string
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

type modifyOrderCmd struct {
	Params     *api.ModifyOrderParams
	ResultChan chan errorCmdResult
}

type errorCmdResult struct {
	Error error
}
