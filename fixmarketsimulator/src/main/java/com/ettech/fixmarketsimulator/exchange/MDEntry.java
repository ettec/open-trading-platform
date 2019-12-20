package com.ettech.fixmarketsimulator.exchange;

import java.math.BigDecimal;

public interface MDEntry {

  String getInstrument();

  MdEntryType getMdEntryType();

  String getOrderId();

  String getClOrderId();

  BigDecimal getPrice();

  double getQuantity();

  Side getSide();

}
