package fixsim

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data/market-data-gateway-fixsim/internal/fix/marketdata"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/staticdata"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
)

type fixSimAdapter struct {
	connectionName    string
	symbolToListingId map[string]int32
	idToQuote         map[int32]*model.ClobQuote
	refreshInChan     chan *marketdata.MarketDataIncrementalRefresh
	listingInChan     chan *model.Listing
	out               chan *model.ClobQuote
	fixSimClient      MarketDataClient
	getListing        staticdata.GetListingFn
	closeChan         chan bool
	log               *log.Logger
	errLog            *log.Logger
}

func (n *fixSimAdapter) GetStream() <-chan *model.ClobQuote {
	return n.out
}

func (n *fixSimAdapter) Close() {
	n.closeChan <- true
}

type MarketDataClient interface {
	subscribe(symbol string) error
	close() error
}

type newMarketDataClient = func(id string, out chan<- *marketdata.MarketDataIncrementalRefresh) (MarketDataClient, error)

func NewFixSimAdapter(
	newClientFn newMarketDataClient, connectionName string, symbolLookup staticdata.GetListingFn,
	sendBufferSize int) (*fixSimAdapter, error) {

	n := &fixSimAdapter{
		out:               make(chan *model.ClobQuote, sendBufferSize),
		connectionName:    connectionName,
		symbolToListingId: make(map[string]int32),
		idToQuote:         make(map[int32]*model.ClobQuote),
		refreshInChan:     make(chan *marketdata.MarketDataIncrementalRefresh, 10000),
		listingInChan:     make(chan *model.Listing, 1000),
		getListing:        symbolLookup,
		closeChan:         make(chan bool),
		log:               log.New(log.Writer(), connectionName+" ", log.Flags()),
		errLog:            log.New(os.Stderr, connectionName+" ", log.Flags()),
	}

	var err error
	n.fixSimClient, err = newClientFn(connectionName, n.refreshInChan)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-n.closeChan:
				log.Print("closed fix sim adapter")
				break
			case l := <-n.listingInChan:
				n.symbolToListingId[l.MarketSymbol] = l.Id
				go func() {
					if err := n.fixSimClient.subscribe(l.MarketSymbol); err != nil {
						n.errLog.Printf("subscribe call failed:%v", err)
					}
				}()
			case r := <-n.refreshInChan:

				if r != nil {
					for _, incGrp := range r.MdIncGrp {
						symbol := incGrp.GetInstrument().GetSymbol()
						if listingId, ok := n.symbolToListingId[symbol]; ok {

							var originalQuote *model.ClobQuote
							if originalQuote, ok = n.idToQuote[listingId]; !ok {
								originalQuote = newClobQuote(listingId)
							}

							newQuote, err := copyQuote(originalQuote)
							newQuote.StreamInterrupted = false
							newQuote.StreamStatusMsg = ""
							if err != nil {
								n.errLog.Print("failed to copy originalQuote:", err)
							}

							linesUpdate := incGrp.MdEntryType == marketdata.MDEntryTypeEnum_MD_ENTRY_TYPE_BID ||
								incGrp.MdEntryType == marketdata.MDEntryTypeEnum_MD_ENTRY_TYPE_OFFER
							if linesUpdate {

								bids := incGrp.MdEntryType == marketdata.MDEntryTypeEnum_MD_ENTRY_TYPE_BID
								updateQuoteDepth(originalQuote, newQuote, incGrp, bids)

							} else if incGrp.MdEntryType == marketdata.MDEntryTypeEnum_MD_ENTRY_TYPE_TRADE {
								newQuote.LastPrice = &model.Decimal64{Mantissa:  incGrp.MdEntryPx.Mantissa, Exponent:  incGrp.MdEntryPx.Exponent}
								newQuote.LastQuantity = &model.Decimal64{Mantissa:  incGrp.MdEntrySize.Mantissa, Exponent:  incGrp.MdEntrySize.Exponent}

							} else if incGrp.MdEntryType == marketdata.MDEntryTypeEnum_MD_ENTRY_TYPE_TRADE_VOLUME {
								newQuote.TradedVolume = &model.Decimal64{Mantissa:  incGrp.MdEntrySize.Mantissa, Exponent:  incGrp.MdEntrySize.Exponent}
							} else {
								continue
							}

							n.idToQuote[listingId] = newQuote
							n.out <- newQuote

						} else {
							n.errLog.Printf("received refresh for unknown symbol: %v", symbol)

						}

					}
				} else {
					for id := range n.idToQuote {
						emptyQuote := newClobQuote(id)
						emptyQuote.StreamInterrupted = true
						emptyQuote.StreamStatusMsg = "fix sim adapter stream interrupted"
						n.idToQuote[id] = emptyQuote
						n.out <- emptyQuote
					}
				}
			}
		}
	}()

	return n, nil
}

