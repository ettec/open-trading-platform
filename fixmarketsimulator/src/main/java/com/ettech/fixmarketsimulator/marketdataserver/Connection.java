package com.ettech.fixmarketsimulator.marketdataserver;

import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.exchange.OrderBook;
import io.grpc.stub.StreamObserver;
import org.fixprotocol.components.InstrmtMDReqGrp;
import org.fixprotocol.components.MarketData;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

public class Connection {
    org.slf4j.Logger log;

    List<MarketDataSubscription> subscriptions = new ArrayList<>();
    Set<String> symbols = new HashSet<>();

    Exchange exchange;
    StreamObserver<MarketData.MarketDataIncrementalRefresh> responseObserver;

    Connection(Exchange exchange, StreamObserver<MarketData.MarketDataIncrementalRefresh> responseObserver, String subscriberId) {
        log = LoggerFactory.getLogger(Connection.class.getSimpleName() + ":" + subscriberId);
        this.responseObserver = responseObserver;
        this.exchange = exchange;
    }

    void setResponseObserver(StreamObserver<MarketData.MarketDataIncrementalRefresh> responseObserver) {
        this.responseObserver = responseObserver;
    }

    void close() {
        this.responseObserver.onCompleted();
    }

    void subscribe(MarketData.MarketDataRequest msg) {

        log.info("Received subscription request", msg);

        for( InstrmtMDReqGrp mdReqGrp : msg.getInstrmtMdReqGrpList() ) {
            String symbol = mdReqGrp.getInstrument().getSymbol();

            log.info("Received subscribe request for symbol {}", symbol);

            if( symbols.contains(symbol)) {
                log.warn("Ignoring subscribe request as already subscribed to " + symbol);
                return;
            }

            symbols.add(symbol);

            OrderBook book = exchange.getOrderBook(symbol);
            subscriptions.add(new MarketDataSubscription(this, book, msg.getMdReqId()));
        }

    }

    public void send(MarketData.MarketDataIncrementalRefresh refresh) {
        responseObserver.onNext(refresh);
    }

}
