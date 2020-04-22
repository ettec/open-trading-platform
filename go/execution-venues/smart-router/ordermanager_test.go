package main

import (
	"context"
	"github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"reflect"
	"testing"
)

type testEvClient struct {
	params []*executionvenue.CreateAndRouteOrderParams
}

func (t *testEvClient) CreateAndRouteOrder(ctx context.Context, in *executionvenue.CreateAndRouteOrderParams, opts ...grpc.CallOption) (*executionvenue.OrderId, error) {
	t.params = append(t.params, in)
	id, _ := uuid.NewUUID()
	return &executionvenue.OrderId{
		OrderId: id.String(),
	}, nil
}

func (t *testEvClient) CancelOrder(ctx context.Context, in *executionvenue.CancelOrderParams, opts ...grpc.CallOption) (*model.Empty, error) {
	panic("implement me")
}

func Test_submitSellOrders(t *testing.T) {

	orderId := "a"
	evId := "testev"

	q := &model.ClobQuote{
		Bids: []*model.ClobLine{
			{Size: model.IasD(10), Price: model.IasD(150), ListingId: 1},
			{Size: model.IasD(10), Price: model.IasD(140), ListingId: 2},
			{Size: model.IasD(10), Price: model.IasD(130), ListingId: 1},
			{Size: model.IasD(10), Price: model.IasD(120), ListingId: 2},
			{Size: model.IasD(10), Price: model.IasD(110), ListingId: 1},
			{Size: model.IasD(10), Price: model.IasD(100), ListingId: 1},
		},
	}

	mo := newManagedOrder(model.NewOrder(orderId, model.Side_SELL, model.IasD(50), model.IasD(120), 0, "oi", "od"))

	listing1 := &model.Listing{Id: 1}
	listing2 := &model.Listing{Id: 2}
	underlyingListings := map[int32]*model.Listing{
		1: listing1,
		2: listing2,
	}

	client := &testEvClient{}
	submitSellOrders(q, mo, underlyingListings, evId, client)

	if len(client.params) != 5 {
		t.FailNow()
	}

	expectedParams := []*executionvenue.CreateAndRouteOrderParams{{
		OrderSide:     model.Side_SELL,
		Quantity:      model.IasD(10),
		Price:         model.IasD(150),
		Listing:       listing1,
		OriginatorId:  evId,
		OriginatorRef: orderId,
	},
		{
			OrderSide:     model.Side_SELL,
			Quantity:      model.IasD(10),
			Price:         model.IasD(140),
			Listing:       listing2,
			OriginatorId:  evId,
			OriginatorRef: orderId,
		},
		{
			OrderSide:     model.Side_SELL,
			Quantity:      model.IasD(10),
			Price:         model.IasD(130),
			Listing:       listing1,
			OriginatorId:  evId,
			OriginatorRef: orderId,
		},
		{
			OrderSide:     model.Side_SELL,
			Quantity:      model.IasD(10),
			Price:         model.IasD(120),
			Listing:       listing2,
			OriginatorId:  evId,
			OriginatorRef: orderId,
		},
		{
			OrderSide:     model.Side_SELL,
			Quantity:      model.IasD(10),
			Price:         model.IasD(120),
			Listing:       listing1,
			OriginatorId:  evId,
			OriginatorRef: orderId,
		},
	}

	for idx, params := range client.params {
		if !reflect.DeepEqual(expectedParams[idx], params) {
			t.Fatalf("expected params at idx %v do not match", idx)
		}
	}

}

