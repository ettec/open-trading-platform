package marketdatasource

import (
	"context"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/staticdata"
	"github.com/ettech/open-trading-platform/go/market-data/market-data-service/marketdatasource/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/quotestream.go -package mocks github.com/ettec/otp-common/marketdata QuoteStream

func TestConnectAndSubscribe(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var getListing getListingFn = func(ctx context.Context, listingId int32, result chan<- staticdata.ListingResult) {
		result <- staticdata.ListingResult{Listing: &model.Listing{Id: listingId, Market: &model.Market{Mic: "XTST"}}}
	}

	inboundQuotes := make(chan *model.ClobQuote, 100)
	quoteStream := mocks.NewMockQuoteStream(mockCtrl)
	quoteStream.EXPECT().Chan().Return(inboundQuotes)
	quoteStream.EXPECT().Subscribe(int32(1))

	gatewayStreamSource := mocks.NewMockGatewayStreamSource(mockCtrl)
	gatewayStreamSource.EXPECT().NewQuoteStreamFromMdSource(ctx, "testMds", "testAddress", 0*time.Second, 100).Return(quoteStream, nil)

	mds := NewMarketDataService(ctx, "testMds", gatewayStreamSource, getListing, 100, 0, 100)
	err := mds.AddMarketDataGateway(TestMarketDataGateway{address: "testAddress", ordinal: 1, marketMic: "XTST"})
	assert.NoError(t, err)

	stream := mds.Connect(ctx, "testSubscriber")

	err = stream.Subscribe(1)
	assert.NoError(t, err)

	inboundQuotes <- &model.ClobQuote{ListingId: 1}
	inboundQuotes <- &model.ClobQuote{ListingId: 2}
	inboundQuotes <- &model.ClobQuote{ListingId: 3}
	inboundQuotes <- &model.ClobQuote{ListingId: 1}

	received := <-stream.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 1}, received)
	received = <-stream.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 1}, received)
}

func TestClosingStream(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var getListing getListingFn = func(ctx context.Context, listingId int32, result chan<- staticdata.ListingResult) {
		result <- staticdata.ListingResult{Listing: &model.Listing{Id: listingId, Market: &model.Market{Mic: "XTST"}}}
	}

	inboundQuotes := make(chan *model.ClobQuote, 100)
	quoteStream := mocks.NewMockQuoteStream(mockCtrl)
	quoteStream.EXPECT().Chan().Return(inboundQuotes)
	quoteStream.EXPECT().Subscribe(int32(1))

	gatewayStreamSource := mocks.NewMockGatewayStreamSource(mockCtrl)
	gatewayStreamSource.EXPECT().NewQuoteStreamFromMdSource(ctx, "testMds", "testAddress", 0*time.Second, 100).Return(quoteStream, nil)

	mds := NewMarketDataService(ctx, "testMds", gatewayStreamSource, getListing, 100, 0, 100)
	err := mds.AddMarketDataGateway(TestMarketDataGateway{address: "testAddress", ordinal: 1, marketMic: "XTST"})
	assert.NoError(t, err)

	stream := mds.Connect(ctx, "testSubscriber")

	err = stream.Subscribe(1)
	assert.NoError(t, err)

	inboundQuotes <- &model.ClobQuote{ListingId: 1}
	inboundQuotes <- &model.ClobQuote{ListingId: 2}
	inboundQuotes <- &model.ClobQuote{ListingId: 3}
	inboundQuotes <- &model.ClobQuote{ListingId: 1}

	received := <-stream.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 1}, received)
	received = <-stream.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 1}, received)

	stream.Close()

	inboundQuotes <- &model.ClobQuote{ListingId: 1}

	timer := time.NewTimer(1 * time.Second)
	select {
	case <-stream.Chan():
		t.Errorf("no more quotes expected")
	case <-timer.C:
	}
}

