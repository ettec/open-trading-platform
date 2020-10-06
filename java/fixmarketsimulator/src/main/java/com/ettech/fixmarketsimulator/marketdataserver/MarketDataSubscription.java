package com.ettech.fixmarketsimulator.marketdataserver;

import com.ettech.fixmarketsimulator.exchange.*;
import com.ettech.fixmarketsimulator.exchange.impl.OrderBookImpl;
import org.fixprotocol.components.Fix;
import org.fixprotocol.components.Instrument;
import org.fixprotocol.components.MarketData;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.Closeable;
import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;

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

        var lastTrade = book.getLastTrade();
        if( lastTrade != null ) {
            MarketData.MDIncGrp.Builder mdEntryBuilder = MarketData.MDIncGrp.newBuilder();
            mdEntryBuilder.setMdUpdateAction(updateType);
            mdEntryBuilder.setMdEntryId(lastTrade.getTradeId());
            setPrice(mdEntryBuilder, lastTrade.getPrice());

            setQuantity(mdEntryBuilder, (long) lastTrade.getQuantity());

            mdEntryBuilder.setMdEntryType(getMdEntryType(MdEntryType.Trade));

            var instrument = Instrument.newBuilder().setSymbol(book.getInstrument());
            mdEntryBuilder.setInstrument(instrument.build());
            incRefreshBuilder.addMdIncGrp(mdEntryBuilder.build());
        }

        MarketData.MDIncGrp.Builder mdTradedVol = MarketData.MDIncGrp.newBuilder();
        mdTradedVol.setMdUpdateAction(updateType);
        mdTradedVol.setMdEntryId(UUID.randomUUID().toString());
        setQuantity(mdTradedVol, (long) book.getTotalTradedVolume());
        mdTradedVol.setMdEntryType(getMdEntryType(MdEntryType.TradeVolume));
        var instrument = Instrument.newBuilder().setSymbol(book.getInstrument());
        mdTradedVol.setInstrument(instrument.build());

        incRefreshBuilder.addMdIncGrp(mdTradedVol.build());



        book.addMdEntryListener(this);

        var incRefresh = incRefreshBuilder.build();
        connection.send(incRefresh);
    }

    private MarketData.MDIncGrp getMdEntryFromOrder(OrderBook book, Order order, MarketData.MDUpdateActionEnum updateType, Side side) {
        MarketData.MDIncGrp.Builder mdEntryBuilder = MarketData.MDIncGrp.newBuilder();
        mdEntryBuilder.setMdUpdateAction(updateType);

        mdEntryBuilder.setMdEntryId(order.getOrderId());

        BigDecimal price = order.getPrice();
        setPrice(mdEntryBuilder, price);

        setQuantity(mdEntryBuilder, (long) order.getRemainingQty());

        mdEntryBuilder.setMdEntryType(getMdEntryType(OrderBookImpl.getMdEntryTypeFromSide(side)));

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

    static Fix.Decimal64 getFixDecimal64(BigDecimal price) {


        var str = price.toPlainString();

        var idx = str.indexOf('.');

        int exp = 0;
        if (idx > -1) {
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
                mdEntryBuilder.setMdUpdateAction(getMDUpdateActionEnum(entry.getMdUpdateAction()));
                mdEntryBuilder.setInstrument(Instrument.newBuilder().setSymbol(entry.getInstrument()).build());
                mdEntryBuilder.setMdEntryId(entry.getId());

                mdEntryBuilder.setMdEntryType(getMdEntryType(entry.getMdEntryType()));

                setPrice(mdEntryBuilder, entry.getPrice());
                setQuantity(mdEntryBuilder, (long) entry.getQuantity());

                incRefresh.addMdIncGrp(mdEntryBuilder.build());
        }


        var refresh = incRefresh.build();

        connection.send(refresh);
    }

    public static MarketData.MDEntryTypeEnum getMdEntryType(MdEntryType type) {

        switch (type) {
            case Bid:
                return MarketData.MDEntryTypeEnum.MD_ENTRY_TYPE_BID;
            case Offer:
                return MarketData.MDEntryTypeEnum.MD_ENTRY_TYPE_OFFER;
            case Trade:
                return MarketData.MDEntryTypeEnum.MD_ENTRY_TYPE_TRADE;
            case TradeVolume:
                return MarketData.MDEntryTypeEnum.MD_ENTRY_TYPE_TRADE_VOLUME;
            default:
                throw new RuntimeException("MdEntry type not supported:" + type);
        }
    }

    public MarketData.MDUpdateActionEnum getMDUpdateActionEnum(MdUpdateActionType entryType) {
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
