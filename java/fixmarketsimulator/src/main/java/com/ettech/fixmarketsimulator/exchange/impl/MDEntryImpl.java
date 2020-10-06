package com.ettech.fixmarketsimulator.exchange.impl;

import com.ettech.fixmarketsimulator.exchange.MdEntryType;
import com.ettech.fixmarketsimulator.exchange.MdUpdateActionType;

import java.math.BigDecimal;

class MDEntryImpl implements com.ettech.fixmarketsimulator.exchange.MDEntry {

  private String instrument;
  private MdUpdateActionType mdUpdateAction;
  private String id;
  private BigDecimal price;
  private double quantity;
  private String clOrderId;
  private MdEntryType mdEntryType;

  public MDEntryImpl(MdUpdateActionType mdUpdateAction, String id, BigDecimal price, double quantity, String instrument, MdEntryType
          mdEntryType, String clOrderId) {
    this.mdUpdateAction = mdUpdateAction;
    this.id = id;
    this.price = price;
    this.quantity = quantity;
    this.instrument = instrument;
    this.mdEntryType = mdEntryType;
    this.clOrderId = clOrderId;
  }

  @Override
  public String getInstrument() {
    return instrument;
  }

  @Override
  public MdUpdateActionType getMdUpdateAction() {
    return mdUpdateAction;
  }

  @Override
  public MdEntryType getMdEntryType() {
    return mdEntryType;
  }

  @Override
  public String getId() {
    return id;
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


}
