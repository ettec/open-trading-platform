package main

import (
	"context"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"testing"
)

type  mock struct {
	SubscribeMock func (ctx context.Context, in *marketdata.MarketDataRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	fetchSymbolMock func (listingId int, onSymbol chan<- listingIdSymbol)

}

func (m *mock)Subscribe(ctx context.Context, in *marketdata.MarketDataRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	return m.SubscribeMock(ctx, in, opts...)
}

func (m *mock)fetchSymbol(listingId int, onSymbol chan<- listingIdSymbol) {
	m.fetchSymbolMock(listingId, onSymbol)
}


func Test_subscriptionHandler_subscribe(t *testing.T) {

	connectionName := "testconn"

	subscribedSymbols := make(map[string]bool)
	listingToSymbol := map[int]string{1:"A", 2:"B", 3:"C"}

	m := &mock{

		SubscribeMock: func(ctx context.Context, in *marketdata.MarketDataRequest, opts ...grpc.CallOption) (e *empty.Empty, err error) {
			if in.Parties[0].PartyId  != connectionName {
				t.Error("expected subscription to have connection name:", connectionName)
			}

			subscribedSymbols[in.InstrmtMdReqGrp[0].Instrument.Symbol] = true
			return &empty.Empty{}, nil
		},
		fetchSymbolMock: func(listingId int, onSymbol chan<- listingIdSymbol) {
			if symbol, ok := listingToSymbol[listingId]; ok {
				onSymbol <- listingIdSymbol{listingId, symbol}
			}
		},
	}

	s := newSubscriptionHandler(connectionName, m ,m )

	s.subscribe(1)
	s.subscribe(2)

	invoke(s.readInputChannels, 4)

	if _, ok := subscribedSymbols["A"]; !ok {
		t.Errorf("expected symbol in subscribe call")
	}

	s.close()
	err := s.readInputChannels()
	if err != closed {
		t.Errorf("expected loop to close")
	}


}

func invoke( f func() error, times int) {

	for i:=0; i<times;i++ {
		f()
	}

}