func Test_submitBuyOrders(t *testing.T) {

	orderId := "a"
	evId := "testev"

	q := &model.ClobQuote{
		Offers: []*model.ClobLine{
			{Size: model.IasD(10), Price: model.IasD(100), ListingId: 1},
			{Size: model.IasD(10), Price: model.IasD(110), ListingId: 2},
			{Size: model.IasD(10), Price: model.IasD(120), ListingId: 1},
			{Size: model.IasD(10), Price: model.IasD(130), ListingId: 2},
			{Size: model.IasD(10), Price: model.IasD(140), ListingId: 1},
			{Size: model.IasD(10), Price: model.IasD(150), ListingId: 1},
		},
	}

	mo := newManagedOrder(model.NewOrder(orderId, model.Side_BUY, model.IasD(50), model.IasD(130), 0, "oi", "od"))

	listing1 := &model.Listing{Id: 1}
	listing2 := &model.Listing{Id: 2}
	underlyingListings := map[int32]*model.Listing{
		1: listing1,
		2: listing2,
	}

	client := &testEvClient{}
	submitBuyOrders(q, mo, underlyingListings, evId, client)

	if len(client.params) != 5 {
		t.FailNow()
	}

	expectedParams := []*executionvenue.CreateAndRouteOrderParams{{
		OrderSide:     model.Side_BUY,
		Quantity:      model.IasD(10),
		Price:         model.IasD(100),
		Listing:       listing1,
		OriginatorId:  evId,
		OriginatorRef: orderId,
	},
		{
			OrderSide:     model.Side_BUY,
			Quantity:      model.IasD(10),
			Price:         model.IasD(110),
			Listing:       listing2,
			OriginatorId:  evId,
			OriginatorRef: orderId,
		},
		{
			OrderSide:     model.Side_BUY,
			Quantity:      model.IasD(10),
			Price:         model.IasD(120),
			Listing:       listing1,
			OriginatorId:  evId,
			OriginatorRef: orderId,
		},
		{
			OrderSide:     model.Side_BUY,
			Quantity:      model.IasD(10),
			Price:         model.IasD(130),
			Listing:       listing2,
			OriginatorId:  evId,
			OriginatorRef: orderId,
		},
		{
			OrderSide:     model.Side_BUY,
			Quantity:      model.IasD(10),
			Price:         model.IasD(130),
			Listing:       listing1,
			OriginatorId:  evId,
			OriginatorRef: orderId,
		},
	}

	for idx, params := range client.params {
		if !reflect.DeepEqual(expectedParams[idx], params) {
			t.Fatalf("expected params at idx %v do not match", idx)
		}
	}

}

type paramsAndId struct {
	params *executionvenue.CreateAndRouteOrderParams
	id     string
}

type testOmClient struct {
	params chan paramsAndId
}

func (t *testOmClient) CreateAndRouteOrder(ctx context.Context, in *executionvenue.CreateAndRouteOrderParams, opts ...grpc.CallOption) (*executionvenue.OrderId, error) {

	id, _ := uuid.NewUUID()

	t.params <- paramsAndId{in, id.String()}

	return &executionvenue.OrderId{
		OrderId: id.String(),
	}, nil
}

func (t *testOmClient) CancelOrder(ctx context.Context, in *executionvenue.CancelOrderParams, opts ...grpc.CallOption) (*model.Empty, error) {
	panic("implement me")
}

func TestNewOrderManager(t *testing.T) {

	evId := "testev"

	q := &model.ClobQuote{
		Offers: []*model.ClobLine{
			{Size: model.IasD(10), Price: model.IasD(100), ListingId: 1},
			{Size: model.IasD(10), Price: model.IasD(110), ListingId: 2},
			{Size: model.IasD(10), Price: model.IasD(120), ListingId: 1},
			{Size: model.IasD(10), Price: model.IasD(130), ListingId: 2},
			{Size: model.IasD(10), Price: model.IasD(140), ListingId: 1},
			{Size: model.IasD(10), Price: model.IasD(150), ListingId: 1},
		},
	}

	srListing := &model.Listing{Id: 3}

	listing1 := &model.Listing{Id: 1}
	listing2 := &model.Listing{Id: 2}
	underlyingListings := map[int32]*model.Listing{
		1: listing1,
		2: listing2,
	}

	params := &executionvenue.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity:      model.IasD(20),
		Price:         model.IasD(130),
		Listing:       srListing,
		OriginatorId:  "oi",
		OriginatorRef: "or",
	}

	done := make(chan string)
	quoteChan := make(chan *model.ClobQuote)
	childOrderUpdates := make(chan *model.Order)

	orderUpdates := make(chan model.Order)

	paramsChan := make(chan paramsAndId)

	_, err := NewOrderManager(params, evId, underlyingListings, done, func(o model.Order) error {
		orderUpdates <- o
		return nil
	}, &testOmClient{paramsChan}, quoteChan, childOrderUpdates)

	if err != nil {
		t.Fatal(err)
	}

	order := <-orderUpdates

	if order.GetTargetStatus() != model.OrderStatus_LIVE {
		t.FailNow()
	}

	quoteChan <- q

	params1 := &executionvenue.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity:      model.IasD(10),
		Price:         model.IasD(100),
		Listing:       listing1,
		OriginatorId:  evId,
		OriginatorRef: order.Id,
	}

	pd := <-paramsChan
	child1Id := pd.id

	if !reflect.DeepEqual(params1, pd.params) {
		t.FailNow()
	}

	params2 := &executionvenue.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity:      model.IasD(10),
		Price:         model.IasD(110),
		Listing:       listing2,
		OriginatorId:  evId,
		OriginatorRef: order.Id,
	}

	pd = <-paramsChan
	//child2Id := pd.id

	if !reflect.DeepEqual(params2, pd.params) {
		t.FailNow()
	}

	order = <-orderUpdates

	if order.GetTargetStatus() != model.OrderStatus_NONE || order.GetStatus() != model.OrderStatus_LIVE {
		t.FailNow()
	}

	if order.GetAvailableQty().GreaterThan(model.IasD(0)) {
		t.Fatalf("no quantity should be left to trade")
	}

	childOrderUpdates <- &model.Order{
		Id:     child1Id,
		Status: model.OrderStatus_LIVE,
	}

	childOrderUpdates <- &model.Order{
		Id:               child1Id,
		Status:           model.OrderStatus_LIVE,
		LastExecQuantity: model.IasD(10),
		LastExecPrice:    model.IasD(100),
		LastExecSeqNo:    1,
		LastExecId:       "e1",
	}

	order = <-orderUpdates

	if !order.GetTradedQuantity().Equal(model.IasD(10)) {
		t.FailNow()
	}

}

