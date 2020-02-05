package actor

import (
	"fmt"
	"log"
	"os"
)



type SubscriptionClient interface {
	Subscribe(symbol string) error
}

type symbolSource interface {
	fetchSymbol(listingId int, onSymbol chan<- ListingIdSymbol)
}

type subscriptionHandler struct {
	connectionId        string
	listingIdToSymbol   map[int]string
	requestedListingIds map[int]bool
	subscribeChan       chan int
	symbolLookupChan    chan ListingIdSymbol
	symbolSource        symbolSource
	subscriptionClient  SubscriptionClient
	closeChan           chan chan<- bool
	log                 *log.Logger
}


func NewSubscriptionHandler(connectionId string, symbolSource symbolSource, subscriptionClient SubscriptionClient) *subscriptionHandler {
	s := &subscriptionHandler{
		connectionId:        connectionId,
		listingIdToSymbol:   make(map[int]string),
		requestedListingIds: make(map[int]bool),
		subscribeChan:       make(chan int, 10000),
		symbolLookupChan:    make(chan ListingIdSymbol, 10000),
		symbolSource:        symbolSource,
		subscriptionClient:  subscriptionClient,
		closeChan:           make(chan chan<- bool, 1),
		log:                 log.New(os.Stdout, connectionId+"-subscriptionHandler:", log.LstdFlags),
	}

	return s
}

func (s *subscriptionHandler) Close(done chan<- bool) {
	s.closeChan <- done
}

func (s *subscriptionHandler) Subscribe(listingId int) {
	s.subscribeChan <- listingId
}

func (s *subscriptionHandler) Start() {

	go func() {

		for {
			if done := s.readInputChannels(); done != nil {
				done <- true
				return
			}
		}
	}()

}

func (s *subscriptionHandler) readInputChannels() chan<- bool {
	select {
	case listingId := <-s.subscribeChan:
		s.requestedListingIds[listingId] = true
		if _, ok := s.listingIdToSymbol[listingId]; !ok {
			s.symbolSource.fetchSymbol(listingId, s.symbolLookupChan)
		}
	case lts := <-s.symbolLookupChan:
		s.log.Println("subscribing to ", lts.Symbol)
		err := s.subscriptionClient.Subscribe(lts.Symbol)
		if err != nil {
			fmt.Errorf("Failed to Subscribe to %v, error: %w ", lts.Symbol, err)
			return nil
		} else {
			s.listingIdToSymbol[lts.ListingId] = lts.Symbol
			s.log.Println("subscribed to ", lts.Symbol)
		}
	case d := <-s.closeChan:
		return d
	}

	return nil
}
