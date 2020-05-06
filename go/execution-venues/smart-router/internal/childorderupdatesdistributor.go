package internal

import (
	"github.com/ettec/open-trading-platform/go/model"
	"log"
)

type childOrderStream struct {
	parentOrderId string
	orderChan     chan *model.Order
	distributor   *childOrderUpdatesDistributor
}

func newChildOrderStream(parentOrderId string, bufferSize int, d *childOrderUpdatesDistributor) *childOrderStream {
	stream := &childOrderStream{parentOrderId: parentOrderId, orderChan: make(chan *model.Order, bufferSize), distributor: d}
	d.openOrderChan <- parentIdAndChan{
		parentId:  parentOrderId,
		orderChan: stream.orderChan,
	}
	return stream
}

func (c *childOrderStream) GetStream() <-chan *model.Order {
	return c.orderChan
}

func (c *childOrderStream) Close() {
	c.distributor.closeOrderChan <- c.parentOrderId
}

type parentIdAndChan struct {
	parentId  string
	orderChan chan *model.Order
}

type childOrderUpdatesDistributor struct {
	openOrderChan  chan parentIdAndChan
	closeOrderChan chan string
	startChan      chan bool
}

type ChildOrderStream interface {
	GetStream() <-chan *model.Order
	Close()
}

func (d *childOrderUpdatesDistributor) NewOrderStream(parentOrderId string, bufferSize int) ChildOrderStream {
	return newChildOrderStream(parentOrderId, bufferSize, d)
}

func NewChildOrderUpdatesDistributor(updates <-chan ChildOrder) *childOrderUpdatesDistributor {

	idToChan := map[string]chan *model.Order{}

	d := &childOrderUpdatesDistributor{
		openOrderChan:  make(chan parentIdAndChan),
		closeOrderChan: make(chan string),
		startChan:      make(chan bool),
	}

	go func() {

	subscribeLoop:
		for {
			select {
			case <-d.startChan:
				break subscribeLoop
			case o := <-d.openOrderChan:
				idToChan[o.parentId] = o.orderChan
			case c := <-d.closeOrderChan:
				delete(idToChan, c)

			}
		}

		for {
			select {
			case u := <-updates:
				if orderChan, ok := idToChan[u.ParentOrderId]; ok {
					select {
					case orderChan <- u.Child:
					default:
						log.Printf("slow consumer, closing child order update channel, parent order id %v", u.ParentOrderId)
						close(orderChan)
						delete(idToChan, u.ParentOrderId)
					}

				}
			case o := <-d.openOrderChan:
				idToChan[o.parentId] = o.orderChan
			case pId := <-d.closeOrderChan:
				updateChan := idToChan[pId]
				delete(idToChan, pId)
				close(updateChan)
			}
		}

	}()

	return d
}

func (d *childOrderUpdatesDistributor) Start() {
	d.startChan <- true
}