func IasD(i int) *model.Decimal64 {
	return model.IasD(i)
}

func Test_managedOrder_onChildOrderUpdate(t *testing.T) {

	mo := newManagedOrder(model.NewOrder("a", model.Side_BUY, IasD(20), IasD(50), 1, "oi", "or"))
	mo.setStatus(model.OrderStatus_LIVE)

	mo.onChildOrderUpdate(&model.Order{Id: "a1", TargetStatus: model.OrderStatus_LIVE, Quantity: IasD(15), RemainingQuantity: IasD(15)})

	if !mo.order.ExposedQuantity.Equal(IasD(15)) {
		t.FailNow()
	}

	if !mo.order.GetAvailableQty().Equal(IasD(5)) {
		t.FailNow()
	}

	mo.onChildOrderUpdate(&model.Order{Id: "a2", TargetStatus: model.OrderStatus_LIVE, Quantity: IasD(5), RemainingQuantity: IasD(5)})

	if !mo.order.ExposedQuantity.Equal(IasD(20)) {
		t.FailNow()
	}

	if !mo.order.GetAvailableQty().Equal(IasD(0)) {
		t.FailNow()
	}

	mo.onChildOrderUpdate(&model.Order{Id: "a1", LastExecPrice: IasD(50), LastExecQuantity: IasD(5), LastExecSeqNo: 1,
		LastExecId: "e1", RemainingQuantity: IasD(10)})

	if !mo.order.TradedQuantity.Equal(IasD(5)) {
		t.FailNow()
	}

	if !mo.order.ExposedQuantity.Equal(IasD(15)) {
		t.FailNow()
	}

	mo.onChildOrderUpdate(&model.Order{Id: "a2", TargetStatus: model.OrderStatus_LIVE, Quantity: IasD(5), RemainingQuantity: IasD(0),
		LastExecPrice: IasD(50), LastExecQuantity: IasD(5), LastExecSeqNo: 1,
		LastExecId: "e1"})

	if !mo.order.TradedQuantity.Equal(IasD(10)) {
		t.FailNow()
	}

	if !mo.order.ExposedQuantity.Equal(IasD(10)) {
		t.FailNow()
	}

	if !mo.order.GetAvailableQty().Equal(IasD(0)) {
		t.FailNow()
	}

	mo.onChildOrderUpdate(&model.Order{Id: "a1", LastExecPrice: IasD(50), LastExecQuantity: IasD(10), LastExecSeqNo: 2,
		LastExecId: "e2", RemainingQuantity: IasD(0)})

	if !mo.order.TradedQuantity.Equal(IasD(20)) {
		t.FailNow()
	}

	if !mo.order.ExposedQuantity.Equal(IasD(0)) {
		t.FailNow()
	}

	if !mo.order.GetAvailableQty().Equal(IasD(0)) {
		t.FailNow()
	}

	if mo.order.GetStatus() != model.OrderStatus_FILLED {
		t.FailNow()
	}

}
