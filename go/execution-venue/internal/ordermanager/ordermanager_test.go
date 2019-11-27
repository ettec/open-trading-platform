package ordermanager

import (
	"github.com/ettec/open-trading-platform/execution-venue/internal/ordercache"
	"github.com/ettec/open-trading-platform/execution-venue/pb"
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

	avgPrice := pb.Decimal64{
		Mantissa:             0,
		Exponent:             0,
	}

	trdQnt := pb.Decimal64{
		Mantissa:             0,
		Exponent:             0,
	}

	lastPrice := pb.Decimal64{
		Mantissa:             10,
		Exponent:             1,
	}

	lastQnt := pb.Decimal64{
		Mantissa:             5,
		Exponent:             0,
	}


	resultAsDecimal64 := calculateAveragePrice(&avgPrice, &trdQnt, &lastPrice, &lastQnt)
	result, _ := resultAsDecimal64.AsDecimal().Float64()

	expectedResult := 100.0
	if !floatEquals(result, expectedResult) {
		t.Fatalf("Expected avg price %v, got %v", expectedResult, result)
	}


	trdQnt = pb.Decimal64{
		Mantissa:             5,
		Exponent:             0,
	}

	lastPrice = pb.Decimal64{
		Mantissa:             80,
		Exponent:             0,
	}

	lastQnt = pb.Decimal64{
		Mantissa:             5,
		Exponent:             0,
	}

	result, _ = calculateAveragePrice(resultAsDecimal64, &trdQnt, &lastPrice, &lastQnt).AsDecimal().Float64()

	expectedResult = 90.0
	if !floatEquals(result, expectedResult) {
		t.Fatalf("Expected avg price %v, got %v", expectedResult, result)
	}


}

var EPSILON float64 = 0.00000001
func floatEquals(a, b float64) bool {
	if (a - b) < EPSILON && (b - a) < EPSILON {
		return true
	}
	return false
}

func TestOrderManagerImpl_UpdateTradedQuantity(t *testing.T) {
	params := &pb.CreateAndRouteOrderParams{
		Side:      pb.Side_BUY,
		Quantity:  pb.IntToDecimal64(15),
		Price:     pb.IntToDecimal64(20),
		ListingId: "listing1",
	}

	id, err := orderManager.CreateAndRouteOrder(params)
	if err != nil {
		t.Fatalf("Create order call failed %v", err)
	}

	order, ok := orderCache.GetOrder(id.OrderId)

	if !ok {
		t.Fatalf("created order not found in store")
	}

	orderManager.SetOrderStatus(id.OrderId, pb.OrderStatus_LIVE)

	err = orderManager.UpdateTradedQuantity(id.OrderId, pb.Decimal64{
		Mantissa:             170,
		Exponent:             -1,
	}, pb.Decimal64{
		Mantissa:             5,
		Exponent:             0,

	})

	if err != nil {
		t.Fatalf("error %v", err)
	}

	err = orderManager.UpdateTradedQuantity(id.OrderId, pb.Decimal64{
		Mantissa:             19,
		Exponent:             0,
	}, pb.Decimal64{
		Mantissa:             5,
		Exponent:             0,

	})

	if err != nil {
		t.Fatalf("error %v", err)
	}

	testOrder := &pb.Order{
		Version:              3,
		Id:                   id.OrderId,
		Side:                 pb.Side_BUY,
		Quantity:             pb.IntToDecimal64(15),
		Price:                pb.IntToDecimal64(20),
		ListingId:            "listing1",
		RemainingQuantity:    pb.IntToDecimal64(5),
		Status:               pb.OrderStatus_LIVE,
		TargetStatus:         pb.OrderStatus_NONE,
		TradedQuantity:       pb.IntToDecimal64(10),
		AvgTradePrice:        &pb.Decimal64{
			Mantissa:             180000000000000000,
			Exponent:             -16,
		},

	}


	if !proto.Equal(testOrder, order) {
		t.Fatalf("Expected order %v, got %v", testOrder, order)
	}


}

