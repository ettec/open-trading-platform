package com.ettech.fixmarketsimulator.exchange;

import java.util.List;

public interface TradeListener {

  void onTrades(List<Trade> trade);

}
