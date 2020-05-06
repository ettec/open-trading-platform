package com.ettech.fixmarketsimulator.exchange;

public interface OrderState {

  String getOrderId();
  double getLeavesQty();
  double getQty();

}
