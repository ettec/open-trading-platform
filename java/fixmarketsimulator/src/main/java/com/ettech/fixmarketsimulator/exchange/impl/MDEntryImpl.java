package com.ettech.fixmarketsimulator.exchange.impl;

import com.ettech.fixmarketsimulator.exchange.MdEntryType;
import com.ettech.fixmarketsimulator.exchange.Side;
import java.math.BigDecimal;
import java.util.Objects;

class MDEntryImpl implements com.ettech.fixmarketsimulator.exchange.MDEntry {

  private String instrument;
  private MdEntryType mdEntryType;
  private String orderId;
  private BigDecimal price;
  private double quantity;
  private Side side;
  private String clOrderId;

  public MDEntryImpl(MdEntryType mdEntryType, String orderId, BigDecimal price, double quantity, String instrument, Side side,
    String clOrderId) {
    this.mdEntryType = mdEntryType;
    this.orderId = orderId;
    this.price = price;
    this.quantity = quantity;
    this.instrument = instrument;
    this.side = side;
    this.clOrderId = clOrderId;
  }

  @Override
  public String getInstrument() {
    return instrument;
  }

  @Override
  public MdEntryType getMdEntryType() {
    return mdEntryType;
  }

  @Override
  public String getOrderId() {
    return orderId;
  }

  @Override
  public String getClOrderId() { return clOrderId; }

  @Override
  public BigDecimal getPrice() {
    return price;
  }

  @Override
  public double getQuantity() {
    return quantity;
  }

  @Override
  public Side getSide() {
    return side;
  }


}
