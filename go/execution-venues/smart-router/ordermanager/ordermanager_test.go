package ordermanager

import (
	"context"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/model"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"testing"
)

func Test_cancelOfUnexposedOrder(t *testing.T) {
	parentOrderUpdatesChan, _, _, _, _, _, om,_ := setupOrderManager()

	order := <-parentOrderUpdatesChan

	if order.GetTargetStatus() != model.OrderStatus_NONE || order.GetStatus() != model.OrderStatus_LIVE {
		t.FailNow()
	}

	if !order.GetExposedQuantity().Equal(model.IasD(0)) {
		t.Fatalf("parent order should not be exposed")
	}

	om.Cancel()

	order = <-parentOrderUpdatesChan

	if order.GetTargetStatus() != model.OrderStatus_NONE || order.GetStatus() != model.OrderStatus_CANCELLED {
		t.FailNow()
	}

	if !order.GetExposedQuantity().Equal(model.IasD(0)) {
		t.Fatalf("parent order should be only partly exposed")
	}

}

func Test_cancelOfPartiallyExposedOrder(t *testing.T) {
	parentOrderUpdatesChan, childOrderOutboundParams, cancelOrderOutboundParams, childOrdersIn, sendChildQty, listing, om, doneChan := setupOrderManager()

	params1 := &api.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity:      model.IasD(10),
		Price:         model.IasD(200),
		Listing:       listing,
		OriginatorId:  om.ExecVenueId,
		OriginatorRef: om.ManagedOrder.Id,
	}

	<-parentOrderUpdatesChan

	sendChildQty <- model.IasD(10)
	pd := <-childOrderOutboundParams

	if !areParamsEqual(params1, pd.params) {
		t.FailNow()
	}

	order := <-parentOrderUpdatesChan

	if order.GetTargetStatus() != model.OrderStatus_NONE || order.GetStatus() != model.OrderStatus_LIVE {
		t.FailNow()
	}

	if !order.GetExposedQuantity().Equal(model.IasD(10)) {
		t.Fatalf("parent order should be only partly exposed")
	}

	childOrdersIn <- &model.Order{
		Id:                pd.id,
		Version:           1,
		ListingId:         listing.Id,
		Status:            model.OrderStatus_LIVE,
		RemainingQuantity: model.IasD(10),
	}

	order = <-parentOrderUpdatesChan

	if order.GetTargetStatus() != model.OrderStatus_NONE || order.GetStatus() != model.OrderStatus_LIVE {
		t.FailNow()
	}

	if !order.GetExposedQuantity().Equal(model.IasD(10)) {
		t.Fatalf("parent order should be only partly exposed")
	}

	om.Cancel()

	cp := <-cancelOrderOutboundParams

	if cp.OrderId != pd.id {
		t.FailNow()
	}

	order = <-parentOrderUpdatesChan

	if order.GetTargetStatus() != model.OrderStatus_CANCELLED || order.GetStatus() != model.OrderStatus_LIVE {
		t.FailNow()
	}

	if !order.GetExposedQuantity().Equal(model.IasD(10)) {
		t.Fatalf("parent order should be only partly exposed")
	}

	childOrdersIn <- &model.Order{
		Id:                pd.id,
		Version:           2,
		ListingId:         listing.Id,
		Status:            model.OrderStatus_CANCELLED,
		RemainingQuantity: model.IasD(10),
	}

	order = <-parentOrderUpdatesChan

	if order.GetTargetStatus() != model.OrderStatus_NONE || order.GetStatus() != model.OrderStatus_CANCELLED {
		t.FailNow()
	}

	if !order.GetExposedQuantity().Equal(model.IasD(0)) {
		t.Fatalf("parent order should be not be exposed")
	}

	id := <- doneChan
	if id != om.ManagedOrder.Id {
		t.FailNow()
	}

}


func TestOrderManagerCompletesWhenChildOrdersFilled(t *testing.T) {


	log.Printf("here")
	parentOrderUpdatesChan, _, _, childOrdersIn, _,
	listing, om, doneChan, child1Id, child2Id := setupOrderManagerAndSendTwoChildOrders(t)


	childOrdersIn <- &model.Order{
		Id:                child1Id,
		Version:           2,
		Status:            model.OrderStatus_LIVE,
		LastExecQuantity:  model.IasD(60),
		LastExecPrice:     model.IasD(100),
		LastExecId:        "c1e1",
		RemainingQuantity: model.IasD(0),
		ListingId: listing.GetId(),
	}

	order := <-parentOrderUpdatesChan

	if !order.GetTradedQuantity().Equal(model.IasD(60)) {
		t.FailNow()
	}

	childOrdersIn <- &model.Order{
		Id:                child2Id,
		Version:           2,
		Status:            model.OrderStatus_LIVE,
		LastExecQuantity:  model.IasD(40),
		LastExecPrice:     model.IasD(110),
		LastExecId:        "c2e1",
		ListingId: listing.GetId(),
		RemainingQuantity: model.IasD(0),
	}

	order = <-parentOrderUpdatesChan

	if !order.GetTradedQuantity().Equal(model.IasD(100)) {
		t.FailNow()
	}

	if order.GetStatus() != model.OrderStatus_FILLED {
		t.FailNow()
	}

	doneId := <-doneChan

	if doneId != om.GetManagedOrderId() {
		t.FailNow()
	}

}


