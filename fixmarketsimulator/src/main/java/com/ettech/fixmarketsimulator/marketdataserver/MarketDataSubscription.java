package com.ettech.fixmarketsimulator.marketdataserver;

import com.ettech.fixmarketsimulator.exchange.*;
import org.fixprotocol.components.Fix;
import org.fixprotocol.components.Instrument;
import org.fixprotocol.components.MarketData;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.Closeable;
import java.io.IOException;
import java.math.BigDecimal;
import java.util.List;

public class MarketDataSubscription implements Closeable, MdEntryListener {

  Connection connection;
  OrderBook book;
  String requestId;

  Logger log = LoggerFactory.getLogger(MarketDataSubscription.class);

  public MarketDataSubscription(Connection connection, OrderBook book, String requestId) {

    this.connection = connection;
    this.book = book;
    this.requestId = requestId;

    MarketData.MarketDataIncrementalRefresh.Builder incRefreshBuilder = MarketData.MarketDataIncrementalRefresh.newBuilder();
    incRefreshBuilder.setMdReqId(requestId);
    var updateType = MarketData.MDUpdateActionEnum.MD_UPDATE_ACTION_NEW;

    for (Order order : book.getBuyOrders()) {
      incRefreshBuilder.addMdIncGrp(getMdEntryFromOrder(book, order, updateType, Side.Buy));
    }

    for (Order order : book.getSellOrders()) {
      incRefreshBuilder.addMdIncGrp(getMdEntryFromOrder(book, order, updateType, Side.Sell));
    }




    book.addMdEntryListener(this);

    var incRefresh = incRefreshBuilder.build();
    connection.send(incRefresh);
    log.info("Sent incremental refresh" + incRefresh);
  }

  private MarketData.MDIncGrp getMdEntryFromOrder(OrderBook book, Order order, MarketData.MDUpdateActionEnum updateType, Side side) {
    MarketData.MDIncGrp.Builder mdEntryBuilder = MarketData.MDIncGrp.newBuilder();
    mdEntryBuilder.setMdUpdateAction(updateType);

    mdEntryBuilder.setMdEntryId(order.getOrderId());

    BigDecimal price = order.getPrice();
    setPrice(mdEntryBuilder, price);

    setQuantity(mdEntryBuilder, (long) order.getRemainingQty());

    mdEntryBuilder.setMdEntryType(getMdEntryType(side));

    var instrument = Instrument.newBuilder().setSymbol(book.getInstrument());
    mdEntryBuilder.setInstrument(instrument.build());
    return mdEntryBuilder.build();
  }

  private void setQuantity(MarketData.MDIncGrp.Builder mdEntryBuilder, long qntAsDouble) {
    long quantity = qntAsDouble;
    var qntBuilder = Fix.Decimal64.newBuilder();
    qntBuilder.setMantissa(quantity);
    qntBuilder.setExponent(0);
    mdEntryBuilder.setMdEntrySize(qntBuilder.build());
  }

  private void setPrice(MarketData.MDIncGrp.Builder mdEntryBuilder, BigDecimal price) {
    Fix.Decimal64 fixPrice = getFixDecimal64(price);

    mdEntryBuilder.setMdEntryPx(fixPrice);
  }

  static  Fix.Decimal64 getFixDecimal64(BigDecimal price) {


    var str = price.toPlainString();

    var idx = str.indexOf('.');

    int exp = 0;
    if( idx > -1) {
      str = str.replace(".", "");
      exp = -(str.length() - idx);
    }

    long mantissa = Long.parseLong(str);

    var priceBuilder = Fix.Decimal64.newBuilder();
    priceBuilder.setExponent(exp);
    priceBuilder.setMantissa(mantissa);

    return priceBuilder.build();
  }

  public void close() {
    book.removeMdEntryListener(this);
  }

  @Override
  public void onMdEntries(List<MDEntry> mdEntries) {
    MarketData.MarketDataIncrementalRefresh.Builder incRefresh = MarketData.MarketDataIncrementalRefresh.newBuilder();
    for (MDEntry entry : mdEntries) {

      var mdEntryBuilder = MarketData.MDIncGrp.newBuilder();
      mdEntryBuilder.setMdUpdateAction(getMDEntryTypeEnum(entry.getMdEntryType()));
      mdEntryBuilder.setInstrument(Instrument.newBuilder().setSymbol(entry.getInstrument()).build());
      mdEntryBuilder.setMdEntryId(entry.getOrderId());

      mdEntryBuilder.setMdEntryType(getMdEntryType(entry.getSide()));

      setPrice(mdEntryBuilder, entry.getPrice());
      setQuantity(mdEntryBuilder, (long) entry.getQuantity());


      incRefresh.addMdIncGrp(mdEntryBuilder.build());
    }


    var refresh = incRefresh.build();

    this.log.info("SENDING INC REFRESH:" + refresh);


    connection.send(refresh);
  }

  private MarketData.MDEntryTypeEnum getMdEntryType(Side side) {
    return side == Side.Buy ? MarketData.MDEntryTypeEnum.MD_ENTRY_TYPE_BID
        : MarketData.MDEntryTypeEnum.MD_ENTRY_TYPE_OFFER;
  }

  public MarketData.MDUpdateActionEnum getMDEntryTypeEnum(MdEntryType entryType) {
    switch (entryType) {
      case Add:
        return MarketData.MDUpdateActionEnum.MD_UPDATE_ACTION_NEW;
      case Modify:
        return MarketData.MDUpdateActionEnum.MD_UPDATE_ACTION_CHANGE;
      case Remove:
        return MarketData.MDUpdateActionEnum.MD_UPDATE_ACTION_DELETE;
      default:
        throw new RuntimeException("Unexpected entry type:" + entryType);
    }

  }

}
