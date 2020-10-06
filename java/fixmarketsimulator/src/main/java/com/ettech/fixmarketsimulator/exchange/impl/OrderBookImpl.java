package com.ettech.fixmarketsimulator.exchange.impl;

import com.ettech.fixmarketsimulator.exchange.*;

import java.math.BigDecimal;
import java.util.*;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import quickfix.field.MDUpdateAction;
import quickfix.fix50.component.BidCompReqGrp;

public class OrderBookImpl implements OrderBook {

    static Logger log = LoggerFactory.getLogger(OrderBookImpl.class);

    private String instrument;

    List<LimitOrderImpl> buyOrders = new ArrayList<>();
    List<LimitOrderImpl> sellOrders = new ArrayList<>();

    List<TradeListener> tradeListeners = new LinkedList<>();
    List<MdEntryListener> mdEntryListeners = new LinkedList<>();

    Trade lastTrade;
    double totalTradedVolume = 0;

    OrderBookImpl(String instrument) {
        this.instrument = instrument;
    }


    @Override
    public synchronized void addTradeListenerIfNotRegistered(TradeListener listener) {
        if (!tradeListeners.contains(listener)) {
            tradeListeners.add(listener);
        }
    }

    @Override
    public synchronized void addMdEntryListener(MdEntryListener listener) {
        mdEntryListeners.add(listener);
    }

    @Override
    public synchronized void removeMdEntryListener(MdEntryListener listener) {
        mdEntryListeners.remove(listener);
    }

    @Override
    public synchronized Order[] getBuyOrders() {
        return buyOrders.toArray(new Order[]{});
    }

    @Override
    public synchronized Order[] getSellOrders() {
        return sellOrders.toArray(new Order[]{});
    }

    @Override
    public Trade getLastTrade() {
        return this.lastTrade;
    }

    @Override
    public double getTotalTradedVolume() {
        return totalTradedVolume;
    }

    @Override
    public synchronized OrderState modifyOrder(String orderId, BigDecimal newPrice, int newQuantity)
            throws OrderModificationException {

        LimitOrderImpl originalOrder = getOrder(orderId);

        if (newPrice == null) {
            newPrice = originalOrder.getPrice();
        }

        if (originalOrder.getQuantity() > newQuantity) {
            double qtyReduction = originalOrder.getQuantity() - newQuantity;
            if (qtyReduction >= originalOrder.getRemainingQty()) {
                try {
                    deleteOrder(orderId);
                } catch (OrderDeletionException e) {
                    throw new RuntimeException("Unexpected failure to delete order during order modification",
                            e);
                }
                return new OrderDeleteImpl(originalOrder.getOrderId(), 0, newQuantity);
            }
        }

        if (originalOrder.getPrice().compareTo(newPrice) == 0) {
            double qntChange = newQuantity - originalOrder.getQuantity();
            originalOrder.setQuantity(originalOrder.getQuantity() + qntChange);
            originalOrder.setRemainingQty(originalOrder.getRemainingQty() + qntChange);

            sendMdEntry(originalOrder, MdUpdateActionType.Modify);

        } else {
            try {
                deleteOrder(orderId);
                addOrderWithId(originalOrder.getSide(), newQuantity, newPrice, originalOrder.getClOrdId(),
                        originalOrder.getOrderId(),
                        true);

            } catch (OrderDeletionException e) {
                throw new RuntimeException("Unexpected failure to delete order during order modification",
                        e);
            }
        }

        return new OrderDeleteImpl(originalOrder.getClOrdId(), originalOrder.getRemainingQty(),
                originalOrder.getQuantity());

    }

    private LimitOrderImpl getOrder(String orderId) throws OrderModificationException {
        Optional<LimitOrderImpl> optFoundOrder = buyOrders.stream().filter(o -> o.getOrderId().equals(orderId))
                .findFirst();
        if (optFoundOrder.isEmpty()) {
            optFoundOrder = sellOrders.stream().filter(o -> o.getOrderId().equals(orderId)).findFirst();
        }

        if (optFoundOrder.isEmpty()) {
            throw new OrderModificationException(
                    "Request to modify order failed as no order with id " + orderId + " exists in the book");
        }

        return optFoundOrder.get();
    }

