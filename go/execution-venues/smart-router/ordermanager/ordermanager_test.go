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
	parentOrderUpdatesChan, _, _, _, _, _, om := setupOrderManager()

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
	parentOrderUpdatesChan, childOrderOutboundParams, cancelOrderOutboundParams, childOrdersIn, sendChildQty, listing, om := setupOrderManager()

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

}

func setupOrderManager() (chan model.Order, chan paramsAndId, chan *api.CancelOrderParams, chan *model.Order, chan *model.Decimal64, *model.Listing,
	*OrderManager) {
	listing := &model.Listing{
		Version: 0,
		Id:      1,
	}

	parentOrderUpdatesChan := make(chan model.Order)

	childOrderOutboundParams := make(chan paramsAndId)
	childOrderCancelParams := make(chan *api.CancelOrderParams)
	orderRouter := &testOmClient{
		croParamsChan:    childOrderOutboundParams,
		cancelParamsChan: childOrderCancelParams,
	}

	childOrdersIn := make(chan *model.Order)
	childOrderStream := testChildOrderStream{stream: childOrdersIn}

	doneChan := make(chan string)

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

	sendChildQty := make(chan *model.Decimal64)
	ExecuteAsDmaOrderManager(om, sendChildQty, listing)
	return parentOrderUpdatesChan, childOrderOutboundParams, childOrderCancelParams, childOrdersIn, sendChildQty, listing, om
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
