package common

import (
	"context"
	"fmt"
	services "github.com/ettec/open-trading-platform/go/common/api/staticdataservice"
	"github.com/ettec/open-trading-platform/go/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
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

type ListingSource interface {
	GetListing(listingId int32, result chan<- *model.Listing)
}

type listingSource struct {
	sdcTaskChan chan staticDataServiceTask
	log         *log.Logger
	errLog      *log.Logger
}

type GrpcConnection interface {
	GetState() connectivity.State
	WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool
}

type GetStaticDataServiceClientFn = func() (services.StaticDataServiceClient, GrpcConnection, error)

func NewStaticDataSource(targetAddress string) (*listingSource, error) {
	return newStaticDataSource(func() (client services.StaticDataServiceClient, connection GrpcConnection, err error) {
		log.Println("connecting to static data service at:" + targetAddress)
		conn, err := grpc.Dial(targetAddress, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(120*time.Second))
		if err != nil {
			return nil, nil, err
		}

		sdc := services.NewStaticDataServiceClient(conn)

		return sdc, conn, nil
	})
}

func newStaticDataSource(getConnection GetStaticDataServiceClientFn) (*listingSource, error) {
	s := &listingSource{
		sdcTaskChan: make(chan staticDataServiceTask, 10000),
		log:         log.New(os.Stdout, "", log.Ltime|log.Lshortfile),
		errLog:      log.New(os.Stdout, "", log.Ltime|log.Lshortfile),
	}

	sdc, conn, err := getConnection()
	if err != nil {
		return nil, err
	}

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
			case t := <-s.sdcTaskChan:
				err := t(sdc)
				if err != nil {
					s.sdcTaskChan <- t
					s.errLog.Printf("error executing static data service task, retry schduled.  Error:%v", err)
				}
			}

		}
	}()

	return s, nil
}

type staticDataServiceTask func(sdc services.StaticDataServiceClient) error

func (s *listingSource) GetListing(listingId int32, result chan<- *model.Listing) {
	s.sdcTaskChan <- func(sdc services.StaticDataServiceClient) error {
		listing, err := sdc.GetListing(context.Background(), &services.ListingId{
			ListingId: listingId,
		})

		if err != nil {
			st, ok := status.FromError(err)
			if !ok || st.Code() != codes.NotFound {
				return fmt.Errorf("error retrieving listing:%v", err)
			} else {
				s.errLog.Printf("no listing found for id:%v", listingId)
			}
		} else {
			s.log.Println("received listing:", listing)
			result <- listing
		}

		return nil
	}
}

func (s *listingSource) GetListingMatching(matchParams *services.ExactMatchParameters, result chan<- *model.Listing) {
	s.sdcTaskChan <- func(sdc services.StaticDataServiceClient) error {
		listing, err := sdc.GetListingMatching(context.Background(), matchParams)

		if err != nil {
			st, ok := status.FromError(err)
			if !ok || st.Code() != codes.NotFound {
				return fmt.Errorf("error retrieving listing:%v", err)
			} else {
				s.errLog.Printf("no listing found for match params:%v", matchParams)
			}
		} else {
			s.log.Printf("received listing:%v for symbol matching:%v and mic:%v", listing, matchParams.Symbol,
				matchParams.Mic)
			result <- listing
		}

		return nil
	}
}

func (s *listingSource) GetListingsWithSameInstrument(listingId int32, listingGroupsIn chan<- []*model.Listing) {

	s.sdcTaskChan <- func(sdc services.StaticDataServiceClient) error {

		listings, err := sdc.GetListingsWithSameInstrument(context.Background(), &services.ListingId{
			ListingId: listingId,
		})

		if err != nil {
			st, ok := status.FromError(err)
			if !ok || st.Code() != codes.NotFound {
				return fmt.Errorf("error retrieving listings :%v", err)
			} else {
				s.errLog.Printf("no listings found for same instrument, listing id:%v", listingId)
			}
		} else {
			s.log.Println("received listings for same instrument:", listings)
			listingGroupsIn <- listings.Listings
		}

		return nil
	}

}
