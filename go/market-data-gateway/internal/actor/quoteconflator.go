package actor

import "github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"

type quoteConflator struct {
	inChan       <-chan *model.ClobQuote
	outChan      chan<- *model.ClobQuote
	closeChan    <-chan bool
}

func NewQuoteConflator(inChan <-chan *model.ClobQuote, outChan chan<- *model.ClobQuote, closeChan <-chan bool) *quoteConflator {
	c := &quoteConflator{
		inChan: inChan, outChan: outChan, closeChan: closeChan}

	pendingQuote := map[int32]*model.ClobQuote{}


	go func() {


		for {
			var eq *model.ClobQuote

			for _, v := range pendingQuote {
				eq = v
			}

			if eq != nil {
				select {
				case q := <-c.inChan:
					if _, ok := pendingQuote[q.ListingId]; !ok {
						//order[writePtr] = q.ListingId
// here
					//	writePtr++
					}
					pendingQuote[q.ListingId] = q

				case c.outChan <- eq:
					delete(pendingQuote, eq.ListingId)
				case <-c.closeChan:
					return
				}

			} else {
				select {
				case q := <-c.inChan:
					pendingQuote[q.ListingId] = q
				case <-c.closeChan:
					return
				}

			}
		}

	}()

	return c
}

type boundedCircularInt32Buffer struct {
	buffer   []int32
	capacity int
	len      int
	readPtr  int
	writePtr int
}

func newBoundedCircularIntBuffer(capacity int) *boundedCircularInt32Buffer {
	b:=  &boundedCircularInt32Buffer{buffer: make([]int32, capacity, capacity), capacity:capacity}


	return b
}

// true if the buffer is not full and the value is added
func (b *boundedCircularInt32Buffer) addHead(i int32) bool {

	if b.len == b.capacity {
		return false
	}


	b.buffer[b.writePtr] = i
	b.len++

	if b.writePtr == b.capacity -1{
		b.writePtr = 0
	} else {
		b.writePtr++
	}


	return true

}

// returns the value and true if a value is available
func (b *boundedCircularInt32Buffer) removeTail() (int32, bool){
	if b.len == 0 {
		return 0, false
	}


	res := b.buffer[b.readPtr]
	b.len--
	b.readPtr++
	if b.readPtr == b.capacity {
		b.readPtr =0
	}

	return res, true

}
