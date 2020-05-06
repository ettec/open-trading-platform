package com.ettech.fixmarketsimulator.exchange;

import java.math.BigDecimal;

public interface Trade {

  String getInstrument();

  String getClOrderId();

  BigDecimal getPrice();

  double getQuantity();

  Side getOrderSide();

  String getOrderId();

  String getTradeId();

  double getLeavesQty();

  double getCumQty();

}
