package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
)



type ListingIdSymbol struct {
	ListingId int
	Symbol    string
}

type ClobQuoteSink interface {
	Send(quote *model.ClobQuote)
}

type Actor interface {
	Start() Actor
	Close(chan<- bool)
}

