package com.ettech.fixmarketsimulator.bookbuilder;

import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.exchange.OrderBook;
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
                        double minQty, double variation, int tickScale, double tradeProbability, int maxDepth) throws Exception {
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
                    updateInterval, minQty, variation, tickScale, tradeProbability, maxDepth));
        });

    }

    static class BookBuilder {
        BookBuilder(OrderBook book, Depth initialDepth, ScheduledThreadPoolExecutor se, long updateInterval,
                    double minQty, double variation, int tickScale, double tradeProbability, int maxDepth) {


            initialDepth.bids.forEach(b->{
                var price = new BigDecimal( b.price).setScale(tickScale,RoundingMode.HALF_EVEN);
                book.addOrder(Side.Buy, b.size, price, UUID.randomUUID().toString());
            });
            initialDepth.asks.forEach(b->{
                var price = new BigDecimal( b.price).setScale(tickScale,RoundingMode.HALF_EVEN);
                book.addOrder(Side.Sell, b.size,  price, UUID.randomUUID().toString());
            });

            int totalBidQty = initialDepth.bids.stream().mapToInt(b->b.size).sum();
            int totalAskQty = initialDepth.asks.stream().mapToInt(b->b.size).sum();

            se.scheduleAtFixedRate(()->{

                // Update based on qty
                updateBookQty(book, minQty, variation, tickScale, totalBidQty, book.getBuyOrders(), initialDepth.bids,
                        Side.Buy, maxDepth);
                updateBookQty(book, minQty, variation, tickScale, totalAskQty, book.getSellOrders(), initialDepth.asks,
                        Side.Sell, maxDepth);

                // Update based on trade probability
                if (Math.random() < tradeProbability) {
                    hitTopOfBook(book, book.getSellOrders(), Side.Buy, tickScale);
                }

                if (Math.random() < tradeProbability) {
                    hitTopOfBook(book, book.getBuyOrders(), Side.Sell, tickScale);
                }

            }, updateInterval, updateInterval, TimeUnit.MILLISECONDS);
        }

        private void hitTopOfBook(OrderBook book, com.ettech.fixmarketsimulator.exchange.Order[] orders, Side side,
                                  int tickScale) {
            if (orders.length > 0) {
                var bestOpp = orders[0];

                var price = bestOpp.getPrice().setScale(tickScale);
                book.addOrder(side, (int) bestOpp.getRemainingQty(), price, UUID.randomUUID().toString());
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

                book.addOrder(side, qty, bdPrice, UUID.randomUUID().toString() );
            }
        }
    }

}
