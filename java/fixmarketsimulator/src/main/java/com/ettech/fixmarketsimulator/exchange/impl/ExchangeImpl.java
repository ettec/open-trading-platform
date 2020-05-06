package com.ettech.fixmarketsimulator.exchange.impl;

import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.exchange.OrderBook;

import com.ettech.fixmarketsimulator.exchange.Trade;
import com.ettech.fixmarketsimulator.exchange.TradeListener;
import java.util.HashMap;
import java.util.Map;

public class ExchangeImpl implements Exchange {

    Map<String, OrderBook> instrumentToOrderBook = new HashMap<>();

    @Override
    public synchronized OrderBook  getOrderBook(String instrument) {

        OrderBook orderBook = instrumentToOrderBook.get(instrument);
        if( orderBook == null) {
            orderBook = new OrderBookImpl(instrument);
            instrumentToOrderBook.put(instrument, orderBook);
        }

        return orderBook;
    }

}
