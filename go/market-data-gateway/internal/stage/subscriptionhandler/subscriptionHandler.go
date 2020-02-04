package subscriptionhandler

import (
	"context"
	"errors"
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/stage"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"os"
)

type symbolSource interface {
	fetchSymbol(listingId int, onSymbol chan<- stage.ListingIdSymbol)
}

type subscriptionClient interface {
	Subscribe(ctx context.Context, in *marketdata.MarketDataRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type subscriptionHandler struct {
	connectionId      string
	listingIdToSymbol map[int]string
	subscribeChan     chan int
	symbolLookupChan  chan stage.ListingIdSymbol
	symbolSource      symbolSource
	simClient         subscriptionClient
	closeChan         chan bool
	log               *log.Logger
}

var closed = errors.New("closed")

func newSubscriptionHandler(connectionId string, simClient subscriptionClient, symbolSource symbolSource) *subscriptionHandler {

	s := &subscriptionHandler{
		connectionId:      connectionId,
		listingIdToSymbol: make(map[int]string),
		subscribeChan:     make(chan int, 10000),
		symbolLookupChan:  make(chan stage.ListingIdSymbol, 10000),
		symbolSource:      symbolSource,
		simClient:         simClient,
		closeChan:         make(chan bool, 1),
		log:               log.New(os.Stdout, connectionId+"-subscriptionHandler:", log.LstdFlags),
	}

	return s
}

func (s *subscriptionHandler) close() {
	s.closeChan <- true
}

func (s *subscriptionHandler) subscribe(listingId int) {
	s.subscribeChan <- listingId
}

func (s *subscriptionHandler) start() {

	go func() {
		log.Println("Connecting to market data server at ", s.simClient)

		for {
			if err := s.readInputChannels(); err != nil {
				if err != closed {
					log.Printf("exiting subscription handler read loop due to error: %v", err)
				} else {
					log.Printf("subscription handler closed")
				}

			}

		}
	}()

}

func (s *subscriptionHandler) readInputChannels() error {
	select {
	case l := <-s.subscribeChan:
		if _, ok := s.listingIdToSymbol[l]; !ok {
			s.symbolSource.fetchSymbol(l, s.symbolLookupChan)
		}
	case ls := <-s.symbolLookupChan:
		s.log.Println("subscribing to ", ls.Symbol)
		request := &marketdata.MarketDataRequest{Parties: []*common.Parties{{PartyId: s.connectionId}},
			InstrmtMdReqGrp: []*common.InstrmtMDReqGrp{{Instrument: &common.Instrument{Symbol: ls.Symbol}}}}
		_, err := s.simClient.Subscribe(context.Background(), request)
		if err != nil {
			fmt.Errorf("Failed to subscribe to %v, error: %w ", ls.Symbol, err)
			return err
		} else {
			s.listingIdToSymbol[ls.ListingId] = ls.Symbol
			s.log.Println("subscribed to ", ls.Symbol)
		}
	case <-s.closeChan:
		return closed
	}
	return nil
}
