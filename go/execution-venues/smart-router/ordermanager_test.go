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

	mo := newParentOrder(*model.NewOrder(orderId, model.Side_SELL, model.IasD(50), model.IasD(120), 0, "oi", "od"))

	listing1 := &model.Listing{Id: 1}
	listing2 := &model.Listing{Id: 2}
	underlyingListings := map[int32]*model.Listing{
		1: listing1,
		2: listing2,
	}

	client := &testEvClient{}

	om := newOrderManager(func(order model.Order) error {
		return nil
	}, mo, underlyingListings, evId, client)

	om.submitSellOrders(q)

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
			t.Fatalf("expected croParamsChan at idx %v do not match", idx)
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

	mo := newParentOrder(*model.NewOrder(orderId, model.Side_BUY, model.IasD(50), model.IasD(130), 0, "oi", "od"))

	listing1 := &model.Listing{Id: 1}
	listing2 := &model.Listing{Id: 2}
	underlyingListings := map[int32]*model.Listing{
		1: listing1,
		2: listing2,
	}

	client := &testEvClient{}
	om := newOrderManager(func(order model.Order) error {
		return nil
	}, mo, underlyingListings, evId, client)

	om.submitBuyOrders(q)

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
			t.Fatalf("expected croParamsChan at idx %v do not match", idx)
		}
	}

}

type paramsAndId struct {
	params *executionvenue.CreateAndRouteOrderParams
	id     string
}

type testOmClient struct {
	croParamsChan chan paramsAndId
	cancelParamsChan chan *executionvenue.CancelOrderParams
	
}

func (t *testOmClient) CreateAndRouteOrder(ctx context.Context, in *executionvenue.CreateAndRouteOrderParams, opts ...grpc.CallOption) (*executionvenue.OrderId, error) {

	id, _ := uuid.NewUUID()

	t.croParamsChan <- paramsAndId{in, id.String()}

	return &executionvenue.OrderId{
		OrderId: id.String(),
	}, nil
}

func (t *testOmClient) CancelOrder(ctx context.Context, in *executionvenue.CancelOrderParams, opts ...grpc.CallOption) (*model.Empty, error) {
	t.cancelParamsChan <- in
	return &model.Empty{},  nil
}

func TestOrderManagerSendsChildOrders(t *testing.T) {
	send2ChildOrders(t)
}

func TestOrderManagerCancel(t *testing.T) {

	done, childOrderUpdates, orderUpdates, om, order, child1Id, child2Id, testExecVenue := send2ChildOrders(t)

	om.Cancel()

	cp1 := <- testExecVenue.cancelParamsChan

	if cp1.OrderId != child1Id {
		t.FailNow()
	}

	cp2 := <- testExecVenue.cancelParamsChan

	if cp2.OrderId != child2Id {
		t.FailNow()
	}

	update := <-orderUpdates
	if update.GetTargetStatus() != model.OrderStatus_CANCELLED {
		t.FailNow()
	}

	childOrderUpdates <- &model.Order{
		Id:                child1Id,
		Status:            model.OrderStatus_CANCELLED,
		RemainingQuantity: IasD(10),
	}

	update = <-orderUpdates
	if !update.GetExposedQuantity().Equal(model.IasD(10)) {
		t.FailNow()
	}

	childOrderUpdates <- &model.Order{
		Id:                child2Id,
		Status:            model.OrderStatus_CANCELLED,
		RemainingQuantity: IasD(10),
	}

	update = <-orderUpdates
	if !update.GetExposedQuantity().Equal(model.IasD(0)) {
		t.FailNow()
	}

	update = <-orderUpdates
	if update.GetStatus() != model.OrderStatus_CANCELLED {
		t.FailNow()
	}
	
	id := <-done
	if id != order.Id {
		t.FailNow()
	}

}



func TestOrderManagerCompletesWhenChildOrdersFilled(t *testing.T) {

	done, childOrderUpdates, orderUpdates, om, order, child1Id, child2Id,_ := send2ChildOrders(t)

	childOrderUpdates <- &model.Order{
		Id:                child1Id,
		Status:            model.OrderStatus_LIVE,
		LastExecQuantity:  model.IasD(10),
		LastExecPrice:     model.IasD(100),
		LastExecSeqNo:     1,
		LastExecId:        "e1",
		RemainingQuantity: IasD(0),
	}

	order = <-orderUpdates

	if !order.GetTradedQuantity().Equal(model.IasD(10)) {
		t.FailNow()
	}

	childOrderUpdates <- &model.Order{
		Id:                child2Id,
		Status:            model.OrderStatus_LIVE,
		LastExecQuantity:  model.IasD(10),
		LastExecPrice:     model.IasD(110),
		LastExecSeqNo:     1,
		LastExecId:        "e1",
		RemainingQuantity: IasD(0),
	}

	order = <-orderUpdates

	if !order.GetTradedQuantity().Equal(model.IasD(20)) {
		t.FailNow()
	}

	if order.GetStatus() != model.OrderStatus_FILLED {
		t.FailNow()
	}

	doneId := <-done

	if doneId != om.GetManagedOrderId() {
		t.FailNow()
	}

}

func send2ChildOrders(t *testing.T) (chan string, chan *model.Order, chan model.Order, *orderManager, model.Order, string, string,
	*testOmClient) {
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
	cancelParamsChan := make(chan *executionvenue.CancelOrderParams)

	testExecVenue := &testOmClient{paramsChan, cancelParamsChan}

	om, err := NewOrderManager(params, evId, underlyingListings, done, func(o model.Order) error {
		orderUpdates <- o
		return nil
	}, testExecVenue, quoteChan, childOrderUpdates)

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

	if !areParamsEqual(params1, pd.params) {
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
	child2Id := pd.id

	if !areParamsEqual(params2, pd.params) {
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
		Id:                child1Id,
		Status:            model.OrderStatus_LIVE,
		RemainingQuantity: IasD(10),
	}

	order = <-orderUpdates
	if !order.GetExposedQuantity().Equal(model.IasD(20)) {
		t.FailNow()
	}

	childOrderUpdates <- &model.Order{
		Id:                child2Id,
		Status:            model.OrderStatus_LIVE,
		RemainingQuantity: IasD(10),
	}

	order = <-orderUpdates
	if !order.GetExposedQuantity().Equal(model.IasD(20)) {
		t.FailNow()
	}
	return done, childOrderUpdates, orderUpdates, om, order, child1Id, child2Id, testExecVenue
}

func areParamsEqual(p1 *executionvenue.CreateAndRouteOrderParams, p2 *executionvenue.CreateAndRouteOrderParams) bool {
	return p1.Quantity.Equal(p2.Quantity) && p1.Listing.Id == p2.Listing.Id && p1.Price.Equal(p2.Price) && p1.OrderSide == p2.OrderSide &&
		p1.OriginatorRef == p2.OriginatorRef && p1.OriginatorId == p2.OriginatorId

}

func IasD(i int) *model.Decimal64 {
	return model.IasD(i)
}

func Test_orderManager_Cancel(t *testing.T) {

}