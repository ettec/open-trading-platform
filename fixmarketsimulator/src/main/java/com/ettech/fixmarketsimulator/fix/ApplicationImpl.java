package com.ettech.fixmarketsimulator.fix;

import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.exchange.OrderBook;
import com.ettech.fixmarketsimulator.exchange.OrderDeletionException;
import com.ettech.fixmarketsimulator.exchange.OrderModificationException;
import com.ettech.fixmarketsimulator.exchange.OrderState;
import com.ettech.fixmarketsimulator.exchange.Trade;
import com.ettech.fixmarketsimulator.exchange.TradeListener;

import java.math.BigDecimal;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import javax.inject.Inject;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import quickfix.Application;
import quickfix.DoNotSend;
import quickfix.FieldNotFound;
import quickfix.IncorrectDataFormat;
import quickfix.IncorrectTagValue;
import quickfix.Message;
import quickfix.MessageCracker;
import quickfix.RejectLogon;
import quickfix.Session;
import quickfix.SessionID;
import quickfix.SessionNotFound;
import quickfix.UnsupportedMessageType;
import quickfix.field.*;
import quickfix.fix50sp2.ExecutionReport;
import quickfix.fix50sp2.NewOrderSingle;
import quickfix.fix50sp2.OrderCancelReject;

public class ApplicationImpl extends MessageCracker implements Application, TradeListener {


    public static final Symbol SYMBOL = new Symbol();
    public static final Price PRICE = new Price();
    public static final Side SIDE = new Side();
    public static final ClOrdID CL_ORD_ID = new ClOrdID();
    public static final OrderQty ORDER_QTY = new OrderQty();


    Map<String, String> clOrderIdToOrderId = new HashMap<>();
    Map<String, String> clOrderIdToSymbol = new HashMap<>();
    Map<String, SessionID> clOrderIdToSession = new HashMap<>();


    @Inject
    Exchange exchange;


    static Logger log = LoggerFactory.getLogger(ApplicationImpl.class);


    @Override
    public void onCreate(SessionID sessionID) {

    }

    @Override
    public void onLogon(SessionID sessionID) {

    }

    @Override
    public void onLogout(SessionID sessionID) {

    }

    @Override
    public void toAdmin(Message message, SessionID sessionID) {

    }

    @Override
    public void fromAdmin(Message message, SessionID sessionID) throws FieldNotFound, IncorrectDataFormat, IncorrectTagValue, RejectLogon {

    }

    @Override
    public void toApp(Message message, SessionID sessionID) throws DoNotSend {

    }

    @Override
    public void fromApp(Message message, SessionID sessionID) throws FieldNotFound, IncorrectDataFormat, IncorrectTagValue, UnsupportedMessageType {
        log.info("Received: " + message);
        crack(message, sessionID);
    }

    @Override
    public void onTrades(List<Trade> trades) {

        for (Trade trade : trades) {

            ExecutionReport executionReport = new ExecutionReport(new OrderID(trade.getOrderId()),
                    new ExecID(trade.getTradeId()),
                    new ExecType(ExecType.TRADE),
                    new OrdStatus(trade.getLeavesQty() > 0 ? OrdStatus.PARTIALLY_FILLED : OrdStatus.FILLED),
                    new Side(trade.getOrderSide() == com.ettech.fixmarketsimulator.exchange.Side.Buy ? Side.BUY : Side.SELL),
                    new LeavesQty(trade.getLeavesQty()),
                    new CumQty(trade.getCumQty()));

            executionReport.set(new ClOrdID(trade.getClOrderId()));
            executionReport.set(new LastQty(trade.getQuantity()));
            executionReport.set(new LastPx(trade.getPrice().doubleValue()));

            try {
                Session.sendToTarget(executionReport, clOrderIdToSession.get(trade.getClOrderId()));
            } catch (SessionNotFound e) {
                log.error("Unable to send trade for client order id {} as the session is not found", trade.getClOrderId());
            }
        }

    }

