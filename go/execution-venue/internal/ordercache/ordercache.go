package ordercache

import (
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/ordercache/orderstore"
)

type OrderCache struct {
	store orderstore.OrderStore
	cache map[string]*model.Order
}

func NewOrderCache(store orderstore.OrderStore) *OrderCache {
	orderCache := OrderCache{
		store: store,
		cache: make(map[string]*model.Order, 100),
	}

	return &orderCache
}


func (oc *OrderCache) Store(order *model.Order) error {

	existingOrder, exists := oc.cache[order.Id]
	if exists {
		order.Version = existingOrder.Version + 1
	} else {
		order.Version = 0
	}

	e := oc.store.Write(order)
	if e != nil {
		return e
	}

	oc.cache[order.Id] = order

	return nil
}

// Returns the order and true if found, otherwise a nil value and false
func (oc *OrderCache) GetOrder(orderId string) (*model.Order, bool) {
	order, ok := oc.cache[orderId]

	if ok {
		return order, true
	} else {
		return nil, false
	}
}

func (oc *OrderCache) Close() {
	oc.store.Close()
}