func TestConnectAndSubscribeFromMultipleClients(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var getListing getListingFn = func(ctx context.Context, listingId int32, result chan<- staticdata.ListingResult) {
		result <- staticdata.ListingResult{Listing: &model.Listing{Id: listingId, Market: &model.Market{Mic: "XTST"}}}
	}

	inboundQuotes := make(chan *model.ClobQuote, 100)
	quoteStream := mocks.NewMockQuoteStream(mockCtrl)
	quoteStream.EXPECT().Chan().Return(inboundQuotes)
	quoteStream.EXPECT().Subscribe(int32(1))

	gatewayStreamSource := mocks.NewMockGatewayStreamSource(mockCtrl)
	gatewayStreamSource.EXPECT().NewQuoteStreamFromMdSource(ctx, "testMds", "testAddress", 0*time.Second, 100).Return(quoteStream, nil)

	mds := NewMarketDataService(ctx, "testMds", gatewayStreamSource, getListing, 100, 0, 100)
	err := mds.AddMarketDataGateway(TestMarketDataGateway{address: "testAddress", ordinal: 1, marketMic: "XTST"})
	assert.NoError(t, err)

	stream1 := mds.Connect(ctx, "testSubscriber1")
	stream2 := mds.Connect(ctx, "testSubscriber2")

	err = stream1.Subscribe(1)
	assert.NoError(t, err)
	err = stream2.Subscribe(1)
	assert.NoError(t, err)

	inboundQuotes <- &model.ClobQuote{ListingId: 1}
	inboundQuotes <- &model.ClobQuote{ListingId: 2}
	inboundQuotes <- &model.ClobQuote{ListingId: 3}
	inboundQuotes <- &model.ClobQuote{ListingId: 1}

	received := <-stream1.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 1}, received)
	received = <-stream1.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 1}, received)

	received = <-stream2.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 1}, received)
	received = <-stream2.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 1}, received)
}

func TestSubscriptionsAcrossMarketSentToCorrectGateways(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var getListing getListingFn = func(ctx context.Context, listingId int32, result chan<- staticdata.ListingResult) {
		switch listingId {
		case 1:
			result <- staticdata.ListingResult{Listing: &model.Listing{Id: listingId, Market: &model.Market{Mic: "XTST"}}}
		case 2:
			result <- staticdata.ListingResult{Listing: &model.Listing{Id: listingId, Market: &model.Market{Mic: "XTST2"}}}
		default:
			t.Errorf("unexpected listing id %v", listingId)
		}
	}

	inboundQuotes1 := make(chan *model.ClobQuote, 100)
	quoteStream1 := mocks.NewMockQuoteStream(mockCtrl)
	quoteStream1.EXPECT().Chan().Return(inboundQuotes1)
	quoteStream1.EXPECT().Subscribe(int32(1))

	inboundQuotes2 := make(chan *model.ClobQuote, 100)
	quoteStream2 := mocks.NewMockQuoteStream(mockCtrl)
	quoteStream2.EXPECT().Chan().Return(inboundQuotes2)
	quoteStream2.EXPECT().Subscribe(int32(2))

	gatewayStreamSource1 := mocks.NewMockGatewayStreamSource(mockCtrl)
	gatewayStreamSource1.EXPECT().NewQuoteStreamFromMdSource(ctx, "testMds", "testAddress1", 0*time.Second, 100).Return(quoteStream1, nil)
	gatewayStreamSource1.EXPECT().NewQuoteStreamFromMdSource(ctx, "testMds", "testAddress2", 0*time.Second, 100).Return(quoteStream2, nil)

	mds := NewMarketDataService(ctx, "testMds", gatewayStreamSource1, getListing, 100, 0, 100)
	err := mds.AddMarketDataGateway(TestMarketDataGateway{address: "testAddress1", ordinal: 0, marketMic: "XTST"})
	assert.NoError(t, err)

	err = mds.AddMarketDataGateway(TestMarketDataGateway{address: "testAddress2", ordinal: 0, marketMic: "XTST2"})
	assert.NoError(t, err)

	stream := mds.Connect(ctx, "testSubscriber")

	err = stream.Subscribe(1)
	assert.NoError(t, err)
	err = stream.Subscribe(2)
	assert.NoError(t, err)

	inboundQuotes1 <- &model.ClobQuote{ListingId: 1}
	inboundQuotes1 <- &model.ClobQuote{ListingId: 2}
	inboundQuotes1 <- &model.ClobQuote{ListingId: 3}
	inboundQuotes1 <- &model.ClobQuote{ListingId: 1}

	received := <-stream.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 1}, received)
	received = <-stream.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 1}, received)

	inboundQuotes2 <- &model.ClobQuote{ListingId: 1}
	inboundQuotes2 <- &model.ClobQuote{ListingId: 2}
	inboundQuotes2 <- &model.ClobQuote{ListingId: 3}
	inboundQuotes2 <- &model.ClobQuote{ListingId: 1}

	received = <-stream.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 2}, received)

}

