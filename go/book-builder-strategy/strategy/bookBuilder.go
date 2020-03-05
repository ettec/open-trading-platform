package strategy

import (
	"context"
	"fmt"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/depth"
	"github.com/ettec/open-trading-platform/go/book-builder-strategy/orderentryapi"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/actor"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/google/uuid"
	"log"
	"os"
	"sync"
)

type bookBuilderState int

const (
	Stopped = iota
	Running
)

type bookBuilder struct {
	listing           *model.Listing
	quoteSource       actor.QuoteDistributor
	initialDepth      depth.Depth
	state             bookBuilderState
	orderEntryService orderentryapi.OrderEntryServiceClient
	stateMux          sync.Mutex
	stopChan          chan bool
	log               *log.Logger
}

func newBookBuilder(listing *model.Listing, distributor actor.QuoteDistributor, initialDepth depth.Depth,
	orderEntryService orderentryapi.OrderEntryServiceClient) *bookBuilder {

	b := &bookBuilder{
		log:               log.New(os.Stdout, fmt.Sprintf(" bookBuilder: %v ", listing.Id), log.Ltime),
		listing:           listing,
		quoteSource:       distributor,
		initialDepth:      initialDepth,
		orderEntryService: orderEntryService,
		stopChan:          make(chan bool),
	}

	return b
}

func (b *bookBuilder) stop() error {
	err := b.setState(Stopped)
	if err != nil {
		return err
	}

	b.stopChan<-true
	return nil
}

func (b *bookBuilder) start() error {

	err := b.setState(Running)
	if err != nil {
		return err
	}

	quotesIn := make(chan *model.ClobQuote, 1000)

	b.quoteSource.AddOutQuoteChan(quotesIn)
	defer b.quoteSource.RemoveOutQuoteChan(quotesIn)
	b.quoteSource.Subscribe(b.listing.Id, quotesIn)

	go func() {

		firstQuote := true

	loop:
		for {
			select {
			case q := <-quotesIn:
				if firstQuote {
					firstQuote = false
					// Clear bookBuilder and then submit initial depth
					totalBidQty, worstBid := getTotalQtyAndLeastCompPrice(q.GetBids(), func(l *model.Decimal64, r *model.Decimal64) bool {
						return l.LessThan(r)
					})

					totalAskQty, worstAsk := getTotalQtyAndLeastCompPrice(q.GetBids(), func(l *model.Decimal64, r *model.Decimal64) bool {
						return l.GreaterThan(r)
					})

					uniqueId, _ := uuid.NewUUID()
					b.orderEntryService.SubmitNewOrder(context.Background(), &orderentryapi.NewOrderParams{
						OrderSide: orderentryapi.Side_SELL,
						Quantity:  toApiDec64(totalBidQty),
						Price:     toApiDec64(worstBid),
						Symbol:    b.listing.MarketSymbol,
						ClOrderId: uniqueId.String(),
					})

					uniqueId, _ = uuid.NewUUID()
					b.orderEntryService.SubmitNewOrder(context.Background(), &orderentryapi.NewOrderParams{
						OrderSide: orderentryapi.Side_BUY,
						Quantity:  toApiDec64(totalAskQty),
						Price:     toApiDec64(worstAsk),
						Symbol:    b.listing.MarketSymbol,
						ClOrderId: uniqueId.String(),
					})
				}
			case <-b.stopChan:
				break loop
			}
		}

	}()

	return nil
}

func toApiDec64(d *model.Decimal64) *orderentryapi.Decimal64 {
	return &orderentryapi.Decimal64{
		Mantissa: d.Mantissa,
		Exponent: d.Exponent,
	}
}

func getTotalQtyAndLeastCompPrice(lines []*model.ClobLine, lessCompetitive func(l *model.Decimal64, r *model.Decimal64) bool) (*model.Decimal64, *model.Decimal64) {
	totalBidQty := &model.Decimal64{}
	var worstBid *model.Decimal64

	firstLine := true

	for _, bid := range lines {
		if firstLine {
			firstLine = false
			worstBid = bid.GetPrice()
		}
		totalBidQty.Add(bid.GetSize())

		if lessCompetitive(bid.GetPrice(), worstBid) {
			worstBid = bid.GetPrice()
		}

	}
	return totalBidQty, worstBid
}

func (b *bookBuilder) setState(newState bookBuilderState) error {
	b.stateMux.Lock()
	defer b.stateMux.Unlock()

	if newState == Running {
		if b.state == Running {
			return fmt.Errorf("bookBuilder for listing id %v is already running", b.listing.Id)
		}
		b.state = Running
	}

	return nil
}