    public void onMessage(quickfix.fix50sp2.OrderCancelReplaceRequest replaceRequest, SessionID sessionID) {


        try {

            String clOrdId = replaceRequest.get(CL_ORD_ID).getValue();
            Side side = replaceRequest.get(SIDE);
            OrderQty qty = replaceRequest.get(ORDER_QTY);
            BigDecimal newPrice = null;
            if (replaceRequest.isSet(PRICE)) {
                Price price = replaceRequest.get(PRICE);
                newPrice = new BigDecimal(price.getObject());
            }

            OrderState orderState;
            String orderId = clOrderIdToOrderId.get(clOrdId);
            var symbol = clOrderIdToSymbol.get(clOrdId);
            try {


                if (orderId == null) {
                    throw new OrderModificationException("No order found for client order id:" + clOrdId);
                }


                OrderBook orderBook = exchange.getOrderBook(symbol);
                orderState = orderBook.modifyOrder(orderId, newPrice, (int) qty.getValue());

            } catch (OrderModificationException e) {
                OrderCancelReject reject = new OrderCancelReject();
                reject.set(new OrderID("NONE"));
                reject.set(new ClOrdID(clOrdId));
                reject.set(new OrdStatus(OrdStatus.REJECTED));
                reject.set(new CxlRejResponseTo(CxlRejResponseTo.ORDER_CANCEL_REPLACE_REQUEST));
                Session.sendToTarget(reject, sessionID);
                return;

            }

            ExecutionReport executionReport = new ExecutionReport();
            executionReport.set(replaceRequest.get(CL_ORD_ID));
            executionReport.set(new OrderID(orderId));
            executionReport.set(new ExecID(UUID.randomUUID().toString()));
            executionReport.set(new ExecType(ExecType.REPLACED));
            executionReport.set(new OrdStatus(OrdStatus.REPLACED));
            executionReport.set(new Symbol(symbol));
            executionReport.set(side);

            executionReport.set(new LeavesQty(orderState.getLeavesQty()));
            executionReport.set(new CumQty(orderState.getQty() - orderState.getLeavesQty()));

            Session.sendToTarget(executionReport, sessionID);

        } catch (Exception e) {
            log.error("Failed to process order replace request:" + replaceRequest, e);
        }


    }


    public void onMessage(quickfix.fix50sp2.OrderCancelRequest cancelRequest, SessionID sessionID) {

        try {

            String clOrdId = cancelRequest.get(CL_ORD_ID).getValue();

            String orderId = clOrderIdToOrderId.get(clOrdId);
            if (orderId == null) {
                rejectCancelRequest(sessionID, clOrdId, "No order found for client order id:" + clOrdId, "NONE");
                return;
            }

            var symbol = clOrderIdToSymbol.get(clOrdId);
            OrderBook orderBook = exchange.getOrderBook(symbol);

            OrderState orderState;
            try {
                orderState = orderBook.deleteOrder(orderId);
            } catch (OrderDeletionException e) {
                rejectCancelRequest(sessionID, clOrdId, e.getMessage(), orderId);
                return;
            }

            ExecutionReport executionReport = new ExecutionReport(new OrderID(orderState.getOrderId()),
                    new ExecID(UUID.randomUUID().toString()),
                    new ExecType(ExecType.CANCELED),
                    new OrdStatus(OrdStatus.CANCELED),
                    cancelRequest.getSide(),
                    new LeavesQty(orderState.getLeavesQty()),
                    new CumQty(orderState.getQty() - orderState.getLeavesQty()));

            executionReport.set(cancelRequest.get(CL_ORD_ID));
            executionReport.set(new Symbol(symbol));

            Session.sendToTarget(executionReport, sessionID);

        } catch (Exception e) {
            log.error("Failed to process order cancel request:" + cancelRequest, e);
        }

    }

