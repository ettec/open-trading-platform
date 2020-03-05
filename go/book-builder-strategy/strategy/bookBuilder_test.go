package strategy

import (
	"context"
	"encoding/json"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/depth"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/orderentryapi"
	"github.com/ettec/open-trading-platform/go/model"
	"google.golang.org/grpc"
	"testing"
)

const depthJson = `{"symbol":"XLF","marketPercent":0.02983,"volume":1569568,"lastSalePrice":27.03,"lastSaleSize":100,
"lastSaleTime":1583340199651,"lastUpdated":1583340210365,

"bids":[{"price":27.03,"size":1100,"timestamp":1583340200208},
{"price":27.02,"size":2500,"timestamp":1583340200865},{"price":27.01,"size":2500,"timestamp":1583340195958},
{"price":27,"size":2500,"timestamp":1583340195475},{"price":26.96,"size":100,"timestamp":1583340188516},
{"price":26.94,"size":100,"timestamp":1583339925086},{"price":26.9,"size":100,"timestamp":1583338223647},
{"price":26.89,"size":100,"timestamp":1583338217434}],

"asks":[{"price":27.06,"size":2500,"timestamp":1583340200308},
{"price":27.07,"size":2500,"timestamp":1583340194981},
{"price":27.08,"size":2500,"timestamp":1583340200206}],

"systemEvent":{"systemEvent":"R","timestamp":1583332200002},
"tradingStatus":{"status":"T","reason":"    ","timestamp":1583324123334},
"opHaltStatus":{"isHalted":false,"timestamp":1583324123334},
"ssrStatus":{"isSSR":false,"detail":" ","timestamp":1583324123334},
"securityEvent":{"securityEvent":"MarketOpen","timestamp":1583332200002},"trades":[],"tradeBreaks":[]}`

type testOrderEntryClient struct {
	submitOrderChan chan *orderentryapi.NewOrderParams


}

func newTestOrderEntryClient() *testOrderEntryClient{
	t := testOrderEntryClient {
		submitOrderChan: make(chan *orderentryapi.NewOrderParams),
	}

	return &t
}

func (o *testOrderEntryClient)	SubmitNewOrder(ctx context.Context, in *orderentryapi.NewOrderParams, opts ...grpc.CallOption) (*orderentryapi.OrderId, error){

}
func (o *testOrderEntryClient) CancelOrder(ctx context.Context, in *orderentryapi.OrderId, opts ...grpc.CallOption) (*orderentryapi.Empty, error) {

}


type testQuoteDist struct {
	sink chan<- *model.ClobQuote
	subscribeChan chan int32
}

func newTestQuoteDist() *testQuoteDist {
	t := &testQuoteDist{}
	t.subscribeChan = make(chan int32)
	return t
}

func (d *testQuoteDist) Subscribe(listingId int32, sink chan<- *model.ClobQuote) {
	d.subscribeChan <- listingId
}

func (d *testQuoteDist) AddOutQuoteChan(sink chan<- *model.ClobQuote) {
	d.sink = sink
}

func (d *testQuoteDist) RemoveOutQuoteChan(sink chan<- *model.ClobQuote) {
	d.sink = nil
}


func Test_bookBuilder_start(t *testing.T) {

	dep := depth.Depth{}
	json.Unmarshal([]byte(depthJson), &dep)

	oec := newTestOrderEntryClient()
	qd := newTestQuoteDist()

	book := newBookBuilder(&model.Listing{Id:1}, qd, dep, oec  )



	d := model.Decimal64{
		Mantissa:             1,
		Exponent:             2,
	}

	qd.sink <- &model.ClobQuote{
		ListingId:            1,
		Bids:                 []*model.ClobLine{{Size: &model.Decimal64{Mantissa:1,Exponent:2,},
												 Price:&model.Decimal64{Mantissa:3,Exponent:1,},},
												{Size: &model.Decimal64{Mantissa:2,Exponent:2,},
												 Price:&model.Decimal64{Mantissa:2,Exponent:1,},},

												 here finish off this testing

		},
		Offers:               nil,
	}


	p := <-oec.submitOrderChan

	//if p.Symbol != "XLF" || p.OrderSide != orderentryapi.Side_SELL || model.Decimal64(p.Price).Equal(modecl)





}