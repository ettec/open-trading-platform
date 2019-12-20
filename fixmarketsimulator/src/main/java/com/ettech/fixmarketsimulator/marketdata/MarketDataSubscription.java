package com.ettech.fixmarketsimulator.marketdata;

import ch.obermuhlner.math.big.BigDecimalMath;
import com.ettech.fixmarketsimulator.exchange.MDEntry;
import com.ettech.fixmarketsimulator.exchange.MdEntryListener;
import com.ettech.fixmarketsimulator.exchange.MdEntryType;
import com.ettech.fixmarketsimulator.exchange.Order;
import com.ettech.fixmarketsimulator.exchange.OrderBook;
import com.ettech.fixmarketsimulator.exchange.Side;
import io.netty.channel.ChannelHandlerContext;
import java.io.Closeable;
import java.io.IOException;
import java.math.BigDecimal;
import java.util.List;
import org.fixprotocol.components.Fix.Decimal64;
import org.fixprotocol.components.Instrument;
import org.fixprotocol.components.MarketData.MDEntryTypeEnum;
import org.fixprotocol.components.MarketData.MDIncGrp;
import org.fixprotocol.components.MarketData.MDIncGrp.Builder;
import org.fixprotocol.components.MarketData.MDUpdateActionEnum;
import org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class MarketDataSubscription implements Closeable, MdEntryListener {

  ChannelHandlerContext ctx;
  OrderBook book;
  String requestId;

  Logger log = LoggerFactory.getLogger(MarketDataSubscription.class);

  MarketDataSubscription(ChannelHandlerContext ctx, OrderBook book, String requestId) {

    this.ctx = ctx;
    this.book = book;
    this.requestId = requestId;

    MarketDataIncrementalRefresh.Builder incRefreshBuilder = MarketDataIncrementalRefresh.newBuilder();
    incRefreshBuilder.setMdReqId(requestId);
    var updateType = MDUpdateActionEnum.MD_UPDATE_ACTION_NEW;

    for (Order order : book.getBuyOrders()) {
      incRefreshBuilder.addMdIncGrp(getMdEntryFromOrder(book, order, updateType, Side.Buy));
    }

    for (Order order : book.getSellOrders()) {
      incRefreshBuilder.addMdIncGrp(getMdEntryFromOrder(book, order, updateType, Side.Sell));
    }

    book.addMdEntryListener(this);

    var incRefresh = incRefreshBuilder.build();
    ctx.write(incRefresh);
    log.info("Sent incremental refresh {}", incRefresh);
  }

  private MDIncGrp getMdEntryFromOrder(OrderBook book, Order order, MDUpdateActionEnum updateType, Side side) {
    var mdEntryBuilder = MDIncGrp.newBuilder();
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

  private void setQuantity(Builder mdEntryBuilder, long qntAsDouble) {
    long quantity = qntAsDouble;
    var qntBuilder = Decimal64.newBuilder();
    qntBuilder.setMantissa(quantity);
    qntBuilder.setExponent(0);
    mdEntryBuilder.setMdEntrySize(qntBuilder.build());
  }

  private void setPrice(Builder mdEntryBuilder, BigDecimal price) {
    var priceBuilder = Decimal64.newBuilder();
    int scale = -1*price.scale();
    priceBuilder.setExponent(scale);

    int unscaledValue = price.unscaledValue().intValue();

    priceBuilder.setMantissa(unscaledValue);

    mdEntryBuilder.setMdEntryPx(priceBuilder.build());
  }

  @Override
  public void close() throws IOException {
    book.removeMdEntryListener(this);
  }

  @Override
  public void onMdEntries(List<MDEntry> mdEntries) {
    MarketDataIncrementalRefresh.Builder incRefresh = MarketDataIncrementalRefresh.newBuilder();
    for (MDEntry entry : mdEntries) {

      var mdEntryBuilder = MDIncGrp.newBuilder();
      mdEntryBuilder.setMdUpdateAction(getMDEntryTypeEnum(entry.getMdEntryType()));
      mdEntryBuilder.setInstrument(Instrument.newBuilder().setSymbol(entry.getInstrument()).build());
      mdEntryBuilder.setMdEntryId(entry.getOrderId());

      mdEntryBuilder.setMdEntryType(getMdEntryType(entry.getSide()));

      setPrice(mdEntryBuilder, entry.getPrice());
      setQuantity(mdEntryBuilder, (long) entry.getQuantity());


      incRefresh.addMdIncGrp(mdEntryBuilder.build());
    }

    ctx.writeAndFlush(incRefresh.build());
  }

  private MDEntryTypeEnum getMdEntryType(Side side) {
    return side == Side.Buy ? MDEntryTypeEnum.MD_ENTRY_TYPE_BID
        : MDEntryTypeEnum.MD_ENTRY_TYPE_OFFER;
  }

  public MDUpdateActionEnum getMDEntryTypeEnum(MdEntryType entryType) {
    switch (entryType) {
      case Add:
        return MDUpdateActionEnum.MD_UPDATE_ACTION_NEW;
      case Modify:
        return MDUpdateActionEnum.MD_UPDATE_ACTION_CHANGE;
      case Remove:
        return MDUpdateActionEnum.MD_UPDATE_ACTION_DELETE;
      default:
        throw new RuntimeException("Unexpected entry type:" + entryType);
    }

  }

}
