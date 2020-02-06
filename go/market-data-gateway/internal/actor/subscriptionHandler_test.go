package actor

import (
	"testing"
)

type mock struct {
	SubscribeMock   func(symbol string) error
	fetchSymbolMock func(listingId int, onSymbol chan<- ListingIdSymbol)
}

func (m *mock) Subscribe(symbol string) error {
	return m.SubscribeMock(symbol)
}

func (m *mock) FetchSymbol(listingId int, onSymbol chan<- ListingIdSymbol) {
	m.fetchSymbolMock(listingId, onSymbol)
}

func Test_subscriptionHandler_subscribe(t *testing.T) {

	connectionName := "testConn"

	subscribedSymbols := make(map[string]bool)
	listingToSymbol := map[int]string{1: "A", 2: "B", 3: "C"}

	m := &mock{

		SubscribeMock: func(symbol string) error {
			subscribedSymbols[symbol] = true
			return nil
		},
		fetchSymbolMock: func(listingId int, onSymbol chan<- ListingIdSymbol) {
			if symbol, ok := listingToSymbol[listingId]; ok {
				onSymbol <- ListingIdSymbol{ListingId: listingId, Symbol: symbol}
			}
		},
	}

	s := NewSubscriptionHandler(connectionName, m, m)

	s.Subscribe(1)
	s.Subscribe(2)

	invoke(s.readInputChannels, 4)

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

	if s.readInputChannels() == nil {
		t.Errorf("expected return close channel")
	}

}


