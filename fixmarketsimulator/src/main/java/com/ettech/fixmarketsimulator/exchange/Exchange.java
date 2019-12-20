package com.ettech.fixmarketsimulator.exchange;

public interface Exchange {

    OrderBook getOrderBook(String instrument);

}
