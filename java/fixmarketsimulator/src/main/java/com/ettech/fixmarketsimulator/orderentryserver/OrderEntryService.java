package com.ettech.fixmarketsimulator.orderentryserver;

import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.exchange.Side;
import com.ettech.fixmarketsimulator.marketdataserver.MarketDataService;
import com.ettech.fixmarketsimulator.orderentryserver.api.OrderEntryServiceGrpc;
import com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi;
import com.google.inject.Inject;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.protobuf.services.ProtoReflectionService;
import io.grpc.stub.StreamObserver;

import java.io.IOException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;

public class OrderEntryService {

    private static final Logger logger = Logger.getLogger(MarketDataService.class.getName());

    private Exchange exchange;

    @Inject
    public OrderEntryService(Exchange exchange) throws Exception {
        this.exchange = exchange;

    }

    private Server server;

    public void start() throws IOException {
        /* The port on which the server should run */
        int port = 50061;
        server = ServerBuilder.forPort(port)
                .addService(new OrderEntryServiceImpl(this.exchange))
                .addService(ProtoReflectionService.newInstance())
                .build()
                .start();
        logger.info("order entry data service started, listening on " + port);
        Runtime.getRuntime().addShutdownHook(new Thread() {
            @Override
            public void run() {
                // Use stderr here since the logger may have been reset by its JVM shutdown hook.
                System.err.println("*** shutting down market data service gRPC server since JVM is shutting down");
                try {
                    OrderEntryService.this.stop();
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

    static class OrderEntryServiceImpl extends OrderEntryServiceGrpc.OrderEntryServiceImplBase {

        Exchange exchange;

        Map<String, String> clOrdIdToSymbol = new ConcurrentHashMap<>();
        Map<String, String> clOrdIdToBookId = new ConcurrentHashMap<>();

        OrderEntryServiceImpl(Exchange exchange) {
            this.exchange = exchange;
        }

        @Override
        public void submitNewOrder(Orderentryapi.NewOrderParams request, StreamObserver<Orderentryapi.OrderId> responseObserver) {

            var book = this.exchange.getOrderBook(request.getSymbol());


            clOrdIdToSymbol.put(request.getClOrderId(), request.getSymbol());

            Side side = request.getOrderSide() == Orderentryapi.Side.BUY ? Side.Buy : Side.Sell;
            var qntB = BigInteger.valueOf(request.getQuantity().getMantissa());
            var exp = BigInteger.valueOf(10).pow(request.getQuantity().getExponent());
            int qty = qntB.multiply(exp).intValue();

            BigDecimal price = BigDecimal.valueOf(request.getPrice().getMantissa(), -1 * request.getPrice().getExponent());

            var id = book.addOrder(side, qty, price, request.getClOrderId());
            clOrdIdToBookId.put(request.getClOrderId(), id);

            var orderId = Orderentryapi.OrderId.newBuilder().setOrderId(id).build();

            responseObserver.onNext(orderId);
            responseObserver.onCompleted();
        }

        @Override
        public void cancelOrder(Orderentryapi.OrderId request, StreamObserver<Orderentryapi.Empty> responseObserver) {

            var symbol = clOrdIdToSymbol.get(request.getOrderId());
            if (symbol == null) {
                responseObserver.onError(new Exception("symbol not found for client order id:" + request.getOrderId()));
            }

            var localOrderId = clOrdIdToBookId.get(request.getOrderId());
            if (localOrderId == null) {
                responseObserver.onError(new Exception("local order id not found for client order id:" + request.getOrderId()));
            }

            var book = this.exchange.getOrderBook(symbol);

            try {
                book.deleteOrder(localOrderId);
            } catch (Exception e) {
                responseObserver.onError(new Exception("failed to cancel order " + request.getOrderId(), e));
            }

            responseObserver.onNext(Orderentryapi.Empty.newBuilder().build());
            responseObserver.onCompleted();

        }
    }

}


