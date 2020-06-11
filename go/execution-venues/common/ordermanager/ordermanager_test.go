package ordermanager

import (
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/execution-venues/common/ordercache"
	"github.com/ettec/open-trading-platform/go/model"
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

var orderCache *ordercache.OrderCache
var orderManager OrderManager

func setup() {
	var err error
	orderCache, err = ordercache.NewOrderCache(NewTestOrderStore())
	if err != nil {
		panic(err)
	}

	orderManager = NewOrderManager(orderCache, &TestOrderManager{}, "")
}

func teardown() {
	defer orderManager.Close()
}

func IntToDecimal64(i int) *model.Decimal64 {
	return &model.Decimal64{
		Mantissa: int64(i),
		Exponent: 0,
	}
}

func TestOrderManagerImpl_CreateAndRouteOrder(t *testing.T) {

	params := &api.CreateAndRouteOrderParams{
		OrderSide: model.Side_BUY,
		Quantity:  IntToDecimal64(10),
		Price:     IntToDecimal64(20),
		Listing:   &model.Listing{Id: 1},
	}

	id, err := orderManager.CreateAndRouteOrder(params)
	if err != nil {
		t.Fatalf("Create order call failed %v", err)
	}

	order, ok := orderCache.GetOrder(id.OrderId)

	if !ok {
		t.Fatalf("created order not found in store")
	}

	orderManager.SetOrderStatus(id.OrderId, model.OrderStatus_LIVE)

	testOrder := &model.Order{
		Version:           1,
		Id:                id.OrderId,
		Side:              model.Side_BUY,
		Quantity:          IntToDecimal64(10),
		Price:             IntToDecimal64(20),
		ListingId:         1,
		RemainingQuantity: IntToDecimal64(10),
		Status:            model.OrderStatus_LIVE,
		TargetStatus:      model.OrderStatus_NONE,
	}

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
		Listing:   listing,
	}

	id, _ := orderManager.CreateAndRouteOrder(params)

	orderManager.SetOrderStatus(id.OrderId, model.OrderStatus_LIVE)

	err := orderManager.CancelOrder(&api.CancelOrderParams{
		OrderId: id.OrderId,
		Listing: listing,
	})
	if err != nil {
		t.Fatalf("cancel order call failed: %v", err)
	}

	orderManager.SetOrderStatus(id.OrderId, model.OrderStatus_CANCELLED)

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
	}

	order, _ := orderCache.GetOrder(id.OrderId)

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

func (f *TestOrderManager)  Modify(order *model.Order, listing *model.Listing, Quantity *model.Decimal64, Price *model.Decimal64) error {
	return nil
}

func NewTestOrderStore() *TestOrderStore {
	t := TestOrderStore{
		orders:    make([]*model.Order, 0, 10),
		ordersMap: make(map[string]*model.Order),
	}

	return &t
}

type TestOrderStore struct {
	orders    []*model.Order
	ordersMap map[string]*model.Order
}

func (t *TestOrderStore) Write(order *model.Order) error {
	t.orders = append(t.orders, order)
	t.ordersMap[order.Id] = order
	return nil
}

func (t *TestOrderStore) RecoverInitialCache() (map[string]*model.Order, error) {
	return map[string]*model.Order{}, nil
}

func (t *TestOrderStore) GetOrder(orderId string) (model.Order, bool) {
	order, ok := t.ordersMap[orderId]
	return *order, ok
}

func (t *TestOrderStore) Close() {

}
