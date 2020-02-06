package actor

import (
	"log"
	"os"
)


type SubscriptionHandler interface {
	Actor
	Subscribe(listingId int)
}

type SubscriptionClient interface {
	Subscribe(symbol string)
}

type symbolSource interface {
	FetchSymbol(listingId int, onSymbol chan<- ListingIdSymbol)
}

type subscriptionHandler struct {
	actorImpl
	connectionId        string
	listingIdToSymbol   map[int]string
	requestedListingIds map[int]bool
	subscribeChan       chan int
	symbolLookupChan    chan ListingIdSymbol
	symbolSource        symbolSource
	subscriptionClient  SubscriptionClient
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
			s.symbolSource.FetchSymbol(listingId, s.symbolLookupChan)
		}
	case lts := <-s.symbolLookupChan:
		s.log.Println("subscribing to ", lts.Symbol)
		s.subscriptionClient.Subscribe(lts.Symbol)
		s.listingIdToSymbol[lts.ListingId] = lts.Symbol
	case d := <-s.closeChan:
		return d, nil
	}

	return nil, nil
}
