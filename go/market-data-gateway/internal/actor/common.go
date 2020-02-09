package actor

import (
	"log"
)


type Actor interface {
	Start()
	Close(chan<- bool)
}

type actorImpl struct {
	id        string
	process   func() (chan<- bool, error)
	closeChan chan chan<- bool
}

func newActorImpl(id string, process func() (chan<- bool, error)) actorImpl {
	return actorImpl{id: id, process: process}
}

func (a *actorImpl) Start() {

	if a.closeChan != nil {
		log.Panic("actor has already been started:", a.id)
	}

	a.closeChan = make(chan chan<- bool, 1)

	go func() {
		for {
			if d, err := a.process(); d != nil {
				log.Println("closing ", a.id)
				select {
				case d <- true:
				default:
					log.Printf("%v unable to send signal on done channel, exiting anyway", a.id)
				}

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
