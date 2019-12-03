package messagesource

import (
	"context"
	"github.com/ettec/open-trading-platform/view-service/internal/model"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"time"
)

type SimulatedMessageSource struct {
	orderChan chan *model.Order
	done      chan bool
}

func NewSimulatedMessageSource() *SimulatedMessageSource {
	ks := &SimulatedMessageSource{}
	ks.orderChan = make(chan *model.Order, 10)
	ks.done = make(chan bool, 1)
	go ks.start()

	return ks
}

func (k *SimulatedMessageSource) start() {
	ticker := time.NewTicker(2000 * time.Millisecond)
	var version int32  = 0
	defer ticker.Stop()
	for {
		select {
		case <-k.done:
			return
		case  <-ticker.C:
			version++
			order := &model.Order{
				Version:              version,
				Id:                   "testorderid",
				Side:                 0,
				Quantity:             &model.Decimal64{
					Mantissa:              int64(rand.Int()),
					Exponent:             0,
				},
				Price:				   &model.Decimal64{
					Mantissa:             int64(rand.Int()),
					Exponent:             0,

				},
				ListingId:            "",
				RemainingQuantity:    nil,
				TradedQuantity:       nil,
				AvgTradePrice:        nil,
				Status:               model.OrderStatus_LIVE,
				TargetStatus:         model.OrderStatus_NONE,

			}

			k.orderChan <- order

		}

	}

}

func (k *SimulatedMessageSource) ReadMessage(ctx context.Context) ([]byte, error) {

	order := <-k.orderChan
	bytes, err := proto.Marshal(order)
	return bytes, err
}

func (k *SimulatedMessageSource) Close() error {
	k.done <- true
	return nil
}
