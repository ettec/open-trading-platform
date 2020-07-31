package main

import (
	"context"
	"github.com/ettec/open-trading-platform/go/smart-router/ordermanager"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/model"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"reflect"
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

	mo := model.NewOrder(orderId, model.Side_SELL, model.IasD(50), model.IasD(120), 0, "oi", "od",
		"ri", "rr")

	listing1 := &model.Listing{Id: 1}
	listing2 := &model.Listing{Id: 2}
	underlyingListings := map[int32]*model.Listing{
		1: listing1,
		2: listing2,
	}

	client := &testEvClient{}

	om := ordermanager.NewOrderManagerFromState(mo, func(order *model.Order) error {
		return nil
	}, evId, client, testChildOrderStream{}, make(chan string))

	submitSellOrders(om, q, underlyingListings)

	if len(client.params) != 4 {
		t.FailNow()
	}

	expectedParams := []*api.CreateAndRouteOrderParams{{
		OrderSide:         model.Side_SELL,
		Quantity:          model.IasD(10),
		Price:             model.IasD(150),
		Listing:           listing1,
		OriginatorId:      evId,
		OriginatorRef:     orderId,
		RootOriginatorId:  "ri",
		RootOriginatorRef: "rr",
	},
		{
			OrderSide:         model.Side_SELL,
			Quantity:          model.IasD(10),
			Price:             model.IasD(140),
			Listing:           listing2,
			OriginatorId:      evId,
			OriginatorRef:     orderId,
			RootOriginatorId:  "ri",
			RootOriginatorRef: "rr",
		},
		{
			OrderSide:         model.Side_SELL,
			Quantity:          model.IasD(10),
			Price:             model.IasD(130),
			Listing:           listing1,
			OriginatorId:      evId,
			OriginatorRef:     orderId,
			RootOriginatorId:  "ri",
			RootOriginatorRef: "rr",
		},
		{
			OrderSide:         model.Side_SELL,
			Quantity:          model.IasD(10),
			Price:             model.IasD(120),
			Listing:           listing2,
			OriginatorId:      evId,
			OriginatorRef:     orderId,
			RootOriginatorId:  "ri",
			RootOriginatorRef: "rr",
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

	listing1 := &model.Listing{Id: 1, Market: &model.Market{Mic: "XNAS"}}
	listing2 := &model.Listing{Id: 2, Market: &model.Market{Mic: "XNAS"}}
	underlyingListings := map[int32]*model.Listing{
		1: listing1,
		2: listing2,
	}

	client := &testEvClient{}
	om := ordermanager.NewOrderManagerFromState(model.NewOrder(orderId, model.Side_BUY, model.IasD(50), model.IasD(130), 0,
		"oi", "od", "ri", "rr"), func(order *model.Order) error {
		return nil
	}, evId, client, testChildOrderStream{}, make(chan string))

	submitBuyOrders(om, q, underlyingListings)

	if len(client.params) != 4 {
		t.FailNow()
	}

	expectedParams := []*api.CreateAndRouteOrderParams{{
		OrderSide:         model.Side_BUY,
		Quantity:          model.IasD(10),
		Price:             model.IasD(100),
		Listing:           listing1,
		OriginatorId:      evId,
		OriginatorRef:     orderId,
		RootOriginatorId:  "ri",
		RootOriginatorRef: "rr",
	},
		{
			OrderSide:         model.Side_BUY,
			Quantity:          model.IasD(10),
			Price:             model.IasD(110),
			Listing:           listing2,
			OriginatorId:      evId,
			OriginatorRef:     orderId,
			RootOriginatorId:  "ri",
			RootOriginatorRef: "rr",
		},
		{
			OrderSide:         model.Side_BUY,
			Quantity:          model.IasD(10),
			Price:             model.IasD(120),
			Listing:           listing1,
			OriginatorId:      evId,
			OriginatorRef:     orderId,
			RootOriginatorId:  "ri",
			RootOriginatorRef: "rr",
		},
		{
			OrderSide:         model.Side_BUY,
			Quantity:          model.IasD(10),
			Price:             model.IasD(130),
			Listing:           listing2,
			OriginatorId:      evId,
			OriginatorRef:     orderId,
			RootOriginatorId:  "ri",
			RootOriginatorRef: "rr",
		},
	}

	for idx, params := range client.params {
		if !reflect.DeepEqual(expectedParams[idx], params) {
			t.Fatalf("expected croParamsChan at idx %v do not match", idx)
		}
	}

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

func TestOrderManagerSendsChildOrders(t *testing.T) {
	setupOrderManagerAndSendTwoChildOrders(t)
}

func TestOrderManagerCancel(t *testing.T) {

	done, childOrderUpdates, orderUpdates, om, order, child1Id, child2Id, testExecVenue := setupOrderManagerAndSendTwoChildOrders(t)

	om.Cancel()

	cp1 := <-testExecVenue.cancelParamsChan

	if cp1.OrderId != child1Id && cp1.OrderId != child2Id {
		t.FailNow()
	}

	cp2 := <-testExecVenue.cancelParamsChan

	if cp2.OrderId != child2Id && cp1.OrderId != child2Id {
		t.FailNow()
	}

	update := <-orderUpdates
	if update.GetTargetStatus() != model.OrderStatus_CANCELLED {
		t.FailNow()
	}

	childOrderUpdates <- &model.Order{
		Id:                child1Id,
		Version:           2,
		Status:            model.OrderStatus_CANCELLED,
		RemainingQuantity: IasD(10),
	}

	update = <-orderUpdates
	if !update.GetExposedQuantity().Equal(model.IasD(10)) {
		t.FailNow()
	}

	childOrderUpdates <- &model.Order{
		Id:                child2Id,
		Version:           2,
		Status:            model.OrderStatus_CANCELLED,
		RemainingQuantity: IasD(10),
	}

	update = <-orderUpdates
	if !update.GetExposedQuantity().Equal(model.IasD(0)) {
		t.FailNow()
	}

	if update.GetStatus() != model.OrderStatus_CANCELLED {
		t.FailNow()
	}

	id := <-done
	if id != order.Id {
		t.FailNow()
	}

}

func TestOrderManagerCompletesWhenChildOrdersFilled(t *testing.T) {

	done, childOrderUpdates, orderUpdates, om, order, child1Id, child2Id, _ := setupOrderManagerAndSendTwoChildOrders(t)

	childOrderUpdates <- &model.Order{
		Id:                child1Id,
		Version:           2,
		Status:            model.OrderStatus_LIVE,
		LastExecQuantity:  model.IasD(10),
		LastExecPrice:     model.IasD(100),
		LastExecId:        "c1e1",
		RemainingQuantity: IasD(0),
	}

	order = <-orderUpdates

	if !order.GetTradedQuantity().Equal(model.IasD(10)) {
		t.FailNow()
	}

	childOrderUpdates <- &model.Order{
		Id:                child2Id,
		Version:           2,
		Status:            model.OrderStatus_LIVE,
		LastExecQuantity:  model.IasD(10),
		LastExecPrice:     model.IasD(110),
		LastExecId:        "c2e1",
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

type testQuoteStream struct {
	stream chan *model.ClobQuote
}

func (t testQuoteStream) Subscribe(listingId int32) {

}

func (t testQuoteStream) GetStream() <-chan *model.ClobQuote {
	return t.stream
}

func (t testQuoteStream) Close() {

}

type testChildOrderStream struct {
	stream chan *model.Order
}

func (t testChildOrderStream) GetStream() <-chan *model.Order {
	return t.stream
}

func (t testChildOrderStream) Close() {
}




func Test_smartRouterSubmitsOrderWhenLiquidityBecomesAvailable(t *testing.T) {
	evId, listing1, listing2, _, quoteChan, _, orderUpdates, paramsChan, _, _, order, _ := setupOrderManager(t)

	q := &model.ClobQuote{
		Offers: []*model.ClobLine{
			{Size: model.IasD(10), Price: model.IasD(100), ListingId: 1},
		},
	}

	quoteChan <- q

	params1 := &api.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity:      model.IasD(10),
		Price:         model.IasD(100),
		Listing:       listing1,
		OriginatorId:  evId,
		OriginatorRef: order.Id,
	}

	pd := <-paramsChan

	if !areParamsEqual(params1, pd.params) {
		t.FailNow()
	}

	q = &model.ClobQuote{
		Offers: []*model.ClobLine{
			{Size: model.IasD(10), Price: model.IasD(110), ListingId: 2},
		},
	}

	order = <-orderUpdates

	if order.GetTargetStatus() != model.OrderStatus_NONE || order.GetStatus() != model.OrderStatus_LIVE {
		t.FailNow()
	}

	if !order.GetExposedQuantity().Equal(model.IasD(10)) {
		t.Fatalf("parent order should be only partly exposed")
	}

	quoteChan <- q

	params2 := &api.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity:      model.IasD(10),
		Price:         model.IasD(110),
		Listing:       listing2,
		OriginatorId:  evId,
		OriginatorRef: order.Id,
	}

	pd = <-paramsChan

	if !areParamsEqual(params2, pd.params) {
		t.FailNow()
	}

	order = <-orderUpdates

	if !order.GetExposedQuantity().Equal(model.IasD(20)) {
		t.Fatalf("parent order should be only partly exposed")
	}

}

func setupOrderManagerAndSendTwoChildOrders(t *testing.T) (chan string, chan *model.Order, chan model.Order, *ordermanager.OrderManager, model.Order, string, string,
	*testOmClient) {
	evId, listing1, listing2, done, quoteChan, childOrderUpdates, orderUpdates, paramsChan, testExecVenue, om, order, _ := setupOrderManager(t)

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

	quoteChan <- q

	params1 := &api.CreateAndRouteOrderParams{
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

	params2 := &api.CreateAndRouteOrderParams{
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
		Version:           0,
		Status:            model.OrderStatus_LIVE,
		RemainingQuantity: IasD(10),
	}

	order = <-orderUpdates
	if !order.GetExposedQuantity().Equal(model.IasD(20)) {
		t.FailNow()
	}

	childOrderUpdates <- &model.Order{
		Id:                child2Id,
		Version:           0,
		Status:            model.OrderStatus_LIVE,
		RemainingQuantity: IasD(10),
	}

	order = <-orderUpdates
	if !order.GetExposedQuantity().Equal(model.IasD(20)) {
		t.FailNow()
	}
	return done, childOrderUpdates, orderUpdates, om, order, child1Id, child2Id, testExecVenue
}

func setupOrderManager(t *testing.T) (string, *model.Listing, *model.Listing, chan string, chan *model.ClobQuote,
	chan *model.Order, chan model.Order, chan paramsAndId, *testOmClient, *ordermanager.OrderManager, model.Order,
	chan *api.CancelOrderParams) {
	evId := "testev"

	srListing := &model.Listing{Id: 3}

	listing1 := &model.Listing{Id: 1, Market: &model.Market{Mic: "XNAS"}}
	listing2 := &model.Listing{Id: 2, Market: &model.Market{Mic: "XNAS"}}
	underlyingListings := []*model.Listing{
		listing1,
		listing2,
	}

	done := make(chan string)
	quoteChan := make(chan *model.ClobQuote)
	childOrderUpdates := make(chan *model.Order)

	orderUpdates := make(chan model.Order)

	paramsChan := make(chan paramsAndId)
	cancelParamsChan := make(chan *api.CancelOrderParams)

	testExecVenue := &testOmClient{paramsChan, cancelParamsChan}

	params := &api.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity:      model.IasD(20),
		Price:         model.IasD(130),
		Listing:       srListing,
		OriginatorId:  "oi",
		OriginatorRef: "or",
	}

	uniqueId, err := uuid.NewUUID()

	if err != nil {
		t.FailNow()
	}

	om, err := ordermanager.NewOrderManagerFromCreateParams(uniqueId.String(), params, evId, func(o *model.Order) error {
		orderUpdates <- *o
		return nil
	}, testExecVenue, &testChildOrderStream{childOrderUpdates}, done)

	if err != nil {
		t.Fatal(err)
	}

	ExecuteAsSmartRouterStrategy(om, func(listingId int32, listingGroupsIn chan<- []*model.Listing) {
		go func() {
			listingGroupsIn <- underlyingListings
		}()
	}, testQuoteStream{stream: quoteChan})

	order := <-orderUpdates

	if order.GetTargetStatus() != model.OrderStatus_LIVE {
		t.FailNow()
	}
	return evId, listing1, listing2, done, quoteChan, childOrderUpdates, orderUpdates, paramsChan, testExecVenue, om, order, cancelParamsChan
}

func areParamsEqual(p1 *api.CreateAndRouteOrderParams, p2 *api.CreateAndRouteOrderParams) bool {
	return p1.Quantity.Equal(p2.Quantity) && p1.Listing.Id == p2.Listing.Id && p1.Price.Equal(p2.Price) && p1.OrderSide == p2.OrderSide &&
		p1.OriginatorRef == p2.OriginatorRef && p1.OriginatorId == p2.OriginatorId

}

func IasD(i int) *model.Decimal64 {
	return model.IasD(i)
}


