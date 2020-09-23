package executionvenue

import (
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/executionvenue"
"github.com/ettec/otp-common/model"
	"github.com/golang/protobuf/proto"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

var orderCache *executionvenue.OrderCache
var om orderManager

func setup() {
	var err error
	orderCache, err = executionvenue.NewOrderCache(newTestOrderStore(), "")
	if err != nil {
		panic(err)
	}

	om = NewOrderManager(orderCache, &TestOrderManager{}, func(listingId int32, result chan<- *model.Listing) {
		result <- &model.Listing{Id:1}
	})
}

func teardown() {
	defer om.Close()
}

func IntToDecimal64(i int) *model.Decimal64 {
	return &model.Decimal64{
		Mantissa: int64(i),
		Exponent: 0,
	}
}

func TestOrderManagerImpl_CreateAndRouteOrderPartialFillImmediate(t *testing.T) {

	params := &api.CreateAndRouteOrderParams{
		OrderSide: model.Side_BUY,
		Quantity:  IntToDecimal64(10),
		Price:     IntToDecimal64(20),
		ListingId:   1,
	}

	id, err := om.CreateAndRouteOrder(params)
	if err != nil {
		t.Fatalf("Create order call failed %v", err)
	}

	order, ok,_ := orderCache.GetOrder(id.OrderId)

	if !ok {
		t.Fatalf("created order not found in store")
	}

	om.AddExecution(id.OrderId, *IntToDecimal64(20), *IntToDecimal64(5), "testexecid" )


	testOrder := &model.Order{
		Version:           1,
		Id:                id.OrderId,
		Side:              model.Side_BUY,
		Quantity:          IntToDecimal64(10),
		Price:             IntToDecimal64(20),
		ListingId:         1,
		RemainingQuantity: IntToDecimal64(0),
		Status:            model.OrderStatus_FILLED,
		TargetStatus:      model.OrderStatus_NONE,
	}

	order, ok,_ = orderCache.GetOrder(id.OrderId)

	order.Created = nil
	testOrder.Created = nil

	if order.Status != model.OrderStatus_LIVE {
		t.Fatalf("Expected order to be live")
	}

}


func TestOrderManagerImpl_CreateAndRouteOrderFullyFilledImmediate(t *testing.T) {

	params := &api.CreateAndRouteOrderParams{
		OrderSide: model.Side_BUY,
		Quantity:  IntToDecimal64(10),
		Price:     IntToDecimal64(20),
		ListingId:   1,
	}

	id, err := om.CreateAndRouteOrder(params)
	if err != nil {
		t.Fatalf("Create order call failed %v", err)
	}

	order, ok, err := orderCache.GetOrder(id.OrderId)

	if !ok {
		t.Fatalf("created order not found in store")
	}

	om.AddExecution(id.OrderId, *IntToDecimal64(20), *IntToDecimal64(10), "testexecid" )


	testOrder := &model.Order{
		Id:                id.OrderId,
		Side:              model.Side_BUY,
		Quantity:          IntToDecimal64(10),
		Price:             IntToDecimal64(20),
		ListingId:         1,
		RemainingQuantity: IntToDecimal64(0),
		Status:            model.OrderStatus_FILLED,
		TargetStatus:      model.OrderStatus_NONE,
	}

	order, ok,_ = orderCache.GetOrder(id.OrderId)

	order.Created = nil
	testOrder.Created = nil

	if order.Status != model.OrderStatus_FILLED {
		t.Fatalf("Expected order to be filled")
	}

}

func TestOrderManagerImpl_CreateAndRouteOrder(t *testing.T) {

	params := &api.CreateAndRouteOrderParams{
		OrderSide: model.Side_BUY,
		Quantity:  IntToDecimal64(10),
		Price:     IntToDecimal64(20),
		ListingId:   1,
	}

	id, err := om.CreateAndRouteOrder(params)
	if err != nil {
		t.Fatalf("Create order call failed %v", err)
	}

	order, ok,_ := orderCache.GetOrder(id.OrderId)

	if !ok {
		t.Fatalf("created order not found in store")
	}

	om.SetOrderStatus(id.OrderId, model.OrderStatus_LIVE)

	testOrder := &model.Order{
        Version: 			1,
		Id:                id.OrderId,
		Side:              model.Side_BUY,
		Quantity:          IntToDecimal64(10),
		Price:             IntToDecimal64(20),
		ListingId:         1,
		RemainingQuantity: IntToDecimal64(10),
		Status:            model.OrderStatus_LIVE,
		TargetStatus:      model.OrderStatus_NONE,
	}

	order, ok,_ = orderCache.GetOrder(id.OrderId)

	order.Created = nil
	testOrder.Created = nil

	if !proto.Equal(testOrder, order) {
		t.Fatalf("Expected order %v, got %v", testOrder, order)
	}

}

func TestOrderManagerImpl_CancelOrder(t *testing.T) {

	listing := &model.Listing{Id: 1}

	params := &api.CreateAndRouteOrderParams{
		OrderSide: model.Side_BUY,
		Quantity:  IntToDecimal64(10),
		Price:     IntToDecimal64(20),
		ListingId:   1,
		Destination: "XNAS",
	}

	id, _ := om.CreateAndRouteOrder(params)

	om.SetOrderStatus(id.OrderId, model.OrderStatus_LIVE)

	err := om.CancelOrder(&api.CancelOrderParams{
		OrderId: id.OrderId,
		ListingId: listing.Id,
		OwnerId: "XNAS",
	})
	if err != nil {
		t.Fatalf("cancel order call failed: %v", err)
	}

	om.SetOrderStatus(id.OrderId, model.OrderStatus_CANCELLED)

	testOrder := &model.Order{
        Version:           3,
		Id:                id.OrderId,
		Side:              model.Side_BUY,
		Quantity:          IntToDecimal64(10),
		Price:             IntToDecimal64(20),
		ListingId:         1,
		RemainingQuantity: IntToDecimal64(10),
		Status:            model.OrderStatus_CANCELLED,
		TargetStatus:      model.OrderStatus_NONE,
		Destination: "XNAS",
	}

	order, _, _ := orderCache.GetOrder(id.OrderId)

	order.Created = nil
	testOrder.Created = nil

	if !proto.Equal(testOrder, order) {
		t.Fatalf("Expected order %v, got %v", testOrder, order)
	}

}

type TestOrderManager struct {
}

func (f *TestOrderManager) Send(order *model.Order, listing *model.Listing) error {
	return nil
}

func (f *TestOrderManager) Cancel(order *model.Order) error {
	return nil
}

func (f *TestOrderManager) Modify(order *model.Order, listing *model.Listing, Quantity *model.Decimal64, Price *model.Decimal64) error {
	return nil
}



func newTestOrderStore() *testOrderStore {
	t := testOrderStore{
		orders:    make([]*model.Order, 0, 10),
		ordersMap: make(map[string]*model.Order),
	}

	return &t
}

type testOrderStore struct {
	orders    []*model.Order
	ordersMap map[string]*model.Order
}

func (t *testOrderStore) Write(order *model.Order) error {


	t.orders = append(t.orders, order)


	t.ordersMap[order.Id] = order



	return nil
}

func (t *testOrderStore) RecoverInitialCache(loadOrder func(order *model.Order) bool) (map[string]*model.Order, error) {
	return map[string]*model.Order{}, nil
}

func (t *testOrderStore) GetOrder(orderId string) (model.Order, bool) {
	order, ok := t.ordersMap[orderId]
	return *order, ok
}

func (t *testOrderStore) Close() {

}