func setupOrderManagerAndSendTwoChildOrders(t *testing.T) (parentOrderUpdatesChan chan model.Order, childOrderOutboundParams chan paramsAndId,
	childOrderCancelParams chan *api.CancelOrderParams, childOrdersIn chan *model.Order,
	sendChildQty chan *model.Decimal64, listing *model.Listing,
	om *OrderManager,  doneChan chan string, child1Id string, child2Id string)  {

	parentOrderUpdatesChan, childOrderOutboundParams, childOrderCancelParams, childOrdersIn, sendChildQty, listing,
	om, doneChan = setupOrderManager()

	<-parentOrderUpdatesChan

	sendChildQty <- &model.Decimal64{Mantissa: 60}

	params1 := &api.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity:      model.IasD(60),
		Price:         model.IasD(200),
		Listing:       listing,
		OriginatorId:  om.ExecVenueId,
		OriginatorRef: om.ManagedOrder.Id,
	}

	pd := <-childOrderOutboundParams
	child1Id = pd.id

	if !areParamsEqual(params1, pd.params) {
		t.FailNow()
	}

	<-parentOrderUpdatesChan

	sendChildQty <- &model.Decimal64{Mantissa: 40}
	params2 := &api.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity:      model.IasD(40),
		Price:         model.IasD(200),
		Listing:       listing,
		OriginatorId:  om.ExecVenueId,
		OriginatorRef: om.ManagedOrder.Id,
	}

	pd = <-childOrderOutboundParams
	child2Id = pd.id

	if !areParamsEqual(params2, pd.params) {
		t.FailNow()
	}

	order := <-parentOrderUpdatesChan

	if order.GetTargetStatus() != model.OrderStatus_NONE || order.GetStatus() != model.OrderStatus_LIVE {
		t.FailNow()
	}

	if order.GetAvailableQty().GreaterThan(model.IasD(0)) {
		t.Fatalf("no quantity should be left to trade")
	}

	childOrdersIn <- &model.Order{
		Id:                child1Id,
		Version:           0,
		Status:            model.OrderStatus_LIVE,
		ListingId:         listing.GetId(),
		RemainingQuantity: model.IasD(60),
	}

	order = <-parentOrderUpdatesChan
	if !order.GetExposedQuantity().Equal(model.IasD(100)) {
		t.FailNow()
	}

	childOrdersIn <- &model.Order{
		Id:                child2Id,
		Version:           0,
		Status:            model.OrderStatus_LIVE,
		ListingId: listing.Id,
		RemainingQuantity: model.IasD(40),
	}

	order = <-parentOrderUpdatesChan
	if !order.GetExposedQuantity().Equal(model.IasD(100)) {
		t.FailNow()
	}

	return parentOrderUpdatesChan, childOrderOutboundParams, childOrderCancelParams, childOrdersIn, sendChildQty,
		listing, om, doneChan, child1Id, child2Id
}

func setupOrderManager() (parentOrderUpdatesChan chan model.Order, childOrderOutboundParams chan paramsAndId,
	childOrderCancelParams chan *api.CancelOrderParams, childOrdersIn chan *model.Order,
	sendChildQty chan *model.Decimal64, listing *model.Listing,
	om *OrderManager,  doneChan chan string) {

	listing = &model.Listing{
		Version: 0,
		Id:      1,
	}

	parentOrderUpdatesChan = make(chan model.Order)

	childOrderOutboundParams = make(chan paramsAndId)
	childOrderCancelParams = make(chan *api.CancelOrderParams)
	orderRouter := &testOmClient{
		croParamsChan:    childOrderOutboundParams,
		cancelParamsChan: childOrderCancelParams,
	}

	childOrdersIn = make(chan *model.Order)
	childOrderStream := testChildOrderStream{stream: childOrdersIn}

	doneChan = make(chan string)

	om, err := NewOrderManagerFromCreateParams("p1", &api.CreateAndRouteOrderParams{
		OrderSide:          model.Side_BUY,
		Quantity:           &model.Decimal64{Mantissa: 100},
		Price:              &model.Decimal64{Mantissa: 200},
		Listing:            listing,
		OriginatorId:       "",
		OriginatorRef:      "",
		RootOriginatorId:   "",
		RootOriginatorRef:  "",
		ExecParametersJson: "",
	}, "e1", func(o *model.Order) error {
		parentOrderUpdatesChan <- *o
		return nil
	}, orderRouter, childOrderStream, doneChan)
	if err != nil {
		panic(err)
	}

	sendChildQty = make(chan *model.Decimal64)
	ExecuteAsDmaOrderManager(om, sendChildQty, listing)
	return parentOrderUpdatesChan, childOrderOutboundParams, childOrderCancelParams, childOrdersIn, sendChildQty, listing, om, doneChan
}

