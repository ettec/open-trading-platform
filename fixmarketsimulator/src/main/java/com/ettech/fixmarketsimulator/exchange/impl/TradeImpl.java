package com.ettech.fixmarketsimulator.exchange.impl;

import com.ettech.fixmarketsimulator.exchange.Side;
import com.ettech.fixmarketsimulator.exchange.Trade;
import java.math.BigDecimal;

class TradeImpl implements Trade {

  public TradeImpl( String tradeId, String clOrderId, BigDecimal price,
      double quantity, String instrument, Side orderSide,
      String orderId,
      double leavesQty,
      double cumQty) {
    this.clOrderId = clOrderId;
    this.price = price;
    this.quantity = quantity;
    this.instrument = instrument;
    this.orderSide = orderSide;
    this.orderId = orderId;
    this.tradeId = tradeId;
    this.leavesQty = leavesQty;
    this.cumQty = cumQty;
  }

  private String tradeId;
  private String clOrderId;
  private BigDecimal price;
  private double quantity;
  private String instrument;
  private Side orderSide;
  private String orderId;
  private double leavesQty;
  private double cumQty;

  @Override
  public String getInstrument() {
    return instrument;
  }

  @Override
  public String getClOrderId() {
    return clOrderId;
  }

  @Override
  public BigDecimal getPrice() {
    return price;
  }

  @Override
  public double getQuantity() {
    return quantity;
  }

  @Override
  public Side getOrderSide() { return orderSide; }

  @Override
  public String getOrderId() { return orderId;  }

  @Override
  public String getTradeId() { return tradeId;  }

  public double getLeavesQty() {
    return leavesQty;
  }

  public double getCumQty() {
    return cumQty;
  }
}
