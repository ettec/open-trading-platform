package com.ettech.fixmarketsimulator.marketdataserver;

import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.marketdataserver.api.FixSimMarketDataServiceGrpc;
import com.google.inject.Inject;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.ServerInterceptors;
import io.grpc.protobuf.services.ProtoReflectionService;
import io.grpc.stub.StreamObserver;
import org.fixprotocol.components.MarketData;

import java.io.IOException;
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

                .addService(ServerInterceptors.intercept(new MarketDataServiceImpl(this.exchange),
                        new MdAuthInterceptor()))
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

        Exchange exchange;

        MarketDataServiceImpl(Exchange exchange) {
            this.exchange = exchange;
        }

        @Override
        public io.grpc.stub.StreamObserver<org.fixprotocol.components.MarketData.MarketDataRequest> connect(
                io.grpc.stub.StreamObserver<org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh> responseObserver) {

            String subscriberId = MdAuthInterceptor.SUBSCRIBER_ID.get();

            logger.info("market data server connect request received for:" + subscriberId);

            var connection = new Connection(exchange, responseObserver, subscriberId);


            StreamObserver<MarketData.MarketDataRequest> so = new StreamObserver<MarketData.MarketDataRequest>() {
                @Override
                public void onNext(MarketData.MarketDataRequest request) {
                    logger.info(subscriberId + " subscription request:" + request);
                    connection.subscribe(request);
                }

                @Override
                public void onError(Throwable t) {
                    connection.close();
                    logger.severe(subscriberId + " connection error:" + t);
                }

                @Override
                public void onCompleted() {
                    connection.close();
                    logger.info(subscriberId + " connection completed");
                }
            };

            return so;
        }
    }


}



