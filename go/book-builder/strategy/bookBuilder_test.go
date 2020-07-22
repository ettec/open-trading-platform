package strategy

import (
	"context"
	"encoding/json"
	"github.com/ettec/open-trading-platform/go/book-builder/depth"
	"github.com/ettec/open-trading-platform/go/book-builder/orderentryapi"
	marketdata "github.com/ettec/otp-mdcommon"

	"github.com/ettec/otp-model"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"testing"
	"time"
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

	return &orderentryapi.OrderId{OrderId: id.String()}, nil
}
func (o *testOrderEntryClient) CancelOrder(ctx context.Context, in *orderentryapi.OrderId, opts ...grpc.CallOption) (*orderentryapi.Empty, error) {
	return nil, nil
}

type testQuoteDist struct {
	subscribeChan    chan int32
	addQuoteChanChan chan chan<- *model.ClobQuote
}

func (t testQuoteDist) GetNewQuoteStream() marketdata.MdsQuoteStream {

	quoteChan := make(chan *model.ClobQuote)

	t.addQuoteChanChan <- quoteChan

	return &testQuoteStream{quoteChan: quoteChan}
}

func newTestQuoteDist() *testQuoteDist {
	t := &testQuoteDist{}
	t.subscribeChan = make(chan int32, 100)
	t.addQuoteChanChan = make(chan chan<- *model.ClobQuote, 100)
	return t
}

type testQuoteStream struct {
	quoteChan chan *model.ClobQuote
}

func (t testQuoteStream) Subscribe(listingId int32) {

}

func (t testQuoteStream) GetStream() <-chan *model.ClobQuote {
	return t.quoteChan
}

func (t testQuoteStream) Close() {

}

func getTestListing() *model.Listing {

	entry := &model.TickSizeEntry{
		LowerPriceBound: model.IasD(0),
		UpperPriceBound: model.IasD(10000000),
		TickSize:        &model.Decimal64{Mantissa: 1, Exponent: -2},
	}

	table := &model.TickSizeTable{
		Entries: []*model.TickSizeEntry{entry},
	}

	return &model.Listing{Id: 1, MarketSymbol: "XLF", TickSize: table}
}

