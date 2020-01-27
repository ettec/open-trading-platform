package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/api"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type connection struct {
	quoteChan     chan *marketdata.MarketDataSnapshotFullRefresh
	stream        api.MarketDataService_ConnectServer
	subscriptions map[int]bool
}

type service struct {
	partyIdToConnection map[string]connection
}

func newService() *service {
	return &service{partyIdToConnection: make(map[string]connection)}
}

func (*service) Subscribe(context.Context, *marketdata.MarketDataRequest) (*empty.Empty, error) {
	return nil, nil
}

func (s *service) Connect(request *api.ConnectRequest, stream api.MarketDataService_ConnectServer) error {

	/*
		partyId := request.GetPartyId()

		log
			con, ok = s.partyIdToConnection[partyId]
		if ok {
			return fmt.Errorf("Connection for part")
		}*/

	return nil
}

/*
func (*server) AddSubscription(context context.Context, subscription *model.Subscription) (*model.AddSubscriptionResponse, error) {

	return &model.AddSubscriptionResponse{
		Message: fmt.Sprintf("Subscription for %v added for subscriber %v", subscription.ListingId, subscription.SubscriberId),
	}, nil
}

func getMetaData(ctx context.Context) (username string, appInstanceId string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)


	if !ok {
		return "", "", fmt.Errorf("failed to read metadata from the context")
	}

	appInstanceIds := md.Get("app-instance-id")
	if len(appInstanceIds) != 1 {
		return "", "", fmt.Errorf("unable to retrieve app-instance-id from metadata")
	}
	appInstanceId = appInstanceIds[0]

	usernames := md.Get("user-name")
	if len(usernames) != 1 {
		return "", "", fmt.Errorf("unable to retrieve user-name from metadata")
	}
	username = usernames[0]

	return username, appInstanceId, nil
}

func (*server) Subscribe(request *model.SubscribeRequest, stream model.MarketDataService_SubscribeServer) error {


		username, appInstanceId, err := getMetaData(stream.Context())
		if err != nil {
			return err
		}



	marketDataChannel := make(chan *model.Quote, 3)

	go func() {

		i := 0
		for {
			i++
			time.Sleep(3 * time.Second)

			if i%2 == 0 {
				marketDataChannel <- &model.Quote{
					ListingId: 54123,
					Depth: []*model.DepthLine{{BidSize: &model.Decimal64{Mantissa: 10}, BidPrice: &model.Decimal64{Mantissa: 13},
						AskPrice: &model.Decimal64{Mantissa: 14}, AskSize: &model.Decimal64{Mantissa: 5}},
						{BidSize: &model.Decimal64{Mantissa: 9}, BidPrice: &model.Decimal64{Mantissa: 12},
							AskPrice: &model.Decimal64{Mantissa: 15}, AskSize: &model.Decimal64{Mantissa: 6}},
					},
				}
			} else {
				marketDataChannel <- &model.Quote{
					ListingId: 54123,
					Depth: []*model.DepthLine{
						{BidSize: &model.Decimal64{Mantissa: 2}, BidPrice: &model.Decimal64{Mantissa: 12},
							AskPrice: &model.Decimal64{Mantissa: 15}, AskSize: &model.Decimal64{Mantissa: 3}},
					},
				}
			}

		}

		defer close(marketDataChannel)
	}()

	for mdUpdate := range marketDataChannel {
		stream.Send(mdUpdate)
	}

	return nil
} */

func main() {

	port := "50551"
	fmt.Println("Starting Client Market Data Gateway on port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()
	api.RegisterMarketDataServiceServer(s, &service{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
