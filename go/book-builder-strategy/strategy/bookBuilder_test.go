package strategy

import (
	"context"
	"encoding/json"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/depth"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/orderentryapi"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"testing"
)

const depthJson = `{"symbol":"XLF","marketPercent":0.02983,"volume":1569568,"lastSalePrice":27.03,"lastSaleSize":100,
"lastSaleTime":1583340199651,"lastUpdated":1583340210365,

"bids":[{"price":27.03,"size":1100,"timestamp":1583340200208},
{"price":27.02,"size":2500,"timestamp":1583340200865},{"price":27.01,"size":2500,"timestamp":1583340195958}],

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

func newTestOrderEntryClient() *testOrderEntryClient {
	t := testOrderEntryClient{
		submitOrderChan: make(chan *orderentryapi.NewOrderParams),
	}

	return &t
}

func (o *testOrderEntryClient) SubmitNewOrder(ctx context.Context, in *orderentryapi.NewOrderParams, opts ...grpc.CallOption) (*orderentryapi.OrderId, error) {
	o.submitOrderChan <- in

	id, _ := uuid.NewUUID()

	return &orderentryapi.OrderId{OrderId: id.String(),}, nil
}
func (o *testOrderEntryClient) CancelOrder(ctx context.Context, in *orderentryapi.OrderId, opts ...grpc.CallOption) (*orderentryapi.Empty, error) {
	return nil, nil
}

type testQuoteDist struct {
	subscribeChan chan int32
	addQuoteChanChan chan chan<- *model.ClobQuote
}

func newTestQuoteDist() *testQuoteDist {
	t := &testQuoteDist{}
	t.subscribeChan = make(chan int32, 100)
	t.addQuoteChanChan = make(chan chan<- *model.ClobQuote, 100)
	return t
}

func (d *testQuoteDist) Subscribe(listingId int32, sink chan<- *model.ClobQuote) {
	d.subscribeChan <- listingId
}

func (d *testQuoteDist) AddOutQuoteChan(sink chan<- *model.ClobQuote) {
	d.addQuoteChanChan <- sink

}

func (d *testQuoteDist) RemoveOutQuoteChan(sink chan<- *model.ClobQuote) {

}

func Test_bookBuilder_start(t *testing.T) {

	dep := depth.Depth{}
	json.Unmarshal([]byte(depthJson), &dep)

	oec := newTestOrderEntryClient()
	qd := newTestQuoteDist()

	book := newBookBuilder(&model.Listing{Id: 1,MarketSymbol: "XLF",}, qd, dep, oec)

	book.start()

	sink := <- qd.addQuoteChanChan

	sink <- &model.ClobQuote{
		ListingId: 1,
		Bids: []*model.ClobLine{{Size: model.IasD(100),
			Price: model.IasD(20),},
			{Size: model.IasD(200),
				Price: model.IasD(30),},
		},
		Offers: []*model.ClobLine{{Size: model.IasD(400),
			Price: model.IasD(35),},
			{Size: model.IasD(200),
				Price: model.IasD(25),},
		}}

	p := <-oec.submitOrderChan
	verifyParams(p, "XLF", orderentryapi.Side_SELL, model.IasD(20), 300, t)

	p = <-oec.submitOrderChan
	verifyParams(p, "XLF", orderentryapi.Side_BUY, model.IasD(35), 600, t)



	p = <-oec.submitOrderChan
	verifyParams(p, "XLF", orderentryapi.Side_BUY, &model.Decimal64{Mantissa: 2703, Exponent: -2}, 1100, t)

	p = <-oec.submitOrderChan
	verifyParams(p, "XLF", orderentryapi.Side_BUY, &model.Decimal64{Mantissa: 2702, Exponent: -2}, 2500, t)

	p = <-oec.submitOrderChan
	verifyParams(p, "XLF", orderentryapi.Side_BUY, &model.Decimal64{Mantissa: 2701, Exponent: -2}, 2500, t)


	p = <-oec.submitOrderChan
	verifyParams(p, "XLF", orderentryapi.Side_SELL, &model.Decimal64{Mantissa: 2706, Exponent: -2}, 2500, t)

	p = <-oec.submitOrderChan
	verifyParams(p, "XLF", orderentryapi.Side_SELL, &model.Decimal64{Mantissa: 2707, Exponent: -2}, 2500, t)

	p = <-oec.submitOrderChan
	verifyParams(p, "XLF", orderentryapi.Side_SELL, &model.Decimal64{Mantissa: 2708, Exponent: -2}, 2500, t)

}

func verifyParams(p *orderentryapi.NewOrderParams, expS string, s orderentryapi.Side, price *model.Decimal64 , qty int, t *testing.T) {
	if p.Symbol != expS || p.OrderSide != s || !asMd64(p.Price).Equal(price) ||
		!asMd64(p.Quantity).Equal(model.IasD(qty)) {
		t.FailNow()
	}
}

func asMd64( d *orderentryapi.Decimal64) *model.Decimal64{
	return &model.Decimal64{
		Mantissa:             d.Mantissa,
		Exponent:             d.Exponent,
	}
}
