package quoteaggregator

import (
	"context"
	"github.com/ettec/otp-common/api/marketdatasource"
	"github.com/ettec/otp-common/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
	"reflect"
	"testing"
)


type testMdsQuoteStream struct {
	subscribeChan  chan int32
	refreshChan chan *model.ClobQuote
}

func newTestQuoteStream() *testMdsQuoteStream {
	s := &testMdsQuoteStream{}
	s.subscribeChan=  make(chan int32, 10)
	s.refreshChan= make(chan *model.ClobQuote)
	return s
}

func(s *testMdsQuoteStream) Subscribe(listingId int32) {
	s.subscribeChan <- listingId
}

func(s *testMdsQuoteStream) GetStream() <-chan *model.ClobQuote {
	return s.refreshChan
}

func(s *testMdsQuoteStream)Close() {}



func Test_quoteAggregator(t *testing.T) {

	mdsqs := newTestQuoteStream()

	qa := New( func(listingId int32, listingGroupsIn chan<- []*model.Listing) {
		if listingId == 1 {
			listingGroupsIn <- []*model.Listing{{
				Version: 0,
				Id:      1,
				Market:  &model.Market{Mic: "XOSR"},
			},
				{
					Version: 0,
					Id:      2,
					Market:  &model.Market{Mic: "IEXG"},
				},
				{
					Version: 0,
					Id:      3,
					Market:  &model.Market{Mic: "XNAS"},
				},
			}
		}

	}, mdsqs, 1000)


	qa.Subscribe(1)

	listingId := <-mdsqs.subscribeChan

	if listingId != 2 {
		t.FailNow()
	}

	listingId = <-mdsqs.subscribeChan

	if listingId != 3 {
		t.FailNow()
	}

	mdsqs.refreshChan <- &model.ClobQuote{

		ListingId: 2,
		Bids: []*model.ClobLine{
			{Size: d64(15), Price: d64(105)},
			{Size: d64(13), Price: d64(103)},
			{Size: d64(10), Price: d64(100)},
		},
		Offers: []*model.ClobLine{
			{Size: d64(10), Price: d64(100)},
			{Size: d64(13), Price: d64(103)},
			{Size: d64(15), Price: d64(105)},
		},
		StreamInterrupted: false,
		StreamStatusMsg:   "",
	}

	q := <-qa.GetStream()

	firstQuote := &model.ClobQuote{

		ListingId: 1,
		Bids: []*model.ClobLine{
			{Size: d64(15), Price: d64(105), ListingId: 2},
			{Size: d64(13), Price: d64(103), ListingId: 2},
			{Size: d64(10), Price: d64(100), ListingId: 2},
		},
		Offers: []*model.ClobLine{
			{Size: d64(10), Price: d64(100), ListingId: 2},
			{Size: d64(13), Price: d64(103), ListingId: 2},
			{Size: d64(15), Price: d64(105), ListingId: 2},
		},
		StreamInterrupted: false,
		StreamStatusMsg:   "",
		TradedVolume: &model.Decimal64{},
	}

	if !reflect.DeepEqual(firstQuote, q) {
		t.FailNow()
	}

	mdsqs.refreshChan <- &model.ClobQuote{
		ListingId: 3,
		Bids: []*model.ClobLine{
			{Size: d64(13), Price: d64(104)},
			{Size: d64(12), Price: d64(102)},
			{Size: d64(11), Price: d64(101)},
		},
		Offers: []*model.ClobLine{
			{Size: d64(11), Price: d64(101)},
			{Size: d64(12), Price: d64(102)},
			{Size: d64(13), Price: d64(104)},
		},
		StreamInterrupted: false,
		StreamStatusMsg:   "",
	}

	q = <-qa.GetStream()

	combinedQuote := &model.ClobQuote{
		ListingId: 1,
		Bids: []*model.ClobLine{
			{Size: d64(15), Price: d64(105), ListingId: 2},
			{Size: d64(13), Price: d64(104), ListingId: 3},
			{Size: d64(13), Price: d64(103), ListingId: 2},
			{Size: d64(12), Price: d64(102), ListingId: 3},
			{Size: d64(11), Price: d64(101), ListingId: 3},
			{Size: d64(10), Price: d64(100), ListingId: 2},
		},
		Offers: []*model.ClobLine{
			{Size: d64(10), Price: d64(100), ListingId: 2},
			{Size: d64(11), Price: d64(101), ListingId: 3},
			{Size: d64(12), Price: d64(102), ListingId: 3},
			{Size: d64(13), Price: d64(103), ListingId: 2},
			{Size: d64(13), Price: d64(104), ListingId: 3},
			{Size: d64(15), Price: d64(105), ListingId: 2},
		},
		StreamInterrupted: false,
		StreamStatusMsg:   "",
		TradedVolume: &model.Decimal64{},
	}

	if !reflect.DeepEqual(combinedQuote, q) {
		t.FailNow()
	}

	mdsqs.refreshChan <- &model.ClobQuote{

		ListingId: 2,
		Bids: []*model.ClobLine{
			{Size: d64(15), Price: d64(106)},
			{Size: d64(13), Price: d64(103)},
			{Size: d64(10), Price: d64(100)},
		},
		Offers: []*model.ClobLine{
			{Size: d64(10), Price: d64(100)},
			{Size: d64(13), Price: d64(103)},
			{Size: d64(15), Price: d64(106)},
		},
		StreamInterrupted: false,
		StreamStatusMsg:   "",

	}

	q = <-qa.GetStream()

	combinedQuote = &model.ClobQuote{
		ListingId: 1,
		Bids: []*model.ClobLine{
			{Size: d64(15), Price: d64(106), ListingId: 2},
			{Size: d64(13), Price: d64(104), ListingId: 3},
			{Size: d64(13), Price: d64(103), ListingId: 2},
			{Size: d64(12), Price: d64(102), ListingId: 3},
			{Size: d64(11), Price: d64(101), ListingId: 3},
			{Size: d64(10), Price: d64(100), ListingId: 2},
		},
		Offers: []*model.ClobLine{
			{Size: d64(10), Price: d64(100), ListingId: 2},
			{Size: d64(11), Price: d64(101), ListingId: 3},
			{Size: d64(12), Price: d64(102), ListingId: 3},
			{Size: d64(13), Price: d64(103), ListingId: 2},
			{Size: d64(13), Price: d64(104), ListingId: 3},
			{Size: d64(15), Price: d64(106), ListingId: 2},
		},
		StreamInterrupted: false,
		StreamStatusMsg:   "",
		TradedVolume: &model.Decimal64{},
	}

	if !reflect.DeepEqual(combinedQuote, q) {
		t.FailNow()
	}

}

