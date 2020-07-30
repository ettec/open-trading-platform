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

func Test_orderManagerCancel(t *testing.T) {

	listing := &model.Listing{
		Version:              0,
		Id:                   1,
	}

	parentOrderUpdatesChan := make(chan model.Order)
	orderRouter := &testOmClient{}

	childOrdersIn := make(chan *model.Order)
	childOrderStream := testChildOrderStream{stream: childOrdersIn}

	doneChan := make(chan string)

	om, err := NewOrderManagerFromCreateParams("p1", &api.CreateAndRouteOrderParams{
		OrderSide:            model.Side_BUY,
		Quantity:             &model.Decimal64{Mantissa: 100},
		Price:                &model.Decimal64{Mantissa: 200},
		Listing:              listing,
		OriginatorId:         "",
		OriginatorRef:        "",
		RootOriginatorId:     "",
		RootOriginatorRef:    "",
		ExecParametersJson:   "",
	}, "e1", func(o *model.Order) error {
		parentOrderUpdatesChan <- *o
		return nil
	}, orderRouter, childOrderStream, doneChan  )
	if err != nil {
		panic(err)
	}

	TickSimpleOrder(om, listing)

	o := <-parentOrderUpdatesChan
	if o.GetStatus() != model.OrderStatus_LIVE {
		t.FailNow()
	}
}




func TickSimpleOrder(om *OrderManager, listing *model.Listing) bool {

	if om.ManagedOrder.GetTargetStatus() == model.OrderStatus_LIVE {
		om.ManagedOrder.SetStatus(model.OrderStatus_LIVE)
	}

	done, err := om.CheckIfDone()
	if err != nil {
		om.ErrLog.Printf("failed to check if done, cancelling order:%v", err)
		om.Cancel()
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

	}

	return done
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