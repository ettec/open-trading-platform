package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/connections/fixsim"
	"testing"
)

type mock struct {
	SubscribeMock   func(symbol string)
	fetchSymbolMock func(listingId int, onSymbol chan<- ListingIdSymbol)
}

func (m *mock) Subscribe(symbol string)  {
	 m.SubscribeMock(symbol)
}

func (m *mock) FetchSymbol(listingId int, onSymbol chan<- ListingIdSymbol) {
	m.fetchSymbolMock(listingId, onSymbol)
}

func Test_subscriptionHandler_subscribe(t *testing.T) {

	connectionName := "testConn"

	subscribedSymbols := make(map[string]bool)
	listingToSymbol := map[int]string{1: "A", 2: "B", 3: "C"}


	s := NewSubscriptionHandler(connectionName, func(listingId int, onSymbol chan<- ListingIdSymbol) {
		if symbol, ok := listingToSymbol[listingId]; ok {
			onSymbol <- ListingIdSymbol{ListingId: listingId, Symbol: symbol}
		}
	}, func(symbol string)  {
		subscribedSymbols[symbol] = true

	})

	s.Subscribe(1)
	s.Subscribe(2)

	fixsim.invoke(s.readInputChannels, 4)

	if _, ok := subscribedSymbols["A"]; !ok {
		t.Errorf("expected symbol in Subscribe call")
	}

	if _, ok := subscribedSymbols["B"]; !ok {
		t.Errorf("expected symbol in Subscribe call")
	}

	if len(subscribedSymbols) != 2 {
		t.Errorf("expected 2 symbols in Subscribe call")
	}

	done := make(chan bool)
	s.Close(done)

	if d, _ := s.readInputChannels(); d == nil {
		t.Errorf("expected return close channel")
	}

}