func TestOrderManagerImpl_UpdateTradedQuantityOnPendingLiveOrder(t *testing.T) {
	params := &pb.CreateAndRouteOrderParams{
		Side:      pb.Side_BUY,
		Quantity:  pb.IntToDecimal64(15),
		Price:     pb.IntToDecimal64(20),
		ListingId: "listing1",
	}

	id, err := orderManager.CreateAndRouteOrder(params)
	if err != nil {
		t.Fatalf("Create order call failed %v", err)
	}


	err = orderManager.UpdateTradedQuantity(id.OrderId, pb.Decimal64{
		Mantissa:             170,
		Exponent:             -1,
	}, pb.Decimal64{
		Mantissa:             15,
		Exponent:             0,

	})

	if err != nil {
		t.Fatalf("error %v", err)
	}

	order, ok := orderCache.GetOrder(id.OrderId)

	if !ok {
		t.Fatalf("created order not found in store")
	}

	testOrder := &pb.Order{
		Version:              1,
		Id:                   id.OrderId,
		Side:                 pb.Side_BUY,
		Quantity:             pb.IntToDecimal64(15),
		Price:                pb.IntToDecimal64(20),
		ListingId:            "listing1",
		RemainingQuantity:    pb.IntToDecimal64(0),
		Status:               pb.OrderStatus_FILLED,
		TargetStatus:         pb.OrderStatus_NONE,
		TradedQuantity:		  pb.IntToDecimal64(15),
		AvgTradePrice:        &pb.Decimal64{
			Mantissa:             170000000000000000,
			Exponent:             -16,
		},

	}


	if !proto.Equal(testOrder, order) {
		t.Fatalf("Expected order %v,\n got %v", testOrder, order)
	}


}

func TestOrderManagerImpl_CreateAndRouteOrder(t *testing.T) {

	params := &pb.CreateAndRouteOrderParams{
		Side:      pb.Side_BUY,
		Quantity:  pb.IntToDecimal64(10),
		Price:     pb.IntToDecimal64(20),
		ListingId: "listing1",
	}

	id, err := orderManager.CreateAndRouteOrder(params)
	if err != nil {
		t.Fatalf("Create order call failed %v", err)
	}

	order, ok := orderCache.GetOrder(id.OrderId)

	if !ok {
		t.Fatalf("created order not found in store")
	}

	orderManager.SetOrderStatus(id.OrderId, pb.OrderStatus_LIVE)

	testOrder := &pb.Order{
		Version:              1,
		Id:                   id.OrderId,
		Side:                 pb.Side_BUY,
		Quantity:             pb.IntToDecimal64(10),
		Price:                pb.IntToDecimal64(20),
		ListingId:            "listing1",
		RemainingQuantity:    pb.IntToDecimal64(10),
		Status:               pb.OrderStatus_LIVE,
		TargetStatus:         pb.OrderStatus_NONE,

	}


	if !proto.Equal(testOrder, order) {
		t.Fatalf("Expected order %v, got %v", testOrder, order)
	}

}

func TestOrderManagerImpl_CancelOrder(t *testing.T) {



	params := &pb.CreateAndRouteOrderParams{
		Side:      pb.Side_BUY,
		Quantity:  pb.IntToDecimal64(10),
		Price:     pb.IntToDecimal64(20),
		ListingId: "listing1",
	}

	id, _ := orderManager.CreateAndRouteOrder(params)

	orderManager.SetOrderStatus(id.OrderId, pb.OrderStatus_LIVE)


	err := orderManager.CancelOrder(id)
	if err != nil {
		t.Fatalf("cancel order call failed: %v", err)
	}

	orderManager.SetOrderStatus(id.OrderId, pb.OrderStatus_CANCELLED)

	testOrder := &pb.Order{
		Version:              3,
		Id:                   id.OrderId,
		Side:                 pb.Side_BUY,
		Quantity:             pb.IntToDecimal64(10),
		Price:                pb.IntToDecimal64(20),
		ListingId:            "listing1",
		RemainingQuantity:    pb.IntToDecimal64(10),
		Status:               pb.OrderStatus_CANCELLED,
		TargetStatus:         pb.OrderStatus_NONE,

	}

	order, _ := orderCache.GetOrder(id.OrderId)

	if !proto.Equal(testOrder, order) {
		t.Fatalf("Expected order %v, got %v", testOrder, order)
	}


}



type TestOrderManager struct {

}

func (f *TestOrderManager) Send(order *pb.Order) error {
	return nil
}


func NewTestOrderStore() *TestOrderStore {
	t := TestOrderStore{
		orders:    make([]*pb.Order,0,10),
		ordersMap:  make(map[string]*pb.Order),
	}

	return &t
}

type TestOrderStore struct {
	orders    []*pb.Order
	ordersMap map[string]*pb.Order
}

func (t *TestOrderStore) Write(order *pb.Order) error {
	t.orders = append(t.orders, order)
	t.ordersMap[order.Id] = order
	return nil
}

func (t *TestOrderStore) GetOrder(orderId string) (pb.Order, bool) {
	order, ok := t.ordersMap[orderId]
	return *order, ok
}

func (t *TestOrderStore) Close() {

}