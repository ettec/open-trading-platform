package com.ettech.fixmarketsimulator.exchange;

import java.math.BigDecimal;

public interface OrderBook {

  void addTradeListenerIfNotRegistered(TradeListener listener);

  void addMdEntryListener(MdEntryListener listener);

  void removeMdEntryListener(MdEntryListener listener);

  String addOrder(Side side, int qty, BigDecimal price, String clOrderId );

  OrderState deleteOrder(String orderId) throws OrderDeletionException;

  Order[] getBuyOrders();

  Order[] getSellOrders();

  OrderState modifyOrder(String orderId, BigDecimal newPrice, int newQuantity) throws OrderModificationException;

  String getInstrument();
}
