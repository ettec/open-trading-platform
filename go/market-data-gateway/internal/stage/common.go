package stage

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
)

type Refresh marketdata.MarketDataIncrementalRefresh


type ListingIdSymbol struct {
	ListingId int
	Symbol    string
}

type ClobQuoteSink interface {
	Send(quote *model.ClobQuote)
}