type testConnection struct {
	state        connectivity.State
	getStateChan chan connectivity.State
}

func (t testConnection) GetState() connectivity.State {
	t.state = <-t.getStateChan
	return t.state
}

func (t testConnection) WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool {

	for {
		if t.state != sourceState {
			return true
		}
		t.state = <-t.getStateChan
	}

}

type testClient struct {
	streamOutChan chan marketdatasource.MarketDataSource_ConnectClient
}

func (t testClient) Connect(ctx context.Context, opts ...grpc.CallOption) (marketdatasource.MarketDataSource_ConnectClient, error) {
	return <-t.streamOutChan, nil
}

type testClientStream struct {
	refreshChan    chan *model.ClobQuote
	refreshErrChan chan error
	subscribeChan  chan *marketdatasource.SubscribeRequest
}

func (t testClientStream) Send(request *marketdatasource.SubscribeRequest) error {
	t.subscribeChan <- request
	return nil
}

func (t testClientStream) Recv() (*model.ClobQuote, error) {
	select {
	case r := <-t.refreshChan:
		return r, nil
	case e := <-t.refreshErrChan:
		return nil, e
	}

	return <-t.refreshChan, nil
}

func (t testClientStream) Header() (metadata.MD, error) {
	panic("implement me")
}

func (t testClientStream) Trailer() metadata.MD {
	panic("implement me")
}

func (t testClientStream) CloseSend() error {
	panic("implement me")
}

func (t testClientStream) Context() context.Context {
	panic("implement me")
}

func (t testClientStream) SendMsg(m interface{}) error {
	panic("implement me")
}

func (t testClientStream) RecvMsg(m interface{}) error {
	panic("implement me")
}



