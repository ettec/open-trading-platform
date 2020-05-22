package source

import (
	"fmt"
	"github.com/emicklei/go-restful/log"
	"github.com/ettec/open-trading-platform/go/common/api/marketdatasource"
	"github.com/ettec/open-trading-platform/go/common/marketdata"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc/metadata"
)

var connections = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "active_connections",
	Help: "The number of active connections",
})

var quotesSent = promauto.NewCounter(prometheus.CounterOpts{
	Name: "quotes_sent",
	Help: "The number of quotes sent across all clients",
})

type marketDataSourceServer struct {
	quoteDistributor marketdata.QuoteDistributor
}

func NewMarketDataSource(quoteDistributor marketdata.QuoteDistributor) marketdatasource.MarketDataSourceServer {
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

	out := make(chan *model.ClobQuote)

	cc := marketdata.NewConflatedQuoteConnection(subscriberId, s.quoteDistributor.GetNewQuoteStream(), out, maxSubscriptions)
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

	connections.Inc()

	for mdUpdate := range out {
		if err := stream.Send(mdUpdate); err != nil {
			log.Printf("error on connection for subscriber %v, closing connection, error:%v ", subscriberId, err)
			break
		}

		quotesSent.Inc()
	}

	connections.Dec()

	return nil
}