func Test_bookBuilder_start(t *testing.T) {

	dep := depth.Depth{}
	json.Unmarshal([]byte(depthJson), &dep)

	oec := newTestOrderEntryClient()
	qd := newTestQuoteDist()

	book, err := NewBookBuilder(getTestListing(), qd, dep, oec, 10*time.Millisecond, 0.0, 0, 0.9)
	if err != nil {
		t.FailNow()
	}

	book.Start()

	p := <-oec.submitOrderChan
	p = <-oec.submitOrderChan

	sink := <-qd.addQuoteChanChan

	sink <- &model.ClobQuote{
		ListingId: 1,
		Bids: []*model.ClobLine{{Size: model.IasD(100),
			Price: model.IasD(20)},
			{Size: model.IasD(200),
				Price: model.IasD(30)},
		},
		Offers: []*model.ClobLine{{Size: model.IasD(400),
			Price: model.IasD(35)},
			{Size: model.IasD(200),
				Price: model.IasD(25)},
		}}

	p = <-oec.submitOrderChan
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

func verifyParams(p *orderentryapi.NewOrderParams, expS string, s orderentryapi.Side, price *model.Decimal64, qty int, t *testing.T) {
	if p.Symbol != expS || p.OrderSide != s || !asMd64(p.Price).Equal(price) ||
		!asMd64(p.Quantity).Equal(model.IasD(qty)) {
		t.FailNow()
	}
}

func asMd64(d *orderentryapi.Decimal64) *model.Decimal64 {
	return &model.Decimal64{
		Mantissa: d.Mantissa,
		Exponent: d.Exponent,
	}
}

func Test_bookBuilderBooksStats(t *testing.T) {
	dep := depth.Depth{}
	json.Unmarshal([]byte(depthJson), &dep)

	qty, best, worst := getBookStats(dep.Bids, model.Side_BUY)
	if qty != 6100 || best != 27.03 || worst != 27.01 {
		t.FailNow()
	}

	qty, best, worst = getBookStats(dep.Asks, model.Side_SELL)
	if qty != 7500 || best != 27.06 || worst != 27.08 {
		t.FailNow()
	}

}

func Test_bookBuilder_rebuildsBook(t *testing.T) {
	dep := depth.Depth{}
	json.Unmarshal([]byte(depthJson), &dep)

	oec := newTestOrderEntryClient()
	qd := newTestQuoteDist()

	book, err := NewBookBuilder(getTestListing(), qd, dep, oec, 10*time.Millisecond, 0.0, 0.01, 0.9)
	if err != nil {
		t.FailNow()
	}

	book.Start()

	p := <-oec.submitOrderChan
	p = <-oec.submitOrderChan

	sink := <-qd.addQuoteChanChan

	sink <- &model.ClobQuote{
		ListingId: 1,
		Bids:      []*model.ClobLine{},
		Offers:    []*model.ClobLine{}}

	<-oec.submitOrderChan
	<-oec.submitOrderChan
	<-oec.submitOrderChan
	<-oec.submitOrderChan
	<-oec.submitOrderChan
	<-oec.submitOrderChan

	quote := getClobQuote([]line{
		{27.03, 1100},
		{27.01, 2500},
	}, []line{
		{27.06, 2500},
		{27.07, 2500},
		{27.08, 2500},
	}, 1)

	sink <- quote

	p = <-oec.submitOrderChan

	if p.Symbol != "XLF" || p.OrderSide != orderentryapi.Side_BUY {
		t.FailNow()
	}

	if asMd64(p.Quantity).LessThan(model.IasD(0)) || asMd64(p.Quantity).GreaterThan(model.IasD(3000)) {
		t.FailNow()
	}

}

func Test_bookBuilder_tradesOppositeSide(t *testing.T) {

	dep := depth.Depth{}
	json.Unmarshal([]byte(depthJson), &dep)

	oec := newTestOrderEntryClient()
	qd := newTestQuoteDist()

	book, err := NewBookBuilder(getTestListing(), qd, dep, oec, 10*time.Millisecond, 1.0, 0.01, 0.0)
	if err != nil {
		t.FailNow()
	}

	book.Start()

	p := <-oec.submitOrderChan
	p = <-oec.submitOrderChan

	sink := <-qd.addQuoteChanChan

	sink <- &model.ClobQuote{
		ListingId: 1,
		Bids:      []*model.ClobLine{},
		Offers:    []*model.ClobLine{}}

	<-oec.submitOrderChan
	<-oec.submitOrderChan
	<-oec.submitOrderChan
	<-oec.submitOrderChan
	<-oec.submitOrderChan
	<-oec.submitOrderChan

	quote := getClobQuote([]line{
		{27.03, 1100},
		{27.01, 2500},
	}, []line{
		{27.06, 2500},
		{27.07, 2500},
		{27.08, 2500},
	}, 1)

	sink <- quote

	p = <-oec.submitOrderChan

	verifyParams(p, "XLF", orderentryapi.Side_BUY, &model.Decimal64{Mantissa: 2706, Exponent: -2}, 2500, t)

	p = <-oec.submitOrderChan

	verifyParams(p, "XLF", orderentryapi.Side_SELL, &model.Decimal64{Mantissa: 2703, Exponent: -2}, 1100, t)

}

func getClobQuote(rawBids []line, rawAsks []line, listingId int32) *model.ClobQuote {
	bids := toClobLines(rawBids)
	asks := toClobLines(rawAsks)

	quote := &model.ClobQuote{
		ListingId: listingId,
		Bids:      bids,
		Offers:    asks}
	return quote
}

func toClobLines(lines []line) []*model.ClobLine {
	clobLines := []*model.ClobLine{}
	for _, line := range lines {
		clobLines = append(clobLines, &model.ClobLine{
			Size:  model.IasD(line.size),
			Price: model.FasD(line.price),
		})
	}
	return clobLines
}

type line struct {
	price float64
	size  int
}
