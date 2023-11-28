package executionvenue

import (
	"context"
	"fmt"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/ordermanagement"
	"github.com/ettec/otp-common/staticdata"
	"github.com/google/uuid"
	"log/slog"
)

type orderGateway interface {
	Send(order *model.Order, listing *model.Listing) error
	Cancel(order *model.Order) error
	Modify(order *model.Order, listing *model.Listing, Quantity *model.Decimal64, Price *model.Decimal64) error
}

// orderManager is responsible for the creation, modification and cancellation of orders.  It depends on two resources
// that are single threaded and IO bound, the order cache and the order gateway. Therefore it is effectively single
// threaded.  To increase throughput additional instances of the execution venue should be deployed.  The order manager
// uses channels to queue commands, primarily to ensure that cancel commands are prioritised above all others and
// additionally to ensure that commands are executed fairly, i.e. in the order they are received, across clients
type orderManagerImpl struct {
	createOrderChan    chan createAndRouteOrderCmd
	cancelOrderChan    chan cancelOrderCmd
	modifyOrderChan    chan modifyOrderCmd
	setOrderStatusChan chan setOrderStatusCmd
	setOrderErrMsgChan chan setOrderErrorMsgCmd
	addExecChan        chan addExecutionCmd

	orderStore *ordermanagement.OrderCache
	gateway    orderGateway
	getListing func(ctx context.Context, listingId int32, result chan<- staticdata.ListingResult)
}

func NewOrderManager(ctx context.Context, cache *ordermanagement.OrderCache, gateway orderGateway,
	getListing func(ctx context.Context, listingId int32, result chan<- staticdata.ListingResult),
	cmdBufferSize int) *orderManagerImpl {

	om := &orderManagerImpl{
		getListing: getListing,
	}

	om.createOrderChan = make(chan createAndRouteOrderCmd, cmdBufferSize)
	om.cancelOrderChan = make(chan cancelOrderCmd, cmdBufferSize)
	om.modifyOrderChan = make(chan modifyOrderCmd, cmdBufferSize)
	om.setOrderStatusChan = make(chan setOrderStatusCmd, cmdBufferSize)
	om.setOrderErrMsgChan = make(chan setOrderErrorMsgCmd, cmdBufferSize)
	om.addExecChan = make(chan addExecutionCmd, cmdBufferSize)

	om.orderStore = cache
	om.gateway = gateway

	go om.executeOrderCommands(ctx)

	return om
}

func (om *orderManagerImpl) executeOrderCommands(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		// Cancel Requests take priority over all other message types
		case oc := <-om.cancelOrderChan:
			om.executeCancelOrderCmd(ctx, oc.Params, oc.ResultChan)
		default:
			select {
			case <-ctx.Done():
				return
			case oc := <-om.cancelOrderChan:
				om.executeCancelOrderCmd(ctx, oc.Params, oc.ResultChan)
			case mp := <-om.modifyOrderChan:
				om.executeModifyOrderCmd(ctx, mp.Params, mp.ResultChan)
			case cro := <-om.createOrderChan:
				om.executeCreateAndRouteOrderCmd(ctx, cro.Params, cro.ResultChan)
			case su := <-om.setOrderStatusChan:
				om.executeSetOrderStatusCmd(ctx, su.orderId, su.status, su.ResultChan)
			case em := <-om.setOrderErrMsgChan:
				om.executeSetErrorMsg(ctx, em.orderId, em.msg, em.ResultChan)
			case tu := <-om.addExecChan:
				om.executeUpdateTradedQntCmd(ctx, tu.orderId, tu.lastPrice, tu.lastQty, tu.execId, tu.ResultChan)
			}
		}
	}

}

func (om *orderManagerImpl) SetErrorMsg(orderId string, msg string) error {
	slog.Info("updating order error message", "orderId", orderId, "newMsg", msg)

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
	slog.Info("updating order status", "orderId", orderId, "newStatus", status)

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
	slog.Info("adding execution to order", "orderId", orderId, "price", lastPrice, "quantity", lastQty)

	resultChan := make(chan errorCmdResult)

	om.addExecChan <- addExecutionCmd{
		orderId:    orderId,
		lastPrice:  lastPrice,
		lastQty:    lastQty,
		execId:     execId,
		ResultChan: resultChan,
	}

	result := <-resultChan

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
	slog.Info("modifying order", "params", params)
	resultChan := make(chan errorCmdResult)

	om.modifyOrderChan <- modifyOrderCmd{
		Params:     params,
		ResultChan: resultChan,
	}

	result := <-resultChan

	return result.Error
}