func TestSubscriptionsAreLoadBalancedAcrossGatewaysForTheSameMarket(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var getListing getListingFn = func(ctx context.Context, listingId int32, result chan<- staticdata.ListingResult) {
		result <- staticdata.ListingResult{Listing: &model.Listing{Id: listingId, Market: &model.Market{Mic: "XTST"}}}
	}

	inboundQuotes1 := make(chan *model.ClobQuote, 100)
	quoteStream1 := mocks.NewMockQuoteStream(mockCtrl)
	quoteStream1.EXPECT().Chan().Return(inboundQuotes1)
	quoteStream1.EXPECT().Subscribe(int32(2))

	inboundQuotes2 := make(chan *model.ClobQuote, 100)
	quoteStream2 := mocks.NewMockQuoteStream(mockCtrl)
	quoteStream2.EXPECT().Chan().Return(inboundQuotes2)
	quoteStream2.EXPECT().Subscribe(int32(1))

	gatewayStreamSource1 := mocks.NewMockGatewayStreamSource(mockCtrl)
	gatewayStreamSource1.EXPECT().NewQuoteStreamFromMdSource(ctx, "testMds", "testAddress1", 0*time.Second, 100).Return(quoteStream1, nil)
	gatewayStreamSource1.EXPECT().NewQuoteStreamFromMdSource(ctx, "testMds", "testAddress2", 0*time.Second, 100).Return(quoteStream2, nil)

	mds := NewMarketDataService(ctx, "testMds", gatewayStreamSource1, getListing, 100, 0, 100)
	err := mds.AddMarketDataGateway(TestMarketDataGateway{address: "testAddress1", ordinal: 0, marketMic: "XTST"})
	assert.NoError(t, err)

	err = mds.AddMarketDataGateway(TestMarketDataGateway{address: "testAddress2", ordinal: 1, marketMic: "XTST"})
	assert.NoError(t, err)

	stream := mds.Connect(ctx, "testSubscriber")

	err = stream.Subscribe(1)
	assert.NoError(t, err)
	err = stream.Subscribe(2)
	assert.NoError(t, err)

	inboundQuotes1 <- &model.ClobQuote{ListingId: 1}
	inboundQuotes1 <- &model.ClobQuote{ListingId: 2}
	inboundQuotes1 <- &model.ClobQuote{ListingId: 3}
	inboundQuotes1 <- &model.ClobQuote{ListingId: 1}

	received := <-stream.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 2}, received)

	inboundQuotes2 <- &model.ClobQuote{ListingId: 1}
	inboundQuotes2 <- &model.ClobQuote{ListingId: 2}
	inboundQuotes2 <- &model.ClobQuote{ListingId: 3}
	inboundQuotes2 <- &model.ClobQuote{ListingId: 1}

	received = <-stream.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 1}, received)
	received = <-stream.Chan()
	assert.Equal(t, &model.ClobQuote{ListingId: 1}, received)

}

type TestMarketDataGateway struct {
	address   string
	ordinal   int
	marketMic string
}

func (t TestMarketDataGateway) GetAddress() string {
	return t.address
}

func (t TestMarketDataGateway) GetOrdinal() int {
	return t.ordinal
}

func (t TestMarketDataGateway) GetMarketMic() string {
	return t.marketMic
}
