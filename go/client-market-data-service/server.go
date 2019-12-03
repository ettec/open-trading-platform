package main

import (
	"context"
	"fmt"
	"github.com/coronationstreet/open-trading-platform/client-market-data-service/cmds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

type server struct{}

func (*server) AddSubscription(context context.Context, subscription *cmds.Subscription) (*cmds.AddSubscriptionResponse, error) {


	return &cmds.AddSubscriptionResponse{
		Message: "Subscription for " + subscription.ListingId + " added for subscriber " + subscription.SubscriberId,
	}, nil
}

func getMetaData(ctx context.Context) (username string, appInstanceId string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", fmt.Errorf("failed to read metadata from the context")
	}

	appInstanceIds := md.Get( "app-instance-id")
	if len(appInstanceIds) != 1 {
		return "", "", fmt.Errorf("unable to retrieve app-instance-id from metadata")
	}
	appInstanceId = appInstanceIds[0]


	usernames := md.Get( "user-name")
	if len(usernames) != 1 {
		return "", "", fmt.Errorf("unable to retrieve user-name from metadata")
	}
	username = usernames[0]

	return username, appInstanceId, nil
}

func (*server) Subscribe(request *cmds.SubscribeRequest, stream cmds.ClientMarketDataService_SubscribeServer) error {


	/*
	username, appInstanceId, err := getMetaData(stream.Context())
	if err != nil {
		return err
	}*/



	marketDataChannel := make(chan *cmds.Book, 3)

	go func() {

		i := 0
		for {
			i++
			time.Sleep(1 * time.Second)

			if i % 2 == 0 {
				marketDataChannel <- &cmds.Book{
					ListingId: "Blah",
					Depth: []*cmds.BookLine{{BidSize:  "10",BidPrice: "13",	AskPrice: "14", AskSize:  "5",},
											{BidSize:  "9",BidPrice: "12",	AskPrice: "15", AskSize:  "6",},
					},
				}
			} else {
				marketDataChannel <- &cmds.Book{
					ListingId: "Blah",
					Depth: []*cmds.BookLine{
											{BidSize:  "2",BidPrice: "12",	AskPrice: "15", AskSize:  "3",},
					},
				}
			}

		}

		defer close(marketDataChannel)
	}()

	for mdUpdate := range marketDataChannel {
		stream.Send(mdUpdate)
	}

	return nil;
}




func main() {

	port := "50551"
	fmt.Println("Starting Client Market Data Server on the port:" + port)
	lis, err := net.Listen("tcp", "0.0.0.0:" + port)

	if err != nil {
		log.Fatalf("Error while listening : %v", err)
	}

	s := grpc.NewServer()
	cmds.RegisterClientMarketDataServiceServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error while serving : %v", err)
	}

}
