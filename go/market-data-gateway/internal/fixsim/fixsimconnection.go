package fixsim

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"google.golang.org/grpc"
	"log"
	"os"
)

type ListingIdSymbol struct {
	ListingId int
	Symbol    string
}

type fixSimConnection struct {
	address           string
	connectionName    string
	symbolToListingId map[string]int
	idToQuote         map[int]*model.ClobQuote
	refreshChan       chan *marketdata.MarketDataIncrementalRefresh
	mappingChan       chan ListingIdSymbol
	out               chan<- *model.ClobQuote
	fixSimClient      marketDataClient
	symbolLookup      symbolLookup
	dial              dial
	log               *log.Logger
}

type marketDataClient interface {
	connect(connectionId string) (receiveIncRefreshFn, error)
	subscribe(symbol string, subscriberId string) error
	close() error
}

type dial func(target string) (marketDataClient, error)
type symbolLookup func(listingId int) (string, error)

func NewFixSimConnection(
	out chan<- *model.ClobQuote, address string, connectionName string, symbolLookup symbolLookup) *fixSimConnection {

	q := &fixSimConnection{
		address:           address,
		connectionName:    connectionName,
		symbolToListingId: make(map[string]int),
		idToQuote:         make(map[int]*model.ClobQuote),
		refreshChan:       make(chan *marketdata.MarketDataIncrementalRefresh, 10000),
		mappingChan:       make(chan ListingIdSymbol, 1000),
		symbolLookup:      symbolLookup,
		out:               out,
		log:               log.New(os.Stdout, "fixSimConnection:"+connectionName, log.LstdFlags),
	}

	q.dial = func(target string) (client marketDataClient, err error) {

		conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock())
		return &fixSimMarketDataClientImpl{
			client: NewFixSimMarketDataServiceClient(conn),
			conn:   conn,
		}, nil

	}

	return q
}

func (n *fixSimConnection) Connect() error {
	n.log.Println("connecting to ", n.address)

	client, err := n.dial(n.address)
	if err != nil {
		return err
	}
	n.fixSimClient = client

	go func() {
		for {
			if err := n.readInputChannel(); err != nil {
				n.log.Printf("closing read loop: %v", err)
				return
			}
		}
	}()

	go n.connect()

	return nil
}

func (n *fixSimConnection) Subscribe(listingId int) {
	go func() {
		if symbol, err :=  n.symbolLookup(listingId); err == nil {
			n.mappingChan<-ListingIdSymbol{
				ListingId: listingId,
				Symbol:    symbol,
			}
		} else {
			n.log.Printf("error lookingup symbol for listing id: %v, error:%v", listingId, err)
		}
	}()
}

func (n *fixSimConnection) Close() error {
	return n.fixSimClient.close()
}

func (n *fixSimConnection) readInputChannel() error {
	select {
	case m := <-n.mappingChan:
		n.symbolToListingId[m.Symbol] = m.ListingId
		go func() {
			if err := n.fixSimClient.subscribe(m.Symbol, n.connectionName); err != nil {
				n.log.Printf("subscribe call failed:%v", err)
			}
		}()
	case r, ok := <-n.refreshChan:
		if ok {
			for _, incGrp := range r.MdIncGrp {
				symbol := incGrp.GetInstrument().GetSymbol()
				bids := incGrp.MdEntryType == marketdata.MDEntryTypeEnum_MD_ENTRY_TYPE_BID
				if listingId, ok := n.symbolToListingId[symbol]; ok {
					if quote, ok := n.idToQuote[listingId]; ok {
						updatedQuote := updateQuote(quote, incGrp, bids)
						n.idToQuote[listingId] = updatedQuote
						n.out <- updatedQuote
					} else {
						quote := newClobQuote(listingId)
						updatedQuote := updateQuote(quote, incGrp, bids)
						n.idToQuote[listingId] = updatedQuote
						n.out <- updatedQuote
					}
				} else {
					n.log.Printf("received refresh for unknown symbol: %v", symbol)
				}
			}
		} else {
			close(n.out)
			return fmt.Errorf("inbound channel is closed")
		}
	}

	return nil
}

func (n *fixSimConnection) connect() {

	defer func() {
		close(n.refreshChan)
		if err := n.fixSimClient.close(); err != nil {
			n.log.Println("error whilst closing:", err)
		}
	}()

	stream, err := n.fixSimClient.connect(n.connectionName)
	if err != nil {
		n.log.Println("Failed to connect to the market data server:", err)
		return
	}

	for {
		incRefresh, err := stream()
		if err != nil {
			n.log.Println("inbound stream error:", err)
			return
		}

		n.refreshChan <- incRefresh
	}

}

func newClobQuote(listingId int) *model.ClobQuote {
	bids := make([]*model.ClobLine, 0)
	offers := make([]*model.ClobLine, 0)

	return &model.ClobQuote{
		ListingId: int32(listingId),
		Bids:      bids,
		Offers:    offers,
	}
}

func updateQuote(quote *model.ClobQuote, update *marketdata.MDIncGrp, bidSide bool) *model.ClobQuote {

	newQuote := model.ClobQuote{
		ListingId: quote.ListingId,
	}

	if bidSide {
		newQuote.Offers = quote.Offers
		newQuote.Bids = updateClobLines(quote.Bids, update, bidSide)
	} else {
		newQuote.Bids = quote.Bids
		newQuote.Offers = updateClobLines(quote.Offers, update, bidSide)
	}

	return &newQuote
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