    @Override
    public synchronized OrderState deleteOrder(String orderId) throws OrderDeletionException {

        Order deletedOrder = deleteAndReturnOrder(orderId);

        if (deletedOrder != null) {
            sendMdEntry(deletedOrder, MdUpdateActionType.Remove);
        } else {
            throw new OrderDeletionException(
                    "Request to delete order failed as no order with id " + orderId + " exists in the book");
        }

        return new OrderDeleteImpl(deletedOrder.getOrderId(), deletedOrder.getRemainingQty(),
                deletedOrder.getQuantity());

    }

    private void sendMdEntry(Order order, MdUpdateActionType mdUpdateActionType) {
        MDEntry entry = createMdEntry(order, mdUpdateActionType);
        List<MDEntry> entries = new ArrayList<>();
        entries.add(entry);
        dispatchMdEntries(entries);
    }


    Order deleteAndReturnOrder(String orderId) {

        Order removedOrder = removeFromOrders(buyOrders, orderId);
        if (removedOrder == null) {
            removedOrder = removeFromOrders(sellOrders, orderId);
        }

        return removedOrder;
    }

    Order removeFromOrders(List<LimitOrderImpl> orders, String orderId) {
        Optional<LimitOrderImpl> optFoundOrder = orders.stream().filter(o -> o.getOrderId().equals(orderId))
                .findFirst();

        if (optFoundOrder.isEmpty()) {
            return null;
        }

        Order order = optFoundOrder.get();
        orders.remove(order);

        return order;
    }


    @Override
    public synchronized String addOrder(Side side, int qty, BigDecimal price, String clOrderId) {

        log.debug("Add Order  symbol:{} qty:{} price:{} side:{} clOrderId:{}", this.instrument, qty, price, side, clOrderId);

        String orderId = UUID.randomUUID().toString();

        return addOrderWithId(side, qty, price, clOrderId, orderId, false).getOrderId();
    }

    private Order addOrderWithId(Side side, int qty, BigDecimal price, String clOrderId,
                                 String orderId, boolean fromModifyOperation) {
        LimitOrderImpl newOrder = new LimitOrderImpl(qty, price, clOrderId, side, orderId);


        checkForCrosses(newOrder, fromModifyOperation);

        if (newOrder.getRemainingQty() > 0) {
            addOrderToBook(newOrder);
        }

        return newOrder;
    }

