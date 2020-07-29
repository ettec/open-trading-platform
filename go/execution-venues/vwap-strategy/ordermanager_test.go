package main

import (
	"context"
	"encoding/json"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/model"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"testing"
)

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

func (t *testEvClient) CancelOrder(ctx context.Context, in *api.CancelOrderParams, opts ...grpc.CallOption) (*model.Empty, error) {
	panic("implement me")
}

func (t *testEvClient) ModifyOrder(ctx context.Context, in *api.ModifyOrderParams, opts ...grpc.CallOption) (*model.Empty, error) {
	panic("implement me")
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


type testChildOrderStream struct {
	stream chan *model.Order
}

func (t testChildOrderStream) GetStream() <-chan *model.Order {
	return t.stream
}

func (t testChildOrderStream) Close() {
}

func Test_submittedChildOrders(t *testing.T) {

}

func newVwapStrategy(t *testing.T, vwapParams vwapParameters, quantity *model.Decimal64,
	price *model.Decimal64, listing *model.Listing) (evId string, done chan string,
	childOrderUpdates chan *model.Order, parentOrderUpdates chan model.Order, paramsChan chan paramsAndId,
	testExecVenue *testOmClient, om *orderManager, parentOrder model.Order,
	cancelParamsChan chan *api.CancelOrderParams) {
	evId = "testev"


	done = make(chan string)
	childOrderUpdates = make(chan *model.Order)

	parentOrderUpdates = make(chan model.Order)

	paramsChan = make(chan paramsAndId)
	cancelParamsChan = make(chan *api.CancelOrderParams)

	testExecVenue = &testOmClient{paramsChan, cancelParamsChan}

	params := &api.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity:      quantity,
		Price:         price,
		Listing:       listing,
		OriginatorId:  "oi",
		OriginatorRef: "or",
	}

	uniqueId, err := uuid.NewUUID()

	if err != nil {
		t.FailNow()
	}

	paramsAsString, _ :=json.Marshal(vwapParams)
	buckets := getBucketsFromParamsString(string(paramsAsString), quantity, listing )

	om, err = NewOrderManagerFromParams(uniqueId.String(), params, evId, buckets, done, func(o *model.Order) error {
		parentOrderUpdates <- *o
		return nil
	}, testExecVenue, &testChildOrderStream{childOrderUpdates})

	if err != nil {
		t.Fatal(err)
	}

	parentOrder = <-parentOrderUpdates

	if parentOrder.GetTargetStatus() != model.OrderStatus_LIVE {
		t.FailNow()
	}
	return evId,  done, childOrderUpdates, parentOrderUpdates, paramsChan, testExecVenue, om, parentOrder, cancelParamsChan
}