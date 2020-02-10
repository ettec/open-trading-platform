package connections

import "github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"

type Connection interface {

	Connect() (<-chan *model.ClobQuote, error)
	Subscribe(listingId int)
	Close() error
}

