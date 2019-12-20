package com.ettech.fixmarketsimulator.exchange.impl;

import com.ettech.fixmarketsimulator.exchange.Side;
import io.swagger.v3.oas.models.media.UUIDSchema;
import java.math.BigDecimal;
import java.util.UUID;

class TestOrder {



  public TestOrder( Side side, int qty, int price) {
    this(side, qty, new BigDecimal(price), UUID.randomUUID().toString());
  }


  public TestOrder( Side side, int qty, int price,
      String clOrderId) {
   this(side, qty, new BigDecimal(price), clOrderId);
  }


  public TestOrder( Side side, int qty, BigDecimal price,
      String clOrderId) {

    this.side = side;
    this.qty = qty;
    this.price = price;
    this.clOrderId = clOrderId;
  }


  Side side;
  int qty;
  BigDecimal price;
  int ordinal;
  String clOrderId;



}
