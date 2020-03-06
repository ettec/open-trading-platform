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

	b.stopChan <- true
	return nil
}

func (b *bookBuilder) start() error {

	err := b.setState(Running)
	if err != nil {
		return err
	}



	go func() {

		quotesIn := make(chan *model.ClobQuote, 1000)
		b.quoteSource.AddOutQuoteChan(quotesIn)
		defer b.quoteSource.RemoveOutQuoteChan(quotesIn)
		b.quoteSource.Subscribe(b.listing.Id, quotesIn)

		firstQuote := true

	loop:
		for {
			select {
			case q := <-quotesIn:
				if firstQuote {
					firstQuote = false

					b.clearBook(q)

					b.sendOrdersForLines(b.initialDepth.Bids, orderentryapi.Side_BUY)
					b.sendOrdersForLines(b.initialDepth.Asks, orderentryapi.Side_SELL)

				}
			case <-b.stopChan:
				break loop
			}
		}

	}()

	return nil
}

func (b *bookBuilder) sendOrdersForLines(bids []struct {
	Price     float64 `json:"price"`
	Size      int     `json:"size"`
	Timestamp int64   `json:"timestamp"`
}, side orderentryapi.Side) {
	for _, bid := range bids {

		uniqueId, _ := uuid.NewUUID()
		b.orderEntryService.SubmitNewOrder(context.Background(), &orderentryapi.NewOrderParams{
			OrderSide: side,
			Quantity:  &orderentryapi.Decimal64{Mantissa: int64(bid.Size), Exponent: 0},
			Price:     toApiDec64(model.NewFromFloat(bid.Price)),
			Symbol:    b.listing.MarketSymbol,
			ClOrderId: uniqueId.String(),
		})

	}
}

func (b *bookBuilder) clearBook(q *model.ClobQuote) {
	totalBidQty, worstBid := getTotalQtyAndLeastCompetitivePrice(q.GetBids(), func(l *model.Decimal64, r *model.Decimal64) bool {
		return l.LessThan(r)
	})

	totalAskQty, worstAsk := getTotalQtyAndLeastCompetitivePrice(q.GetOffers(), func(l *model.Decimal64, r *model.Decimal64) bool {
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

func toApiDec64(d *model.Decimal64) *orderentryapi.Decimal64 {
	return &orderentryapi.Decimal64{
		Mantissa: d.Mantissa,
		Exponent: d.Exponent,
	}
}

func getTotalQtyAndLeastCompetitivePrice(lines []*model.ClobLine, lessCompetitive func(l *model.Decimal64, r *model.Decimal64) bool) (*model.Decimal64, *model.Decimal64) {
	totalQty := &model.Decimal64{}
	var worstPrice *model.Decimal64

	firstLine := true

	for _, line := range lines {
		if firstLine {
			firstLine = false
			worstPrice = line.GetPrice()
		}
		totalQty.Add(line.GetSize())

		if lessCompetitive(line.GetPrice(), worstPrice) {
			worstPrice = line.GetPrice()
		}

	}
	return totalQty, worstPrice
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