    private void checkForCrosses(LimitOrderImpl newOrder, boolean newOrderFromModifyOperation) {
        List<MDEntry> mdEntries = new ArrayList<>();
        List<Trade> trades = new ArrayList<>();

        List<LimitOrderImpl> oppSideOrders = newOrder.getSide() == Side.Buy ? sellOrders : buyOrders;
        List<Integer> oppSideIndexOfOrdersToRemove = new ArrayList<>();

        for (int i = 0; i < oppSideOrders.size(); i++) {
            LimitOrderImpl oppSideOrder = oppSideOrders.get(i);

            if (newOrder.getRemainingQty() > 0 &&
                    ((newOrder.getSide() == Side.Buy && newOrder.getPrice().compareTo(oppSideOrder.getPrice()) != -1) ||
                            (newOrder.getSide() == Side.Sell && newOrder.getPrice()
                                    .compareTo(oppSideOrder.getPrice()) != 1))) {

                double quantity =
                        newOrder.getRemainingQty() > oppSideOrder.getRemainingQty() ? oppSideOrder
                                .getRemainingQty()
                                : newOrder.getRemainingQty();
                BigDecimal price = oppSideOrder.getPrice();

                newOrder.setRemainingQty(newOrder.getRemainingQty() - quantity);


                if (oppSideOrder.getRemainingQty() - quantity == 0) {
                    oppSideIndexOfOrdersToRemove.add(i);
                    mdEntries.add(createMdEntry(oppSideOrder, MdUpdateActionType.Remove));
                    // The remove entry requires the entries quantity, hence the adjustment to the remaining qty
                    // is done after the creation the mdEntry
                    oppSideOrder.setRemainingQty(oppSideOrder.getRemainingQty() - quantity);
                } else {
                    oppSideOrder.setRemainingQty(oppSideOrder.getRemainingQty() - quantity);
                    mdEntries.add(createMdEntry(oppSideOrder, MdUpdateActionType.Modify));
                }

                String tradeId = UUID.randomUUID().toString();
                trades.add(new TradeImpl(tradeId, oppSideOrder.getClOrdId(), price, quantity, instrument,
                        oppSideOrder.getSide(), oppSideOrder.getOrderId(), oppSideOrder.getRemainingQty(),
                        oppSideOrder.getQuantity() - oppSideOrder.getRemainingQty()));
                trades.add(
                        new TradeImpl(tradeId, newOrder.getClOrdId(), price, quantity, instrument,
                                newOrder.getSide(),
                                newOrder.getOrderId(), newOrder.getRemainingQty(),
                                newOrder.getQuantity() - newOrder.getRemainingQty()));

            } else {
                break;
            }

        }

        if (!oppSideIndexOfOrdersToRemove.isEmpty()) {
            // Remove empty orders
            List<LimitOrderImpl> updatedOrders = new ArrayList<>();
            for (int i = 0; i < oppSideOrders.size(); i++) {
                if (!oppSideIndexOfOrdersToRemove.contains(i)) {
                    updatedOrders.add(oppSideOrders.get(i));
                }
            }

            if (newOrder.getSide() == Side.Buy) {
                sellOrders = updatedOrders;
            } else {
                buyOrders = updatedOrders;
            }

        }

        if (newOrder.getRemainingQty() > 0) {
            mdEntries.add(createMdEntry(newOrder,
                    newOrderFromModifyOperation ? MdUpdateActionType.Modify : MdUpdateActionType.Add));
        }


        List<Trade> uniqueTrades = new ArrayList<>();
        Set<String> encountered = new HashSet<String>();
        for (var trade : trades) {
            if (!encountered.contains(trade.getTradeId())) {
                encountered.add(trade.getTradeId());
                uniqueTrades.add(trade);
            }
        }

        for (var trade : uniqueTrades) {
           lastTrade = trade;
           totalTradedVolume += trade.getQuantity();

           var mdTradeEntry =  new MDEntryImpl(MdUpdateActionType.Modify, trade.getTradeId(), trade.getPrice(),
                   trade.getQuantity(),
                    instrument,
                    MdEntryType.Trade, "");

           mdEntries.add(mdTradeEntry);
        }

        if( uniqueTrades.size() > 0 ) {
            var mdTotalTradeEntry =  new MDEntryImpl(MdUpdateActionType.Modify, UUID.randomUUID().toString(), lastTrade.getPrice(),
                    totalTradedVolume,
                    instrument,
                    MdEntryType.TradeVolume, "");

            mdEntries.add(mdTotalTradeEntry);
        }

        dispatchTrades(trades);
        dispatchMdEntries(mdEntries);
    }

    private MDEntryImpl createMdEntry(Order order, MdUpdateActionType entryType) {
        return new MDEntryImpl(entryType, order.getOrderId(), order.getPrice(),
                order.getRemainingQty(),
                instrument,
                getMdEntryTypeFromSide(order.getSide()), order.getClOrdId());
    }

    public static MdEntryType getMdEntryTypeFromSide(Side side) {
        return side == Side.Buy ? MdEntryType.Bid : MdEntryType.Offer;
    }

    private void addOrderToBook(LimitOrderImpl newOrder) {
        if (newOrder.getSide() == Side.Buy) {

            int insertionIndex = 0;
            for (; insertionIndex < buyOrders.size(); insertionIndex++) {
                if (buyOrders.get(insertionIndex).getPrice().compareTo(newOrder.getPrice()) == -1) {
                    break;
                }
            }

            buyOrders.add(insertionIndex, newOrder);
        } else {

            int insertionIndex = 0;
            for (; insertionIndex < sellOrders.size(); insertionIndex++) {
                if (sellOrders.get(insertionIndex).getPrice().compareTo(newOrder.getPrice()) == 1) {
                    break;
                }
            }

            sellOrders.add(insertionIndex, newOrder);
        }
    }

    void dispatchMdEntries(List<MDEntry> mdEntries) {
        mdEntryListeners.forEach(l -> l.onMdEntries(mdEntries));
    }

    void dispatchTrades(List<Trade> trades) {
        tradeListeners.forEach(l -> l.onTrades(trades));
    }


    @Override
    public String getInstrument() {
        return instrument;
    }
}