func copyQuote(quote *model.ClobQuote) (*model.ClobQuote, error) {

	bytes, err := proto.Marshal(quote)
	if err != nil {
		return nil, fmt.Errorf("failed to copy quote:%w", err)
	}

	quoteCopy := &model.ClobQuote{}
	err = proto.Unmarshal(bytes, quoteCopy)
	if err != nil {
		return nil, fmt.Errorf("failed to copy quote:%w", err)
	}

	return quoteCopy, nil
}

func (n *fixSimAdapter) Subscribe(listingId int32) {
	n.getListing(listingId, n.listingInChan)
}

func newClobQuote(listingId int32) *model.ClobQuote {
	bids := make([]*model.ClobLine, 0)
	offers := make([]*model.ClobLine, 0)

	return &model.ClobQuote{
		ListingId: listingId,
		Bids:      bids,
		Offers:    offers,
	}
}

func updateQuoteDepth(originalQuote *model.ClobQuote, newQuote *model.ClobQuote, update *marketdata.MDIncGrp, bidSide bool) {

	if bidSide {
		newQuote.Bids = updateClobLines(originalQuote.Bids, update, bidSide)
	} else {
		newQuote.Offers = updateClobLines(originalQuote.Offers, update, bidSide)
	}

}

func updateClobLines(lines []*model.ClobLine, update *marketdata.MDIncGrp, bids bool) []*model.ClobLine {

	updateAction := update.GetMdUpdateAction()
	newClobLines := make([]*model.ClobLine, 0, len(lines)+1)

	compareTest := 1
	if bids {
		compareTest = -1
	}

	switch updateAction {
	case marketdata.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW:
		inserted := false

		newLine := &model.ClobLine{
			Size:    &model.Decimal64{Mantissa: update.MdEntrySize.Mantissa, Exponent: update.MdEntrySize.Exponent},
			Price:   &model.Decimal64{Mantissa: update.MdEntryPx.Mantissa, Exponent: update.MdEntryPx.Exponent},
			EntryId: update.MdEntryId,
		}

		for _, line := range lines {
			compareResult := model.Compare(*line.Price, model.Decimal64(*update.GetMdEntryPx()))
			if !inserted && compareResult == compareTest {
				newClobLines = append(newClobLines, newLine)
				inserted = true
			}
			newClobLines = append(newClobLines, line)
		}

		if !inserted {
			newClobLines = append(newClobLines, newLine)
		}

	case marketdata.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE:
		inserted := false

		newLine := &model.ClobLine{
			Size:    &model.Decimal64{Mantissa: update.MdEntrySize.Mantissa, Exponent: update.MdEntrySize.Exponent},
			Price:   &model.Decimal64{Mantissa: update.MdEntryPx.Mantissa, Exponent: update.MdEntryPx.Exponent},
			EntryId: update.MdEntryId,
		}

		for _, line := range lines {
			compareResult := model.Compare(*line.Price, model.Decimal64(*update.GetMdEntryPx()))
			if !inserted && compareResult == compareTest {
				newClobLines = append(newClobLines, newLine)
				inserted = true
			}
			if line.EntryId != newLine.EntryId {
				newClobLines = append(newClobLines, line)
			}

		}

		if !inserted {
			newClobLines = append(newClobLines, newLine)
		}

	case marketdata.MDUpdateActionEnum_MD_UPDATE_ACTION_DELETE:
		for _, line := range lines {
			if line.EntryId != update.MdEntryId {
				newClobLines = append(newClobLines, line)
			}
		}
	}

	return newClobLines

}
