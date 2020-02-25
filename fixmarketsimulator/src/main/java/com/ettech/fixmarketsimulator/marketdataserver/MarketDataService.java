package com.ettech.fixmarketsimulator.marketdataserver;

import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.marketdataserver.api.FixSimMarketDataServiceGrpc;
import com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver;
import com.google.inject.Inject;
import com.google.protobuf.Empty;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.protobuf.services.ProtoReflectionService;
import io.grpc.stub.StreamObserver;
import org.fixprotocol.components.MarketData;

import java.io.IOException;
import java.util.*;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;

public class MarketDataService {

    private static final Logger logger = Logger.getLogger(MarketDataService.class.getName());

    private Exchange exchange;

    @Inject
    public MarketDataService(Exchange exchange) throws Exception {
        this.exchange = exchange;

    }

    private Server server;

    public void start() throws IOException {
        /* The port on which the server should run */
        int port = 50051;
        server = ServerBuilder.forPort(port)
                .addService(new MarketDataServiceImpl(this.exchange))
                .addService(ProtoReflectionService.newInstance())
                .build()
                .start();
        logger.info("Market data service started, listening on " + port);
        Runtime.getRuntime().addShutdownHook(new Thread() {
            @Override
            public void run() {
                // Use stderr here since the logger may have been reset by its JVM shutdown hook.
                System.err.println("*** shutting down market data service gRPC server since JVM is shutting down");
                try {
                    MarketDataService.this.stop();
                } catch (InterruptedException e) {
                    e.printStackTrace(System.err);
                }
                System.err.println("*** server shut down");
            }
        });
    }

    public void stop() throws InterruptedException {
        if (server != null) {
            server.shutdown().awaitTermination(30, TimeUnit.SECONDS);
        }
    }

    /**
     * Await termination on the main thread since the grpc library uses daemon threads.
     */
    public void blockUntilShutdown() throws InterruptedException {
        if (server != null) {
            server.awaitTermination();
        }
    }

    static class MarketDataServiceImpl extends FixSimMarketDataServiceGrpc.FixSimMarketDataServiceImplBase {


        Map<String, Connection>  partyIdToConnection = new HashMap<>();
        Exchange exchange;

        MarketDataServiceImpl(Exchange exchange) {
            this.exchange = exchange;
        }

        @Override
        public void subscribe(MarketData.MarketDataRequest request, StreamObserver<Empty> responseObserver) {

            logger.info("subscription request:" + request);

            try {
                if (request.getPartiesCount() != 1) {
                    responseObserver.onError(new IllegalArgumentException("request must specify exactly one party"));
                    return;
                }

                var parties = request.getParties(0);
                String partyId = parties.getPartyId();
                synchronized ( partyIdToConnection ) {
                    var connection = partyIdToConnection.get(partyId);
                    if( connection == null ) {
                        responseObserver.onError(new IllegalArgumentException("must connect before attempting to subscribe"));
                        return;
                    }
                    connection.subscribe(request);
                }




            } catch (Exception e) {
                responseObserver.onError(new IllegalArgumentException(e));
            }

            responseObserver.onNext(Empty.newBuilder().build());
            responseObserver.onCompleted();

        }

        @Override
        public void connect(Marketdataserver.Party request, StreamObserver<MarketData.MarketDataIncrementalRefresh> responseObserver) {

            logger.info("market data server connect request received for:" + request);

            synchronized ( partyIdToConnection ) {
                var partyId = request.getPartyId();
                var connection = partyIdToConnection.get(partyId);

                if( connection != null) {
                    try {
                        connection.close();
                    } catch(Exception e) {
                        logger.info("exception when closing connection:" + e);
                    }
                }


                connection = new Connection(exchange, responseObserver, partyId);
                partyIdToConnection.put(partyId, connection);

            }
        }
    }



}



