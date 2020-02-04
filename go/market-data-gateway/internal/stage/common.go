package stage

import "github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"

type Refresh marketdata.MarketDataIncrementalRefresh


type ListingIdSymbol struct {
	ListingId int
	Symbol    string
}