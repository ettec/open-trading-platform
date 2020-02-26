package ordermanager

import (
	"github.com/ettec/open-trading-platform/go/execution-venue/api"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/ordercache"
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
	orderCache = ordercache.NewOrderCache(NewTestOrderStore())
	orderManager = NewOrderManager(orderCache, &TestOrderManager{})
}

func teardown() {
	defer orderManager.Close()
}

func TestOrderManagerImpl_CalculateAveragePrice(t *testing.T) {

	avgPrice := model.Decimal64{
		Mantissa: 0,
		Exponent: 0,
	}

	trdQnt := model.Decimal64{
		Mantissa: 0,
		Exponent: 0,
	}

	lastPrice := model.Decimal64{
		Mantissa: 10,
		Exponent: 1,
	}

	lastQnt := model.Decimal64{
		Mantissa: 5,
		Exponent: 0,
	}

	resultAsDecimal64 := calculateAveragePrice(&avgPrice, &trdQnt, &lastPrice, &lastQnt)
	result, _ := resultAsDecimal64.AsDecimal().Float64()

	expectedResult := 100.0
	if !floatEquals(result, expectedResult) {
		t.Fatalf("Expected avg price %v, got %v", expectedResult, result)
	}

	trdQnt = model.Decimal64{
		Mantissa: 5,
		Exponent: 0,
	}

	lastPrice = model.Decimal64{
		Mantissa: 80,
		Exponent: 0,
	}

	lastQnt = model.Decimal64{
		Mantissa: 5,
		Exponent: 0,
	}

	result, _ = calculateAveragePrice(resultAsDecimal64, &trdQnt, &lastPrice, &lastQnt).AsDecimal().Float64()

	expectedResult = 90.0
	if !floatEquals(result, expectedResult) {
		t.Fatalf("Expected avg price %v, got %v", expectedResult, result)
	}

}

var EPSILON = 0.00000001

func floatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func TestOrderManagerImpl_UpdateTradedQuantity(t *testing.T) {
	params := &api.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity: model.IntToDecimal64(15),
		Price:    model.IntToDecimal64(20),
		Listing:  &model.Listing{Id: 1},
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

	err = orderManager.UpdateTradedQuantity(id.OrderId, model.Decimal64{
		Mantissa: 170,
		Exponent: -1,
	}, model.Decimal64{
		Mantissa: 5,
		Exponent: 0,
	})

	if err != nil {
		t.Fatalf("error %v", err)
	}

	err = orderManager.UpdateTradedQuantity(id.OrderId, model.Decimal64{
		Mantissa: 19,
		Exponent: 0,
	}, model.Decimal64{
		Mantissa: 5,
		Exponent: 0,
	})

	if err != nil {
		t.Fatalf("error %v", err)
	}

	testOrder := &model.Order{
		Version:           3,
		Id:                id.OrderId,
		Side:              model.Side_BUY,
		Quantity:          model.IntToDecimal64(15),
		Price:             model.IntToDecimal64(20),
		ListingId:         1,
		RemainingQuantity: model.IntToDecimal64(5),
		Status:            model.OrderStatus_LIVE,
		TargetStatus:      model.OrderStatus_NONE,
		TradedQuantity:    model.IntToDecimal64(10),
		AvgTradePrice: &model.Decimal64{
			Mantissa: 180000000000000000,
			Exponent: -16,
		},
	}

	order.Created = nil
	testOrder.Created = nil

	if !proto.Equal(testOrder, order) {
		t.Fatalf("Expected order %v, got %v", testOrder, order)
	}

}

func TestOrderManagerImpl_UpdateTradedQuantityOnPendingLiveOrder(t *testing.T) {
	params := &api.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity: model.IntToDecimal64(15),
		Price:    model.IntToDecimal64(20),
		Listing:  &model.Listing{Id: 1},
	}

	id, err := orderManager.CreateAndRouteOrder(params)
	if err != nil {
		t.Fatalf("Create order call failed %v", err)
	}

	err = orderManager.UpdateTradedQuantity(id.OrderId, model.Decimal64{
		Mantissa: 170,
		Exponent: -1,
	}, model.Decimal64{
		Mantissa: 15,
		Exponent: 0,
	})

	if err != nil {
		t.Fatalf("error %v", err)
	}

	order, ok := orderCache.GetOrder(id.OrderId)

	if !ok {
		t.Fatalf("created order not found in store")
	}

	testOrder := &model.Order{
		Version:           1,
		Id:                id.OrderId,
		Side:              model.Side_BUY,
		Quantity:          model.IntToDecimal64(15),
		Price:             model.IntToDecimal64(20),
		ListingId:         1,
		RemainingQuantity: model.IntToDecimal64(0),
		Status:            model.OrderStatus_FILLED,
		TargetStatus:      model.OrderStatus_NONE,
		TradedQuantity:    model.IntToDecimal64(15),
		AvgTradePrice: &model.Decimal64{
			Mantissa: 170000000000000000,
			Exponent: -16,
		},
	}

	order.Created = nil
	testOrder.Created = nil

	if !proto.Equal(testOrder, order) {
		t.Fatalf("Expected order %v,\n got %v", testOrder, order)
	}

}

func TestOrderManagerImpl_CreateAndRouteOrder(t *testing.T) {

	params := &api.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity: model.IntToDecimal64(10),
		Price:    model.IntToDecimal64(20),
		Listing:  &model.Listing{Id: 1},
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
		Quantity:          model.IntToDecimal64(10),
		Price:             model.IntToDecimal64(20),
		ListingId:         1,
		RemainingQuantity: model.IntToDecimal64(10),
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

	params := &api.CreateAndRouteOrderParams{
		OrderSide:     model.Side_BUY,
		Quantity: model.IntToDecimal64(10),
		Price:    model.IntToDecimal64(20),
		Listing:  &model.Listing{Id: 1},
	}

	id, _ := orderManager.CreateAndRouteOrder(params)

	orderManager.SetOrderStatus(id.OrderId, model.OrderStatus_LIVE)

	err := orderManager.CancelOrder(id)
	if err != nil {
		t.Fatalf("cancel order call failed: %v", err)
	}

	orderManager.SetOrderStatus(id.OrderId, model.OrderStatus_CANCELLED)

	testOrder := &model.Order{
		Version:           3,
		Id:                id.OrderId,
		Side:              model.Side_BUY,
		Quantity:          model.IntToDecimal64(10),
		Price:             model.IntToDecimal64(20),
		ListingId:         1,
		RemainingQuantity: model.IntToDecimal64(10),
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

func (f *TestOrderManager)  Cancel(order *model.Order) error {
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

func (t *TestOrderStore) GetOrder(orderId string) (model.Order, bool) {
	order, ok := t.ordersMap[orderId]
	return *order, ok
}

func (t *TestOrderStore) Close() {

}
