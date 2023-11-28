package executionvenue

import (
	"context"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/ordermanagement"
	"github.com/ettec/otp-common/staticdata"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setup(ctx)
	code := m.Run()
	os.Exit(code)
}

var orderCache *ordermanagement.OrderCache
var om orderManager

func setup(ctx context.Context) {
	var err error
	orderCache, err = ordermanagement.NewOwnerOrderCache(ctx, "", newTestOrderStore())
	if err != nil {
		panic(err)
	}

	om = NewOrderManager(ctx, orderCache, &TestOrderManager{}, func(ctx context.Context, listingId int32, result chan<- staticdata.ListingResult) {
		result <- staticdata.ListingResult{Listing: &model.Listing{Id: 1}}
	}, 100)
}

func IntToDecimal64(i int) *model.Decimal64 {
	return &model.Decimal64{
		Mantissa: int64(i),
		Exponent: 0,
	}
}

func TestCreateAndRouteOrderPartialFillImmediate(t *testing.T) {

	params := &api.CreateAndRouteOrderParams{
		OrderSide: model.Side_BUY,
		Quantity:  IntToDecimal64(10),
		Price:     IntToDecimal64(20),
		ListingId: 1,
	}

	id, err := om.CreateAndRouteOrder(params)
	if err != nil {
		t.Fatalf("Create order call failed %v", err)
	}

	order, ok, _ := orderCache.GetOrder(id.OrderId)

	if !ok {
		t.Fatalf("created order not found in store")
	}

	err = om.AddExecution(id.OrderId, *IntToDecimal64(20), *IntToDecimal64(5), "testexecid")
	assert.NoError(t, err)

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

	order, ok, _ = orderCache.GetOrder(id.OrderId)

	order.Created = nil
	testOrder.Created = nil

	if order.Status != model.OrderStatus_LIVE {
		t.Fatalf("Expected order to be live")
	}

}

func TestCreateAndRouteOrderFullyFilledImmediate(t *testing.T) {

	params := &api.CreateAndRouteOrderParams{
		OrderSide: model.Side_BUY,
		Quantity:  IntToDecimal64(10),
		Price:     IntToDecimal64(20),
		ListingId: 1,
	}

	id, err := om.CreateAndRouteOrder(params)
	if err != nil {
		t.Fatalf("Create order call failed %v", err)
	}

	order, ok, err := orderCache.GetOrder(id.OrderId)

	if !ok {
		t.Fatalf("created order not found in store")
	}

	err = om.AddExecution(id.OrderId, *IntToDecimal64(20), *IntToDecimal64(10), "testexecid")
	assert.NoError(t, err)

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

	order, ok, _ = orderCache.GetOrder(id.OrderId)

	order.Created = nil
	testOrder.Created = nil

	if order.Status != model.OrderStatus_FILLED {
		t.Fatalf("Expected order to be filled")
	}

}

func TestCreateAndRouteOrder(t *testing.T) {

	params := &api.CreateAndRouteOrderParams{
		OrderSide: model.Side_BUY,
		Quantity:  IntToDecimal64(10),
		Price:     IntToDecimal64(20),
		ListingId: 1,
	}

	id, err := om.CreateAndRouteOrder(params)
	if err != nil {
		t.Fatalf("Create order call failed %v", err)
	}

	order, ok, _ := orderCache.GetOrder(id.OrderId)

	if !ok {
		t.Fatalf("created order not found in store")
	}

	err = om.SetOrderStatus(id.OrderId, model.OrderStatus_LIVE)
	assert.NoError(t, err)

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

	order, ok, _ = orderCache.GetOrder(id.OrderId)

	order.Created = nil
	testOrder.Created = nil

	if !proto.Equal(testOrder, order) {
		t.Fatalf("Expected order %v, got %v", testOrder, order)
	}

}

func TestCancelOrder(t *testing.T) {

	listing := &model.Listing{Id: 1}

	params := &api.CreateAndRouteOrderParams{
		OrderSide:   model.Side_BUY,
		Quantity:    IntToDecimal64(10),
		Price:       IntToDecimal64(20),
		ListingId:   1,
		Destination: "XNAS",
	}

	id, _ := om.CreateAndRouteOrder(params)

	err := om.SetOrderStatus(id.OrderId, model.OrderStatus_LIVE)
	assert.NoError(t, err)

	err = om.CancelOrder(&api.CancelOrderParams{
		OrderId:   id.OrderId,
		ListingId: listing.Id,
		OwnerId:   "XNAS",
	})
	assert.NoError(t, err)

	err = om.SetOrderStatus(id.OrderId, model.OrderStatus_CANCELLED)
	assert.NoError(t, err)

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
		Destination:       "XNAS",
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

func (f *TestOrderManager) Send(_ *model.Order, _ *model.Listing) error {
	return nil
}

func (f *TestOrderManager) Cancel(_ *model.Order) error {
	return nil
}

func (f *TestOrderManager) Modify(_ *model.Order, _ *model.Listing, _ *model.Decimal64, _ *model.Decimal64) error {
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

func (t *testOrderStore) Write(_ context.Context, order *model.Order) error {

	t.orders = append(t.orders, order)

	t.ordersMap[order.Id] = order

	return nil
}

func (t *testOrderStore) LoadOrders(_ context.Context, _ func(order *model.Order) bool) (map[string]*model.Order, error) {
	return map[string]*model.Order{}, nil
}

func (t *testOrderStore) GetOrder(orderId string) (model.Order, bool) {
	order, ok := t.ordersMap[orderId]
	return *order, ok
}

func (t *testOrderStore) Close() {

}