    private void rejectCancelRequest(SessionID sessionID, String clOrdId, String message,
                                     String rejectedOrderId) throws SessionNotFound {
        OrderCancelReject reject = new OrderCancelReject();
        reject.set(new OrderID(rejectedOrderId));
        reject.set(new ClOrdID(clOrdId));
        reject.set(new OrdStatus(OrdStatus.REJECTED));
        reject.set(new CxlRejResponseTo(CxlRejResponseTo.ORDER_CANCEL_REQUEST));
        reject.set(new Text(message));
        Session.sendToTarget(reject, sessionID);
    }


    public void onMessage(quickfix.fix50sp2.NewOrderSingle order, SessionID sessionID) {

        try {
            handleNewOrderSingle(order, sessionID);
        } catch (RejectOrderException e) {
            rejectNewOrder(order, sessionID, e.getMessage());
        }

    }

    private void handleNewOrderSingle(NewOrderSingle order, SessionID sessionID) throws RejectOrderException {
        try {

            OrderQty qty = order.get(ORDER_QTY);

            if (qty.getValue() <= 0) {
                throw new RejectOrderException("Quantity must be greater than 0");
            }

            ClOrdID clOrdId = order.get(CL_ORD_ID);
            Symbol symbol = order.get(SYMBOL);

            clOrderIdToSession.put(clOrdId.getValue(), sessionID);

            if (!order.isSet(PRICE)) {
                throw new RejectOrderException("Price must be set");
            }

            Price price = order.get(PRICE);

            Side side = order.get(SIDE);

            com.ettech.fixmarketsimulator.exchange.Side exSide = null;

            if (side.getValue() == Side.BUY) {
                exSide = com.ettech.fixmarketsimulator.exchange.Side.Buy;
            } else if (side.getValue() == Side.SELL) {
                exSide = com.ettech.fixmarketsimulator.exchange.Side.Sell;
            } else {
                throw new RejectOrderException("Side not supported:" + side);
            }


            OrderBook orderBook = exchange.getOrderBook(symbol.getValue());
            orderBook.addTradeListenerIfNotRegistered(this);

            String orderId = orderBook.addOrder(exSide, (int) qty.getValue(), new BigDecimal(price.getValue()), clOrdId.getValue());

            if (clOrderIdToOrderId.containsKey(clOrdId)) {
                throw new RejectOrderException("Already received clOrdId:" + clOrdId);
            }

            clOrderIdToOrderId.put(clOrdId.getValue(), orderId);
            clOrderIdToSymbol.put(clOrdId.getValue(), symbol.getValue());


            ExecutionReport executionReport = new ExecutionReport();
            executionReport.set(clOrdId);
            executionReport.set(new OrderID(orderId));
            executionReport.set(new ExecID(UUID.randomUUID().toString()));
            executionReport.set(new ExecType(ExecType.NEW));
            executionReport.set(new OrdStatus(OrdStatus.NEW));
            executionReport.set(symbol);
            executionReport.set(side);
            executionReport.set(new LeavesQty(qty.getValue()));
            executionReport.set(new CumQty(0));
            executionReport.set(new AvgPx(0));

            Session.sendToTarget(executionReport, sessionID);


        } catch (Exception e) {
            log.error("Failed to process new order single:" + order, e);
        }
    }

    public void rejectNewOrder(NewOrderSingle order, SessionID sessionID, String reason) {

        try {

            ExecutionReport executionReport = new ExecutionReport();
            executionReport.set(order.get(CL_ORD_ID));
            executionReport.set(new OrderID(""));
            executionReport.set(new ExecID(UUID.randomUUID().toString()));
            executionReport.set(new ExecType(ExecType.REJECTED));
            executionReport.set(new OrdStatus(OrdStatus.REJECTED));
            executionReport.set(order.get(SYMBOL));
            executionReport.set(order.get(SIDE));
            executionReport.set(new LeavesQty(0));
            executionReport.set(new CumQty(0));
            executionReport.set(new AvgPx(0));

            executionReport.set(new OrdRejReason(OrdRejReason.OTHER));
            executionReport.set(new Text(reason));

            Session.sendToTarget(executionReport, sessionID);
        } catch (Exception e) {
            log.error("Failed to reject order new order single:" + order, e);
        }

    }


}
