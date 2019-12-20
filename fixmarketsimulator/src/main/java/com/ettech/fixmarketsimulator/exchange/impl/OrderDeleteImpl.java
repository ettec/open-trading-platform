package com.ettech.fixmarketsimulator.exchange.impl;

import com.ettech.fixmarketsimulator.exchange.OrderState;

public class OrderDeleteImpl implements OrderState {

  public OrderDeleteImpl(String orderId, double leavesQty, double quantity) {
    this.orderId = orderId;
    this.leavesQty = leavesQty;
    this.qty = quantity;
  }

  String orderId;
  double leavesQty;
  double qty;




  @Override
  public String getOrderId() {
    return orderId;
  }

  @Override
  public double getLeavesQty() {
    return leavesQty;
  }

  @Override
  public double getQty() {
    return qty;
  }
}