func (om *orderManagerImpl) CancelOrder(params *api.CancelOrderParams) error {

	slog.Info("cancelling order", "params", params)

	resultChan := make(chan errorCmdResult)

	om.cancelOrderChan <- cancelOrderCmd{
		Params:     params,
		ResultChan: resultChan,
	}

	result := <-resultChan

	return result.Error
}

func (om *orderManagerImpl) executeUpdateTradedQntCmd(ctx context.Context, id string, lastPrice model.Decimal64,
	lastQty model.Decimal64,
	execId string, resultChan chan errorCmdResult) {

	order, exists, err := om.orderStore.GetOrder(id)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("update traded quantity failed, no order found for id %v", id)}
		return
	}

	err = order.AddExecution(model.Execution{
		Id:    execId,
		Price: lastPrice,
		Qty:   lastQty,
	})

	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	err = om.orderStore.Store(ctx, order)
	resultChan <- errorCmdResult{Error: err}

}

func (om *orderManagerImpl) executeSetErrorMsg(ctx context.Context, id string, msg string, resultChan chan errorCmdResult) {

	order, exists, err := om.orderStore.GetOrder(id)
	if err != nil {
		resultChan <- errorCmdResult{Error: fmt.Errorf("failed to get order for id %s from cache: %w", id, err)}
	}

	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("set order error message failed, no order found for id %s", id)}
		return
	}

	order.ErrorMessage = msg

	err = om.orderStore.Store(ctx, order)
}

func (om *orderManagerImpl) executeSetOrderStatusCmd(ctx context.Context, id string, status model.OrderStatus,
	resultChan chan errorCmdResult) {

	order, exists, err := om.orderStore.GetOrder(id)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("set order status failed, no order found for id %v", id)}
		return
	}

	err = order.SetStatus(status)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	err = om.orderStore.Store(ctx, order)

	resultChan <- errorCmdResult{Error: err}

}

func (om *orderManagerImpl) executeModifyOrderCmd(ctx context.Context, params *api.ModifyOrderParams, resultChan chan errorCmdResult) {

	order, exists, err := om.orderStore.GetOrder(params.OrderId)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("modify order failed, no order found for id %v", params.OrderId)}
		return
	}

	err = order.SetTargetStatus(model.OrderStatus_LIVE)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	order.Price = params.Price
	order.Quantity = params.Quantity

	err = om.orderStore.Store(ctx, order)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	listingChan := make(chan staticdata.ListingResult, 1)
	om.getListing(ctx, params.ListingId, listingChan)
	listingResult := <-listingChan
	if listingResult.Err != nil {
		resultChan <- errorCmdResult{Error: listingResult.Err}
		return
	}

	err = om.gateway.Modify(order, listingResult.Listing, params.Quantity, params.Price)

	resultChan <- errorCmdResult{Error: err}
}

func (om *orderManagerImpl) executeCancelOrderCmd(ctx context.Context,
	params *api.CancelOrderParams, resultChan chan errorCmdResult) {

	order, exists, err := om.orderStore.GetOrder(params.OrderId)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	if !exists {
		resultChan <- errorCmdResult{Error: fmt.Errorf("cancel order failed, no order found for id %s", params.OrderId)}
		return
	}

	err = order.SetTargetStatus(model.OrderStatus_CANCELLED)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	err = om.orderStore.Store(ctx, order)
	if err != nil {
		resultChan <- errorCmdResult{Error: err}
		return
	}

	err = om.gateway.Cancel(order)

	resultChan <- errorCmdResult{Error: err}
}

func (om *orderManagerImpl) executeCreateAndRouteOrderCmd(ctx context.Context, params *api.CreateAndRouteOrderParams,
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

	if err = order.SetTargetStatus(model.OrderStatus_LIVE); err != nil {
		slog.Error("failed to set target status to live", "error", err)
	}

	err = om.orderStore.Store(ctx, order)
	if err != nil {
		resultChan <- createAndRouteOrderCmdResult{
			OrderId: nil,
			Error:   err,
		}

		return
	}

	listingChan := make(chan staticdata.ListingResult, 1)
	om.getListing(ctx, params.ListingId, listingChan)

	listingResult := <-listingChan
	if listingResult.Err != nil {
		resultChan <- createAndRouteOrderCmdResult{
			OrderId: nil,
			Error:   listingResult.Err,
		}
		return
	}

	err = om.gateway.Send(order, listingResult.Listing)

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