func Test_combineQuotes(t *testing.T) {
	type args struct {
		combinedListingId int32
		quotes            []*model.ClobQuote
	}

	tests := []struct {
		name string
		args args
		want *model.ClobQuote
	}{
		{name: "test combine 2 quotes",
			args: args{combinedListingId: 1, quotes: []*model.ClobQuote{
				{
					ListingId: 2,
					Bids: []*model.ClobLine{
						{Size: d64(15), Price: d64(105)},
						{Size: d64(13), Price: d64(103)},
						{Size: d64(10), Price: d64(100)},
					},
					Offers: []*model.ClobLine{
						{Size: d64(10), Price: d64(100)},
						{Size: d64(13), Price: d64(103)},
						{Size: d64(15), Price: d64(105)},
					},
					StreamInterrupted: false,
					StreamStatusMsg:   "",
				},
				{
					ListingId: 3,
					Bids: []*model.ClobLine{
						{Size: d64(13), Price: d64(103)},
						{Size: d64(12), Price: d64(102)},
						{Size: d64(11), Price: d64(101)},
					},
					Offers: []*model.ClobLine{
						{Size: d64(11), Price: d64(101)},
						{Size: d64(12), Price: d64(102)},
						{Size: d64(13), Price: d64(103)},
					},
					StreamInterrupted: false,
					StreamStatusMsg:   "",
				},
			}},

			want: &model.ClobQuote{
				ListingId: 1,
				Bids: []*model.ClobLine{
					{Size: d64(15), Price: d64(105), ListingId: 2},
					{Size: d64(13), Price: d64(103), ListingId: 2},
					{Size: d64(13), Price: d64(103), ListingId: 3},
					{Size: d64(12), Price: d64(102), ListingId: 3},
					{Size: d64(11), Price: d64(101), ListingId: 3},
					{Size: d64(10), Price: d64(100), ListingId: 2},
				},
				Offers: []*model.ClobLine{
					{Size: d64(10), Price: d64(100), ListingId: 2},
					{Size: d64(11), Price: d64(101), ListingId: 3},
					{Size: d64(12), Price: d64(102), ListingId: 3},
					{Size: d64(13), Price: d64(103), ListingId: 2},
					{Size: d64(13), Price: d64(103), ListingId: 3},
					{Size: d64(15), Price: d64(105), ListingId: 2},
				},
				StreamInterrupted: false,
				StreamStatusMsg:   "",
				TradedVolume: &model.Decimal64{},
			}},
		{name: "test combine 3 quote",
			args: args{combinedListingId: 1, quotes: []*model.ClobQuote{
				{
					ListingId: 2,
					Bids: []*model.ClobLine{
						{Size: d64(15), Price: d64(105)},
						{Size: d64(13), Price: d64(103)},
						{Size: d64(10), Price: d64(100)},
					},
					Offers: []*model.ClobLine{
						{Size: d64(10), Price: d64(100)},
						{Size: d64(13), Price: d64(103)},
						{Size: d64(15), Price: d64(105)},
					},
					StreamInterrupted: false,
					StreamStatusMsg:   "",
				},
				{
					ListingId: 3,
					Bids: []*model.ClobLine{
						{Size: d64(13), Price: d64(103)},
						{Size: d64(12), Price: d64(102)},
						{Size: d64(11), Price: d64(101)},
					},
					Offers: []*model.ClobLine{
						{Size: d64(11), Price: d64(101)},
						{Size: d64(12), Price: d64(102)},
						{Size: d64(13), Price: d64(103)},
					},
					StreamInterrupted: false,
					StreamStatusMsg:   "",
				},
				{
					ListingId: 4,
					Bids: []*model.ClobLine{
						{Size: d64(16), Price: d64(106)},
						{Size: d64(12), Price: d64(102)},
						{Size: d64(9), Price: d64(99)},
					},
					Offers: []*model.ClobLine{
						{Size: d64(9), Price: d64(99)},
						{Size: d64(12), Price: d64(102)},
						{Size: d64(16), Price: d64(106)},
					},
					StreamInterrupted: false,
					StreamStatusMsg:   "",
				},
			}},

			want: &model.ClobQuote{
				ListingId: 1,
				Bids: []*model.ClobLine{
					{Size: d64(16), Price: d64(106), ListingId: 4},
					{Size: d64(15), Price: d64(105), ListingId: 2},
					{Size: d64(13), Price: d64(103), ListingId: 2},
					{Size: d64(13), Price: d64(103), ListingId: 3},
					{Size: d64(12), Price: d64(102), ListingId: 3},
					{Size: d64(12), Price: d64(102), ListingId: 4},
					{Size: d64(11), Price: d64(101), ListingId: 3},
					{Size: d64(10), Price: d64(100), ListingId: 2},
					{Size: d64(9), Price: d64(99), ListingId: 4},
				},
				Offers: []*model.ClobLine{
					{Size: d64(9), Price: d64(99), ListingId: 4},
					{Size: d64(10), Price: d64(100), ListingId: 2},
					{Size: d64(11), Price: d64(101), ListingId: 3},
					{Size: d64(12), Price: d64(102), ListingId: 3},
					{Size: d64(12), Price: d64(102), ListingId: 4},
					{Size: d64(13), Price: d64(103), ListingId: 2},
					{Size: d64(13), Price: d64(103), ListingId: 3},
					{Size: d64(15), Price: d64(105), ListingId: 2},
					{Size: d64(16), Price: d64(106), ListingId: 4},
				},
				StreamInterrupted: false,
				StreamStatusMsg:   "",
				TradedVolume: &model.Decimal64{},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combineQuotes(tt.args.combinedListingId, tt.args.quotes, tt.args.quotes[len(tt.args.quotes)-1]); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combineQuotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_combineQuoteStatus(t *testing.T) {
	type args struct {
		combinedListingId int32
		quotes            []*model.ClobQuote
	}

	tests := []struct {
		name string
		args args
		want *model.ClobQuote
	}{
		{name: "test combine 2 quotes",
			args: args{combinedListingId: 1, quotes: []*model.ClobQuote{
				{
					ListingId: 2,
					Bids: []*model.ClobLine{
						{Size: d64(15), Price: d64(105)},
					},
					Offers: []*model.ClobLine{
						{Size: d64(10), Price: d64(100)},
					},
					StreamInterrupted: true,
					StreamStatusMsg:   "connection lost",
				},
				{
					ListingId: 3,
					Bids: []*model.ClobLine{
						{Size: d64(13), Price: d64(103)},
					},
					Offers: []*model.ClobLine{
						{Size: d64(11), Price: d64(101)},
					},
					StreamInterrupted: false,
					StreamStatusMsg:   "",
				},
			}},

			want: &model.ClobQuote{
				ListingId: 1,
				Bids: []*model.ClobLine{
					{Size: d64(15), Price: d64(105), ListingId: 2},
					{Size: d64(13), Price: d64(103), ListingId: 3},
				},
				Offers: []*model.ClobLine{
					{Size: d64(10), Price: d64(100), ListingId: 2},
					{Size: d64(11), Price: d64(101), ListingId: 3},
				},
				StreamInterrupted: true,
				StreamStatusMsg:   "connection lost",
				TradedVolume: &model.Decimal64{},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combineQuotes(tt.args.combinedListingId, tt.args.quotes, tt.args.quotes[len(tt.args.quotes)-1]); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combineQuotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_combineTradeData(t *testing.T) {
	type args struct {
		combinedListingId int32
		quotes            []*model.ClobQuote
	}

	tests := []struct {
		name string
		args args
		want *model.ClobQuote
	}{
		{name: "test combine 2 quotes",
			args: args{combinedListingId: 1, quotes: []*model.ClobQuote{
				{
					ListingId: 2,
					Bids: []*model.ClobLine{
						{Size: d64(15), Price: d64(105)},
					},
					Offers: []*model.ClobLine{
						{Size: d64(10), Price: d64(100)},
					},
					StreamInterrupted: true,
					StreamStatusMsg:   "connection lost",
					LastPrice:  &model.Decimal64{Mantissa: 50, Exponent: 0},
					LastQuantity:  &model.Decimal64{Mantissa: 15, Exponent: 0},
					TradedVolume: &model.Decimal64{Mantissa: 10, Exponent: 0},
				},
				{
					ListingId: 3,
					Bids: []*model.ClobLine{
						{Size: d64(13), Price: d64(103)},
					},
					Offers: []*model.ClobLine{
						{Size: d64(11), Price: d64(101)},
					},
					StreamInterrupted: false,
					StreamStatusMsg:   "",
					LastPrice:  &model.Decimal64{Mantissa: 60, Exponent: 0},
					LastQuantity:  &model.Decimal64{Mantissa: 25, Exponent: 0},
					TradedVolume: &model.Decimal64{Mantissa: 20, Exponent: 0},
				},
			}},

			want: &model.ClobQuote{
				ListingId: 1,
				Bids: []*model.ClobLine{
					{Size: d64(15), Price: d64(105), ListingId: 2},
					{Size: d64(13), Price: d64(103), ListingId: 3},
				},
				Offers: []*model.ClobLine{
					{Size: d64(10), Price: d64(100), ListingId: 2},
					{Size: d64(11), Price: d64(101), ListingId: 3},
				},
				StreamInterrupted: true,
				StreamStatusMsg:   "connection lost",
				LastPrice:  &model.Decimal64{Mantissa: 60, Exponent: 0},
				LastQuantity:  &model.Decimal64{Mantissa: 25, Exponent: 0},
				TradedVolume: &model.Decimal64{Mantissa: 30, Exponent: 0},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combineQuotes(tt.args.combinedListingId, tt.args.quotes, tt.args.quotes[len(tt.args.quotes)-1]); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combineQuotes() = %v, want %v", got, tt.want)
			}
		})
	}
}



func d64(mantissa int) *model.Decimal64 {
	return &model.Decimal64{Mantissa: int64(mantissa), Exponent: 0}
}
