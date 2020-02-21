package actor

import (
	"log"
	"os"
)


type SubscriptionHandler interface {
	Actor
	Subscribe(listingId int)
}

type subscribeFn = func(symbol string)
type subscribeToListing = func(listingId int)

type SubscriptionClient interface {
	Subscribe(symbol string)
}


type ListingIdSymbol struct {
	ListingId int
	Symbol    string
}

type fetchSymbolForListingFn = func(listingId int, onSymbol chan<- ListingIdSymbol)

type subscriptionHandler struct {
	actorImpl
	connectionId        string
	listingIdToSymbol   map[int]string
	requestedListingIds map[int]bool
	subscribeChan       chan int
	symbolLookupChan    chan ListingIdSymbol
	fetchSymbolForListingFn        fetchSymbolForListingFn
	subscribeFn         subscribeFn
	log                 *log.Logger
}


func NewSubscriptionHandler(connectionId string, fetchSymbolForListingFn fetchSymbolForListingFn, subscribeFn subscribeFn) *subscriptionHandler {
	s := &subscriptionHandler{
		connectionId:        connectionId,
		listingIdToSymbol:   make(map[int]string),
		requestedListingIds: make(map[int]bool),
		subscribeChan:       make(chan int, 10000),
		symbolLookupChan:    make(chan ListingIdSymbol, 10000),
		fetchSymbolForListingFn:        fetchSymbolForListingFn,
		subscribeFn:         subscribeFn,
		log:                 log.New(os.Stdout, connectionId+"-subscriptionHandler:", log.LstdFlags),
	}

	s.actorImpl = newActorImpl("subscriptionHandler for connection " + connectionId, s.readInputChannels)

	return s
}


func (s *subscriptionHandler) Subscribe(listingId int) {
	s.subscribeChan <- listingId
}


func (s *subscriptionHandler) readInputChannels() (chan<- bool, error) {
	select {
	case listingId := <-s.subscribeChan:
		s.requestedListingIds[listingId] = true
		if _, ok := s.listingIdToSymbol[listingId]; !ok {
			s.fetchSymbolForListingFn(listingId, s.symbolLookupChan)
		}
	case lts := <-s.symbolLookupChan:
		s.log.Println("subscribing to ", lts.Symbol)
		s.subscribeFn(lts.Symbol)
		s.listingIdToSymbol[lts.ListingId] = lts.Symbol
	case d := <-s.closeChan:
		return d, nil
	}

	return nil, nil
}
