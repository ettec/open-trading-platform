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
	"math/rand"
	"os"
	"sync"
	"time"
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
	bookScanInterval  time.Duration
	tradeProbability  float32
	variation         float64
	minQty            float64
	log               *log.Logger
	errLog            *log.Logger
}

func newBookBuilder(listing *model.Listing, distributor actor.QuoteDistributor, initialDepth depth.Depth,
	orderEntryService orderentryapi.OrderEntryServiceClient,
	bookScanInterval time.Duration, tradeProbability float32, variation float64, minQtyPercent float64) (*bookBuilder, error) {

	b := &bookBuilder{
		log:               log.New(os.Stdout, fmt.Sprintf(" bookBuilder: %v ", listing.Id), log.Ltime),
		errLog:            log.New(os.Stderr, fmt.Sprintf(" bookBuilder: %v ", listing.Id), log.Ltime),
		listing:           listing,
		quoteSource:       distributor,
		initialDepth:      initialDepth,
		orderEntryService: orderEntryService,
		stopChan:          make(chan bool),
		bookScanInterval:  bookScanInterval,
		tradeProbability:  tradeProbability,
		minQty:            minQtyPercent,
		variation:         variation,
	}

	if len(b.initialDepth.Bids) == 0 {
		return nil, fmt.Errorf("initial depth for listing id %v, symbol %v has no bids", listing.Id, listing.GetMarketSymbol())
	}

	if len(b.initialDepth.Asks) == 0 {
		return nil, fmt.Errorf("initial depth for listing id %v, symbol %v has no asks", listing.Id, listing.GetMarketSymbol())
	}

	return b, nil
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

		ticker := time.NewTicker(500 * time.Millisecond)

		bidsQty, _, _ := getBookStats(b.initialDepth.Bids, model.Side_BUY)
		asksQty, _, _ := getBookStats(b.initialDepth.Asks, model.Side_SELL)

		var lastQuote *model.ClobQuote

	loop:
		for {
			select {
			case q := <-quotesIn:
				lastQuote = q
				if firstQuote {
					firstQuote = false

					b.clearBook(q)
					b.sendOrdersForLines(b.initialDepth.Bids, orderentryapi.Side_BUY)
					b.sendOrdersForLines(b.initialDepth.Asks, orderentryapi.Side_SELL)
				}
			case <-ticker.C:
				if lastQuote != nil {
					b.updateBookSide(orderentryapi.Side_BUY, bidsQty, b.initialDepth.Bids,
						lastQuote.Bids, lastQuote.Offers)

					b.updateBookSide(orderentryapi.Side_SELL, asksQty, b.initialDepth.Asks,
						lastQuote.Offers, lastQuote.Bids)

				}
			case <-b.stopChan:
				break loop
			}
		}

	}()

	return nil
}

func (b *bookBuilder) updateBookSide(side orderentryapi.Side, totalInitialQty float64, initialDepth []struct {
	Price     float64 `json:"price"`
	Size      int     `json:"size"`
	Timestamp int64   `json:"timestamp"`
}, lastQuoteSameSide []*model.ClobLine, lastQuoteOppositeSide []*model.ClobLine) {

	qQty, _, _ := getQuoteStats(lastQuoteSameSide, side)

	if qQty < totalInitialQty*b.minQty {

		idx := rand.Intn(len(initialDepth))
		line := initialDepth[idx]

		price := line.Price - (line.Price * rand.Float64() * b.variation)
		qty := float64(line.Size) - (float64(line.Size) * rand.Float64() * b.variation)

		roundedPrice, err := b.listing.RoundToNearestTick(price)
		if err != nil {
			panic(err)
		}

		roundedQty := b.listing.RoundToLotSize(qty)

		uniqueId, _ := uuid.NewUUID()
		b.orderEntryService.SubmitNewOrder(context.Background(), &orderentryapi.NewOrderParams{
			OrderSide: side,
			Quantity:  toApiDec64(roundedQty),
			Price:     toApiDec64(roundedPrice),
			Symbol:    b.listing.MarketSymbol,
			ClOrderId: uniqueId.String(),
		})
	}

	if rand.Float32() < b.tradeProbability {
		if len(lastQuoteOppositeSide) > 0 {
			bestOpp := lastQuoteOppositeSide[0]
			uniqueId, _ := uuid.NewUUID()
			b.orderEntryService.SubmitNewOrder(context.Background(), &orderentryapi.NewOrderParams{
				OrderSide: side,
				Quantity:  toApiDec64(bestOpp.Size),
				Price:     toApiDec64(bestOpp.Price),
				Symbol:    b.listing.MarketSymbol,
				ClOrderId: uniqueId.String(),
			})
		}
	}
}

func getQuoteStats(lines []*model.ClobLine, side orderentryapi.Side) (float64,
	float64, float64) {

	qty := &model.Decimal64{}
	bestPrice := &model.Decimal64{}
	worstPrice := &model.Decimal64{}

	zero := &model.Decimal64{}

	for _, line := range lines {
		qty.Add(line.Size)
		if bestPrice.Equal(zero) {
			bestPrice = line.Price
		}
		if worstPrice.Equal(zero) {
			worstPrice = line.Price
		}

		if line.Price.GreaterThan(bestPrice) {
			if side == orderentryapi.Side_BUY {
				bestPrice = line.Price
			} else {
				worstPrice = line.Price
			}
		}

		if line.Price.LessThan(bestPrice) {
			if side == orderentryapi.Side_SELL {
				worstPrice = line.Price
			} else {
				bestPrice = line.Price
			}
		}
	}
	return qty.ToFloat(), bestPrice.ToFloat(), worstPrice.ToFloat()

}

func getBookStats(lines []struct {
	Price     float64 `json:"price"`
	Size      int     `json:"size"`
	Timestamp int64   `json:"timestamp"`
}, side model.Side) (initialQty float64, bestPrice float64, worstPrice float64) {

	for _, line := range lines {
		initialQty += float64(line.Size)
		if bestPrice == 0 {
			bestPrice = line.Price
		}
		if worstPrice == 0 {
			worstPrice = line.Price
		}

		if line.Price > bestPrice {
			if side == model.Side_BUY {
				bestPrice = line.Price
			} else {
				worstPrice = line.Price
			}
		}

		if line.Price < bestPrice {
			if side == model.Side_BUY {
				worstPrice = line.Price
			} else {
				bestPrice = line.Price
			}
		}
	}
	return initialQty, bestPrice, worstPrice
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
			Price:     toApiDec64(model.FasD(bid.Price)),
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

	if totalBidQty.GreaterThan(model.IasD(0)) {
		uniqueId, _ := uuid.NewUUID()
		b.orderEntryService.SubmitNewOrder(context.Background(), &orderentryapi.NewOrderParams{
			OrderSide: orderentryapi.Side_SELL,
			Quantity:  toApiDec64(totalBidQty),
			Price:     toApiDec64(worstBid),
			Symbol:    b.listing.MarketSymbol,
			ClOrderId: uniqueId.String(),
		})

	}

	if totalAskQty.GreaterThan(model.IasD(0)) {
		uniqueId, _ := uuid.NewUUID()
		b.orderEntryService.SubmitNewOrder(context.Background(), &orderentryapi.NewOrderParams{
			OrderSide: orderentryapi.Side_BUY,
			Quantity:  toApiDec64(totalAskQty),
			Price:     toApiDec64(worstAsk),
			Symbol:    b.listing.MarketSymbol,
			ClOrderId: uniqueId.String(),
		})
	}
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
