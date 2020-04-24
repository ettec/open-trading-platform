package marketdata

import (
	"fmt"
	"github.com/emicklei/go-restful/log"
	"github.com/ettec/open-trading-platform/go/common/api/marketdatasource"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type marketDataSourceServer struct {
	quoteDistributor QuoteDistributor
}

func NewMarketDataSource(quoteDistributor QuoteDistributor) marketdatasource.MarketDataSourceServer {
	return &marketDataSourceServer{quoteDistributor}
}

var maxSubscriptions = 10000

const SubscriberIdKey = "subscriber_id"

func (s *marketDataSourceServer) Connect(stream marketdatasource.MarketDataSource_ConnectServer) error {

	ctx, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return fmt.Errorf("failed to retrieve call context")
	}

	values := ctx.Get(SubscriberIdKey)
	if len(values) != 1 {
		return fmt.Errorf("must specify string value for %v", SubscriberIdKey)
	}

	fromClientId := values[0]
	subscriberId := fromClientId + ":" + uuid.New().String()

	log.Printf("connect request received for subscriber %v, unique connection id: %v ", fromClientId, subscriberId)


	cc := NewConflatedQuoteConnection(subscriberId,  s.quoteDistributor.GetNewQuoteStream(), maxSubscriptions)
	defer cc.Close()

	go func() {
		for {
			subscription, err := stream.Recv()

			if err != nil {
				log.Printf("subscriber:%v inbound stream error:%v ", subscriberId, err)
				break
			} else {
				log.Printf("subscribe request, subscriber id:%v, listing id:%v", subscriberId, subscription.ListingId)
				cc.Subscribe(subscription.ListingId)
			}
		}
	}()

	for mdUpdate := range cc.GetStream() {
		if err := stream.Send(mdUpdate); err != nil {
			log.Printf("error on connection for subscriber %v, closing connection, error:%v ", subscriberId, err)
			break
		}
	}

	return nil
}
