package common

import (
	"context"
	services "github.com/ettec/open-trading-platform/go/common/services"
	"github.com/ettec/open-trading-platform/go/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"os"
	"time"
)

type SubscriptionHandler interface {
	Subscribe(listingId int32)
}

type SubscriptionClient interface {
	Subscribe(symbol string)
}


type GetListingFn = func(listingId int32, onSymbol chan<- *model.Listing)

type listingSource struct {
	fetchReqChan chan fetchRequest
	log          *log.Logger
	errLog       *log.Logger
}

func NewListingSource(targetAddress string) (*listingSource, error) {
	s := &listingSource{
		fetchReqChan: make(chan fetchRequest, 10000),
		log:          log.New(os.Stdout, "", log.Ltime|log.Lshortfile),
		errLog:       log.New(os.Stdout, "", log.Ltime|log.Lshortfile),
	}

	conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(120*time.Second))
	if err != nil {
		return nil, err
	}

	sdc := services.NewStaticDataServiceClient(conn)

	go func() {

		for {
			state := conn.GetState()
			for state != connectivity.Ready {
				s.log.Printf("waiting for static data service connection to be ready....")
				conn.WaitForStateChange(context.Background(), state)
				state = conn.GetState()
				s.log.Println("static data service connection state is:", state)
			}

			select {
			case fr := <-s.fetchReqChan:
				listing, err := sdc.GetListing(context.Background(), &services.ListingId{
					ListingId: fr.listingId,
				})

				if err != nil {
					s.errLog.Printf("error retrieving listing:%v", err)
					s.fetchReqChan <- fr
					break
				}

				s.log.Println("received listing:", listing )

				fr.resultChan <- listing
			}

		}
	}()

	return s, nil
}

type fetchRequest struct {
	listingId  int32
	resultChan chan<- *model.Listing
}

func (s *listingSource) GetListing(listingId int32, result chan<- *model.Listing) {
	s.fetchReqChan <- fetchRequest{listingId: listingId, resultChan: result}
}
