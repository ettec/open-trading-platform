package marketdata

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/common/api/marketdatasource"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	logger "log"
	"os"
)

type marketDataSourceServer struct {
	quoteDistributor QuoteDistributor
}

func NewMarketDataSource(quoteDistributor QuoteDistributor ) marketdatasource.MarketDataSourceServer {
	return &marketDataSourceServer{quoteDistributor}
}

var maxSubscriptions = 10000

var log = logger.New(os.Stdout, "", logger.Ltime|logger.Lshortfile)
var errLog = logger.New(os.Stderr, "", logger.Ltime|logger.Lshortfile)

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

	out := make(chan *model.ClobQuote, 100)
	cc := NewConflatedQuoteConnection(subscriberId, out, s.quoteDistributor, maxSubscriptions)
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

	for mdUpdate := range out {
		if err := stream.Send(mdUpdate); err != nil {
			log.Printf("error on connection for subscriber %v, closing connection, error:%v ", subscriberId, err)
			break
		}
	}

	return nil
}