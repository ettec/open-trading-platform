package com.ettech.fixmarketsimulator.bookbuilder;

import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.exchange.OrderBook;
import com.ettech.fixmarketsimulator.exchange.OrderDeletionException;
import com.ettech.fixmarketsimulator.exchange.Side;
import com.google.gson.Gson;

import java.io.FileReader;
import java.math.BigDecimal;
import java.math.RoundingMode;
import java.util.*;
import java.util.concurrent.ScheduledThreadPoolExecutor;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;

public class BooksBuilder {

    private static final Logger logger = Logger.getLogger(BooksBuilder.class.getName());

    private Exchange exchange;
    private ScheduledThreadPoolExecutor executor;

    private List<BookBuilder> builders = new ArrayList<>();

    public BooksBuilder(Exchange exchange, String depthDataPath, String symbolsToRunStr, int numExecutorThreads,
                        long updateInterval,
                        double minQty, double variation, int tickScale, double tradeProbability, int maxDepth,
                        double cancelProbability) throws Exception {
        this.exchange = exchange;

        FileReader fr = new FileReader(depthDataPath);
        Gson g = new Gson();
        Depth[] depths = g.fromJson(fr, Depth[].class);
        Map<String, Depth> symToDepth = new HashMap();
        Arrays.stream(depths).forEach(depth -> symToDepth.put(depth.symbol, depth));

        Set<String> symsToRun = new HashSet<>();
        if( symbolsToRunStr.equals("*")) {
            symsToRun.addAll(symToDepth.keySet());
        } else {
            symsToRun.addAll(Arrays.asList(symbolsToRunStr.split(",")));
        }

        this.executor = new ScheduledThreadPoolExecutor(numExecutorThreads);


        symsToRun.forEach(s->{
            builders.add(new BookBuilder(exchange.getOrderBook(s), symToDepth.get(s), executor,
                    updateInterval, minQty, variation, tickScale, tradeProbability, maxDepth,
                    cancelProbability));
        });

    }


    static class BookBuilder {

        private static final String BOOK_BUILDER_ORDERID_PREPEND = "BOOKBUILDER:";

        BookBuilder(OrderBook book, Depth initialDepth, ScheduledThreadPoolExecutor se, long updateInterval,
                    double minQty, double variation, int tickScale, double tradeProbability, int maxDepth,
                    double cancelProbability) {

            initialDepth.bids.forEach(b->{
                var price = new BigDecimal( b.price).setScale(tickScale,RoundingMode.HALF_EVEN);
                book.addOrder(Side.Buy, b.size, price, newOrderId());
            });
            initialDepth.asks.forEach(b->{
                var price = new BigDecimal( b.price).setScale(tickScale,RoundingMode.HALF_EVEN);
                book.addOrder(Side.Sell, b.size,  price, newOrderId());
            });

            int totalBidQty = initialDepth.bids.stream().mapToInt(b->b.size).sum();
            int totalAskQty = initialDepth.asks.stream().mapToInt(b->b.size).sum();

            se.scheduleAtFixedRate(()->{
                try {

                    // Update based on qty
                    updateBookQty(book, minQty, variation, tickScale, totalBidQty, book.getBuyOrders(), initialDepth.bids,
                            Side.Buy, maxDepth);
                    updateBookQty(book, minQty, variation, tickScale, totalAskQty, book.getSellOrders(), initialDepth.asks,
                            Side.Sell, maxDepth);

                    // Update based on trade probability
                    if (Math.random() < tradeProbability) {
                        hitTopOfBook(book, book.getSellOrders(), Side.Buy);
                    }

                    if (Math.random() < tradeProbability) {
                        hitTopOfBook(book, book.getBuyOrders(), Side.Sell);
                    }

                    if (Math.random() < cancelProbability) {
                        cancelOrder(book, book.getSellOrders());
                    }

                    if (Math.random() < cancelProbability) {
                        cancelOrder(book, book.getBuyOrders());
                    }
                } catch( Exception e) {
                    logger.severe("Exception in book builder" + e);
                }


            }, updateInterval, updateInterval, TimeUnit.MILLISECONDS);
        }


        private String newOrderId() {
            return BOOK_BUILDER_ORDERID_PREPEND  + UUID.randomUUID().toString();
        }

        private void cancelOrder(OrderBook book, com.ettech.fixmarketsimulator.exchange.Order[] orders) {
            if(orders.length > 0) {
                int idx = (int) Math.round((orders.length - 1) * Math.random());
                var order = orders[idx];
                if (order.getClOrdId().startsWith(BOOK_BUILDER_ORDERID_PREPEND)) {
                    try {
                        book.deleteOrder(order.getOrderId());
                    } catch (OrderDeletionException e) {
                        logger.info("failed to delete order: " + e);
                    }
                }
            }
        }

        private void hitTopOfBook(OrderBook book, com.ettech.fixmarketsimulator.exchange.Order[] orders, Side side) {
            if (orders.length > 0) {
                int numOrders = (int) (orders.length * 0.5 * Math.random());

                long qty=0;
                BigDecimal price = new BigDecimal(0);
                for( int i=0; i < numOrders; i++) {
                     qty += Math.round(orders[i].getRemainingQty());
                    price = orders[i].getPrice();
                }

                if( qty > 0 ) {
                    book.addOrder(side, (int) qty, price, newOrderId());
                }

            }
        }

        private void updateBookQty(OrderBook book, double minQty, double variation, int tickScale, int totalQty,
                                   com.ettech.fixmarketsimulator.exchange.Order[] orders,
                                   List<Depth.Line> lines, Side side, int maxDepth) {
            var qQty = Arrays.stream(orders).findFirst().stream().mapToDouble(o->o.getRemainingQty()).sum();
            if ( qQty < totalQty * minQty && orders.length < maxDepth) {
                int idx =(int)Math.random() * lines.size();
                var line = lines.get(idx);

                var price = line.price - (line.price * (Math.random()-0.5) * variation);
                var qty = (int) (line.size - (line.size * (Math.random()-0.5) * variation));

                var bdPrice = new BigDecimal(price);
                bdPrice = bdPrice.setScale(tickScale, RoundingMode.HALF_EVEN);

                book.addOrder(side, qty, bdPrice, newOrderId() );
            }
        }
    }

}
