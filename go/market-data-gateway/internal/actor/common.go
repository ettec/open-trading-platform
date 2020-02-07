package actor

import (
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"log"
)



type ListingIdSymbol struct {
	ListingId int
	Symbol    string
}

type ClobQuoteSink interface {
	Send(quote *model.ClobQuote)
}

type Actor interface {
	Start()
	Close(chan<- bool)
}


type actorImpl struct {
	id string
	process  func() (chan<- bool, error)
	closeChan chan chan<- bool

}

func newActorImpl(id string, process  func() (chan<- bool, error)) actorImpl {
	return actorImpl{id:id, process:process}
}


func (a *actorImpl) Start()  {

	if a.closeChan != nil {
		log.Panic("actor has already been started:", a.id)
	}

	a.closeChan = make(chan chan<-bool,1)

	go func() {
		for {
			if d, err := a.process(); d != nil {
				log.Println("closing ", a.id)
				d<-true
				return
			} else if err != nil {
				log.Printf("closing %v due to error %v", a.id, err)
				return
			}
		}
	}()
}

func (a *actorImpl) Close(d chan<- bool) {
	a.closeChan <- d
}
