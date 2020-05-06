package com.ettech.fixmarketsimulator.marketdataserver;

import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.exchange.OrderBook;
import io.grpc.stub.StreamObserver;
import org.fixprotocol.components.InstrmtMDReqGrp;
import org.fixprotocol.components.MarketData;
import org.slf4j.LoggerFactory;

import java.util.*;

public class Connection {
    org.slf4j.Logger log;

    Map<String,MarketDataSubscription> subscriptions = new HashMap<String,MarketDataSubscription>();


    Exchange exchange;
    StreamObserver<MarketData.MarketDataIncrementalRefresh> responseObserver;

    Connection(Exchange exchange, StreamObserver<MarketData.MarketDataIncrementalRefresh> responseObserver, String subscriberId) {
        log = LoggerFactory.getLogger(Connection.class.getSimpleName() + ":" + subscriberId);
        this.responseObserver = responseObserver;
        this.exchange = exchange;
    }


    void close() {
        subscriptions.values().forEach(s->s.close());
    }

    void subscribe(MarketData.MarketDataRequest msg) {

        log.info("Received subscription request", msg);

        for( InstrmtMDReqGrp mdReqGrp : msg.getInstrmtMdReqGrpList() ) {
            String symbol = mdReqGrp.getInstrument().getSymbol();

            log.info("Received subscribe request for symbol {}", symbol);

            if( subscriptions.containsKey(symbol)) {
                var subscription = subscriptions.remove(symbol);
                subscription.close();
                return;
            }

            OrderBook book = exchange.getOrderBook(symbol);
            subscriptions.put(symbol, new MarketDataSubscription(this, book, msg.getMdReqId()));
        }

    }

    public void send(MarketData.MarketDataIncrementalRefresh refresh) {
        try {
            responseObserver.onNext(refresh);


        } catch (Throwable t) {
            log.error("failed to send refresh, closing connection", t);
            this.close();
            responseObserver.onError(t);
        }

    }

}