func ExecuteAsDmaOrderManager(om *OrderManager, sendChildQty chan *model.Decimal64, listing *model.Listing) {

	if om.ManagedOrder.GetTargetStatus() == model.OrderStatus_LIVE {
		om.ManagedOrder.SetStatus(model.OrderStatus_LIVE)
	}
	go func() {
		for {
			done, err := om.CheckIfDone()
			if err != nil {
				om.ErrLog.Printf("failed to check if done, cancelling order:%v", err)
				om.Cancel()
			}

			if done {
				break
			}

			select {
			case <-om.CancelChan:
				err := om.CancelManagedOrder(func(listingId int32) *model.Listing {
					if listingId != listing.Id {
						panic("unexpected listing id")
					}
					return listing
				})
				if err != nil {
					log.Panicf("failed to cancel order:%v", err)
				}
			case co, ok := <-om.ChildOrderUpdateChan:
				om.OnChildOrderUpdate(ok, co)
			case q := <-sendChildQty:
				om.SendChildOrder(om.ManagedOrder.Side, q, om.ManagedOrder.Price, listing)
			}
		}
	}()

}

func areParamsEqual(p1 *api.CreateAndRouteOrderParams, p2 *api.CreateAndRouteOrderParams) bool {
	return p1.Quantity.Equal(p2.Quantity) && p1.Listing.Id == p2.Listing.Id && p1.Price.Equal(p2.Price) && p1.OrderSide == p2.OrderSide &&
		p1.OriginatorRef == p2.OriginatorRef && p1.OriginatorId == p2.OriginatorId

}

type testEvClient struct {
	params []*api.CreateAndRouteOrderParams
}

func (t *testEvClient) GetExecutionParametersMetaData(ctx context.Context, empty *model.Empty, opts ...grpc.CallOption) (*api.ExecParamsMetaDataJson, error) {
	panic("implement me")
}

func (t *testEvClient) CreateAndRouteOrder(ctx context.Context, in *api.CreateAndRouteOrderParams, opts ...grpc.CallOption) (*api.OrderId, error) {
	t.params = append(t.params, in)
	id, _ := uuid.NewUUID()
	return &api.OrderId{
		OrderId: id.String(),
	}, nil
}

type paramsAndId struct {
	params *api.CreateAndRouteOrderParams
	id     string
}

type testOmClient struct {
	croParamsChan    chan paramsAndId
	cancelParamsChan chan *api.CancelOrderParams
}

func (t *testOmClient) GetExecutionParametersMetaData(ctx context.Context, empty *model.Empty, opts ...grpc.CallOption) (*api.ExecParamsMetaDataJson, error) {
	panic("implement me")
}

func (t *testOmClient) CreateAndRouteOrder(ctx context.Context, in *api.CreateAndRouteOrderParams, opts ...grpc.CallOption) (*api.OrderId, error) {

	id, _ := uuid.NewUUID()

	t.croParamsChan <- paramsAndId{in, id.String()}

	return &api.OrderId{
		OrderId: id.String(),
	}, nil
}

func (t *testOmClient) CancelOrder(ctx context.Context, in *api.CancelOrderParams, opts ...grpc.CallOption) (*model.Empty, error) {
	t.cancelParamsChan <- in
	return &model.Empty{}, nil
}

func (t *testOmClient) ModifyOrder(ctx context.Context, in *api.ModifyOrderParams, opts ...grpc.CallOption) (*model.Empty, error) {
	panic("implement me")
}

func (t *testEvClient) CancelOrder(ctx context.Context, in *api.CancelOrderParams, opts ...grpc.CallOption) (*model.Empty, error) {
	panic("implement me")
}

func (t *testEvClient) ModifyOrder(ctx context.Context, in *api.ModifyOrderParams, opts ...grpc.CallOption) (*model.Empty, error) {
	panic("implement me")
}

type testChildOrderStream struct {
	stream chan *model.Order
}

func (t testChildOrderStream) GetStream() <-chan *model.Order {
	return t.stream
}

func (t testChildOrderStream) Close() {
}
