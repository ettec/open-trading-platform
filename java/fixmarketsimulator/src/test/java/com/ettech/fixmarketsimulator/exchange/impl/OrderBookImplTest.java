package com.ettech.fixmarketsimulator.exchange.impl;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assertions.fail;

import com.ettech.fixmarketsimulator.exchange.*;

import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

class OrderBookImplTest {

  public static final String TESTINSTRUMENT = "testinstrument";

  List<Trade> trades ;
  List<MDEntry> mdEntries ;
  OrderBookImpl orderBook ;

  @BeforeEach
  public void setupBeforeEachTest() {

    trades = new ArrayList<>();
    mdEntries = new ArrayList<>();
    orderBook = new OrderBookImpl(TESTINSTRUMENT);

    orderBook.addTradeListenerIfNotRegistered(t->trades.addAll(t));
    orderBook.addMdEntryListener(m->mdEntries.addAll(m));
  }


  @Test
  void emptyOrderAreRemovedFromBook() {

    var testOrder = new TestOrder(Side.Buy, 10, 50);
    var testOrder1 = new TestOrder(Side.Buy, 10, 50);
    var testOrder2 = new TestOrder(Side.Buy, 10, 50);
    var testOrder3 = new TestOrder(Side.Buy, 10, 50);

    var testOrder4 = new TestOrder(Side.Sell, 35, 50);

    addOrder(testOrder);
    addOrder(testOrder1);
    addOrder(testOrder2);
    addOrder(testOrder3);
    addOrder(testOrder4);

    assertEquals(1,orderBook.buyOrders.size());
    assertEquals(0,orderBook.sellOrders.size());


  }


  @Test
  void addOrder() {
    var testOrder = new TestOrder(Side.Buy, 10, 10);

    addOrder(testOrder);

    ComparisonBook testBook = new ComparisonBook();
    testBook.addOrder( testOrder);

    testBook.testEquals(mdEntries);
  }

  @Test
  void addMultipleOrders() {
    var testOrder = new TestOrder(Side.Buy, 10, 10);
    var testOrder1 = new TestOrder(Side.Buy, 10, 9);

    var testOrder2 = new TestOrder(Side.Sell, 10, 12);
    var testOrder3 = new TestOrder(Side.Sell, 10, 14);

    addOrder(testOrder);
    addOrder(testOrder1);
    addOrder(testOrder2);
    addOrder(testOrder3);


    ComparisonBook testBook = new ComparisonBook();
    testBook.addOrder( testOrder);
    testBook.addOrder( testOrder1);
    testBook.addOrder( testOrder2);
    testBook.addOrder( testOrder3);


    testBook.testEquals(mdEntries);
  }


  @Test
  void addMultipleOrdersInRandomPriceOrder() {
    var testOrder = new TestOrder(Side.Buy, 10, 8);
    var testOrder1 = new TestOrder(Side.Buy, 10, 7);
    var testOrder2 = new TestOrder(Side.Buy, 10, 10);
    var testOrder3 = new TestOrder(Side.Buy, 10, 6);
    var testOrder4 = new TestOrder(Side.Buy, 5, 7);
    var testOrder5 = new TestOrder(Side.Buy, 10, 9);
    var testOrder6 = new TestOrder(Side.Buy, 10, 10);

    var testOrder7 = new TestOrder(Side.Sell, 10, 12);
    var testOrder8 = new TestOrder(Side.Sell, 10, 14);
    var testOrder9 = new TestOrder(Side.Sell, 10, 12);
    var testOrder10 = new TestOrder(Side.Sell, 10, 18);
    var testOrder11 = new TestOrder(Side.Sell, 5, 14);
    var testOrder12 = new TestOrder(Side.Sell, 10, 14);
    var testOrder13 = new TestOrder(Side.Sell, 10, 12);


    addOrder(testOrder);
    addOrder(testOrder1);
    addOrder(testOrder2);
    addOrder(testOrder3);
    addOrder(testOrder4);
    addOrder(testOrder5);
    addOrder(testOrder6);
    addOrder(testOrder7);
    addOrder(testOrder8);
    addOrder(testOrder9);
    addOrder(testOrder10);
    addOrder(testOrder11);
    addOrder(testOrder12);
    addOrder(testOrder13);


    ComparisonBook testBook = new ComparisonBook();
    testBook.addOrder( testOrder);
    testBook.addOrder( testOrder1);
    testBook.addOrder( testOrder2);
    testBook.addOrder( testOrder3);
    testBook.addOrder( testOrder4);
    testBook.addOrder( testOrder5);
    testBook.addOrder( testOrder6);
    testBook.addOrder( testOrder7);
    testBook.addOrder( testOrder8);
    testBook.addOrder( testOrder9);
    testBook.addOrder( testOrder10);
    testBook.addOrder( testOrder11);
    testBook.addOrder( testOrder12);
    testBook.addOrder( testOrder13);

    for( int i=0; i < orderBook.buyOrders.size() - 1;i++) {
      assertTrue(
          orderBook.buyOrders.get(i).getPrice().compareTo(orderBook.buyOrders.get(i + 1).getPrice()) != -1);
    }

    for( int i=0; i < orderBook.sellOrders.size() - 1;i++) {
      assertTrue(orderBook.sellOrders.get(i).getPrice()
          .compareTo(orderBook.sellOrders.get(i + 1).getPrice()) != 1);
    }

    testBook.testEquals(mdEntries);
  }




  @Test
  void executeTrade() {
    var buyOrder = new TestOrder(Side.Buy, 15, 10);
    var sellOrder = new TestOrder(Side.Sell, 8, 10);

    addOrder(buyOrder);
    addOrder(sellOrder);

    ComparisonBook compBook = new ComparisonBook();
    compBook.addOrder( Side.Buy, 7, 10);

    compBook.testEquals(mdEntries);

    assertEquals(2, trades.size() );
    Trade buyTrade = trades.get(0);
    Trade sellTrade = trades.get(1);

    assertEquals(8, buyTrade.getQuantity());
    assertEquals(new BigDecimal(10), buyTrade.getPrice());
    assertEquals( buyOrder.clOrderId, buyTrade.getClOrderId());

    assertEquals(8, sellTrade.getQuantity());
    assertEquals(new BigDecimal(10), sellTrade.getPrice());
    assertEquals( sellOrder.clOrderId, sellTrade.getClOrderId());
  }

  @Test
  void executeTradeOfSameSizeOrders() {
    var buyOrder = new TestOrder(Side.Buy, 8, 10);
    var sellOrder = new TestOrder(Side.Sell, 8, 10);

    addOrder(buyOrder);
    addOrder(sellOrder);

    ComparisonBook compBook = new ComparisonBook();

    compBook.testEquals(mdEntries);

    assertEquals(2, trades.size() );
    Trade buyTrade = trades.get(0);
    Trade sellTrade = trades.get(1);

    assertEquals(8, buyTrade.getQuantity());
    assertEquals(new BigDecimal(10), buyTrade.getPrice());
    assertEquals( buyOrder.clOrderId, buyTrade.getClOrderId());
    assertEquals( 0, buyTrade.getLeavesQty());
    assertEquals( 8, buyTrade.getCumQty());

    assertEquals(8, sellTrade.getQuantity());
    assertEquals(new BigDecimal(10), sellTrade.getPrice());
    assertEquals( sellOrder.clOrderId, sellTrade.getClOrderId());
    assertEquals( 0, sellTrade.getLeavesQty());
    assertEquals( 8, sellTrade.getCumQty());

    var tradeEntry = mdEntries.stream().filter(md->md.getMdEntryType()== MdEntryType.Trade).findFirst().get();

    assertEquals(8, tradeEntry.getQuantity());
    assertEquals(new BigDecimal(10), tradeEntry.getPrice());
    assertEquals( buyTrade.getTradeId(), tradeEntry.getId());

    var totalTraded = mdEntries.stream().filter(md->md.getMdEntryType()== MdEntryType.TradeVolume).findFirst().get();

    assertEquals(8, totalTraded.getQuantity());

  }




  @Test
  void quantityOnRemoveEntryIsCorrect() {
    var buyOrder = new TestOrder(Side.Buy, 15, 10);
    var sellOrder = new TestOrder(Side.Sell, 8, 10);

    addOrder(sellOrder);
    addOrder(buyOrder);


    ComparisonBook compBook = new ComparisonBook();
    compBook.addOrder( Side.Buy, 7, 10);

    compBook.testEquals(mdEntries);

    MDEntry entry = mdEntries.stream().filter(e->e.getMdUpdateAction() == MdUpdateActionType.Remove).findFirst().get();

    assertEquals(8, entry.getQuantity());
  }



  @Test
  void tradeMultipleOpppositeSideOrderWithSingleOrder() {
    var buyOrder = new TestOrder(Side.Buy, 22, 12);
    var sellOrder = new TestOrder(Side.Sell, 8, 10);
    var sellOrder1 = new TestOrder(Side.Sell, 8, 11);
    var sellOrder2 = new TestOrder(Side.Sell, 8, 12);
    var sellOrder3 = new TestOrder(Side.Sell, 8, 13);

    addOrder(sellOrder);
    addOrder(sellOrder1);
    addOrder(sellOrder2);
    addOrder(sellOrder3);

    addOrder(buyOrder);

    ComparisonBook compBook = new ComparisonBook();
    compBook.addOrder( Side.Sell, 2, 12);
    compBook.addOrder( Side.Sell, 8, 13);

    compBook.testEquals(mdEntries);

    assertEquals(6, trades.size() );

    List<Trade> buyTrades = trades.stream().filter(t->t.getOrderSide() == Side.Buy).collect(
        Collectors.toList());

    assertEquals(3, buyTrades.size() );

    List<Trade> sellTrades = trades.stream().filter(t->t.getOrderSide() == Side.Sell).collect(
        Collectors.toList());

    assertEquals(3, sellTrades.size() );

    Trade buyTrade = buyTrades.get(0);
    Trade buyTrade1 = buyTrades.get(1);
    Trade buyTrade2 = buyTrades.get(2);

    assertEquals(8, buyTrade.getQuantity());
    assertEquals(new BigDecimal(10), buyTrade.getPrice());
    assertEquals( buyOrder.clOrderId, buyTrade.getClOrderId());
    assertEquals( 14, buyTrade.getLeavesQty());
    assertEquals( 8, buyTrade.getCumQty());

    assertEquals(8, buyTrade1.getQuantity());
    assertEquals(new BigDecimal(11), buyTrade1.getPrice());
    assertEquals( buyOrder.clOrderId, buyTrade1.getClOrderId());
    assertEquals( 6, buyTrade1.getLeavesQty());
    assertEquals( 16, buyTrade1.getCumQty());

    assertEquals(6, buyTrade2.getQuantity());
    assertEquals(new BigDecimal(12), buyTrade2.getPrice());
    assertEquals( buyOrder.clOrderId, buyTrade2.getClOrderId());
    assertEquals( 0, buyTrade2.getLeavesQty());
    assertEquals( 22, buyTrade2.getCumQty());

    Trade sellTrade = sellTrades.get(0);
    Trade sellTrade1 = sellTrades.get(1);
    Trade sellTrade2 = sellTrades.get(2);

    assertEquals(8, sellTrade.getQuantity());
    assertEquals(new BigDecimal(10), sellTrade.getPrice());
    assertEquals( sellOrder.clOrderId, sellTrade.getClOrderId());
    assertEquals( 0, sellTrade.getLeavesQty());
    assertEquals( 8, sellTrade.getCumQty());

    assertEquals(8, sellTrade1.getQuantity());
    assertEquals(new BigDecimal(11), sellTrade1.getPrice());
    assertEquals( sellOrder1.clOrderId, sellTrade1.getClOrderId());
    assertEquals( 0, sellTrade1.getLeavesQty());
    assertEquals( 8, sellTrade1.getCumQty());

    assertEquals(6, sellTrade2.getQuantity());
    assertEquals(new BigDecimal(12), sellTrade2.getPrice());
    assertEquals( sellOrder2.clOrderId, sellTrade2.getClOrderId());
    assertEquals( 2, sellTrade2.getLeavesQty());
    assertEquals( 6, sellTrade2.getCumQty());

    var tradeEntries = mdEntries.stream().filter(md->md.getMdEntryType()== MdEntryType.Trade).collect(Collectors.toList());
    assertEquals(3, tradeEntries.size());
    var trade = tradeEntries.get(0);
    var trade1 = tradeEntries.get(1);
    var trade2 = tradeEntries.get(2);

    assertEquals(8, trade.getQuantity());
    assertEquals(new BigDecimal(10), trade.getPrice());
    assertEquals( sellTrade.getTradeId(), trade.getId());

    assertEquals(8, trade1.getQuantity());
    assertEquals(new BigDecimal(11), trade1.getPrice());
    assertEquals( sellTrade1.getTradeId(), trade1.getId());

    assertEquals(6, trade2.getQuantity());
    assertEquals(new BigDecimal(12), trade2.getPrice());
    assertEquals( sellTrade2.getTradeId(), trade2.getId());


    var totalTraded = mdEntries.stream().filter(md->md.getMdEntryType()== MdEntryType.TradeVolume).findFirst().get();

    assertEquals(22, totalTraded.getQuantity());



  }

  @Test
  void modifyOrderQuantity() throws Exception {

    var testOrder = new TestOrder(Side.Buy, 10, 10);
    var testOrder1 = new TestOrder(Side.Buy, 5, 9);
    var testOrder2 = new TestOrder(Side.Buy, 8, 9);
    var testOrder3 = new TestOrder(Side.Buy, 12, 9);


    addOrder(testOrder);
    addOrder(testOrder1);
    String orderId = addOrder(testOrder2);
    addOrder(testOrder3);

    orderBook.modifyOrder(orderId, new BigDecimal(9), 6);

    ComparisonBook testBook = new ComparisonBook();
    testBook.addOrder( testOrder);
    testBook.addOrder( testOrder1);
    testBook.addOrder( new TestOrder(Side.Buy, 6, 9));
    testBook.addOrder( testOrder3);

    testBook.testEquals(mdEntries);

  }

  @Test
  void modifyOrderToTrade() throws Exception {

    var testOrder = new TestOrder(Side.Buy, 10, 10);
    var testOrder1 = new TestOrder(Side.Buy, 10, 9);

    var testOrder2 = new TestOrder(Side.Sell, 8, 12);
    var testOrder3 = new TestOrder(Side.Sell, 10, 14);

    addOrder(testOrder);
    addOrder(testOrder1);
    String sellOrderToModifyId = addOrder(testOrder2);
    addOrder(testOrder3);

    orderBook.modifyOrder(sellOrderToModifyId, new BigDecimal(10), 8);


    ComparisonBook testBook = new ComparisonBook();
    testBook.addOrder( new TestOrder(Side.Buy, 2, 10));
    testBook.addOrder( testOrder1);
    testBook.addOrder( testOrder3);


    testBook.testEquals(mdEntries);

    assertEquals(2, trades.size());

    Trade buyTrade = trades.stream().filter(t->t.getOrderSide() == Side.Buy).findFirst().get();
    Trade sellTrade = trades.stream().filter(t->t.getOrderSide() == Side.Sell).findFirst().get();

    assertEquals(8, sellTrade.getQuantity());
    assertEquals(new BigDecimal(10), sellTrade.getPrice());

    assertEquals(8, buyTrade.getQuantity());
    assertEquals(new BigDecimal(10), buyTrade.getPrice());

  }

  @Test
  void testModifyToTradeMultipleOrders() throws Exception {

    var testOrder = new TestOrder(Side.Buy, 10, 10);
    var testOrder1 = new TestOrder(Side.Buy, 10, 9);
    var testOrder2 = new TestOrder(Side.Buy, 10, 9);

    var testOrder3 = new TestOrder(Side.Sell, 8, 12);
    var testOrder4 = new TestOrder(Side.Sell, 10, 14);
    var testOrder5 = new TestOrder(Side.Sell, 10, 15);

    String idOfOrderToModify = addOrder(testOrder);
    addOrder(testOrder1);
    addOrder(testOrder2);
    addOrder(testOrder3);
    addOrder(testOrder4);
    addOrder(testOrder5);

    orderBook.modifyOrder(idOfOrderToModify, new BigDecimal(15), 22);


    ComparisonBook testBook = new ComparisonBook();

    testBook.addOrder( testOrder1);
    testBook.addOrder( testOrder2);
    testBook.addOrder( new TestOrder(Side.Sell, 6, 15));

    testBook.testEquals(mdEntries);

    assertEquals(6, trades.size());

    List<Trade> buyTrades = trades.stream().filter(t->t.getOrderSide() == Side.Buy).collect(
        Collectors.toList());

    List<Trade> sellTrades = trades.stream().filter(t->t.getOrderSide() == Side.Sell).collect(
        Collectors.toList());

    tradeMatches(buyTrades.get(0), 8,12);
    tradeMatches(sellTrades.get(0), 8,12);

    tradeMatches(buyTrades.get(1), 10,14);
    tradeMatches(sellTrades.get(1), 10,14);

    tradeMatches(buyTrades.get(2), 4,15);
    tradeMatches(sellTrades.get(2), 4,15);

  }

  void tradeMatches(Trade trade, int qnt, int price) {
    assertEquals(qnt, trade.getQuantity());
    assertEquals(new BigDecimal(price), trade.getPrice());

  }


  @Test
  void tradingOrderWithMultipleSamePriceOrdersGivesCorrectTrades() {


    var testOrder = new TestOrder(Side.Buy, 12, 10);
    var testOrder1 = new TestOrder(Side.Buy, 10, 10);
    var testOrder2 = new TestOrder(Side.Buy, 10, 9);

    var testOrder3 = new TestOrder(Side.Sell, 8, 12);
    var testOrder4 = new TestOrder(Side.Sell, 10, 14);
    var testOrder5 = new TestOrder(Side.Sell, 10, 15);

    addOrder(testOrder);
    addOrder(testOrder1);
    addOrder(testOrder2);
    addOrder(testOrder3);
    addOrder(testOrder4);
    addOrder(testOrder5);

    addOrder(new TestOrder(Side.Sell, 8, 10));

    assertEquals(2, trades.size());

    List<Trade> buyTrades = trades.stream().filter(t->t.getOrderSide() == Side.Buy).collect(
        Collectors.toList());

    List<Trade> sellTrades = trades.stream().filter(t->t.getOrderSide() == Side.Sell).collect(
        Collectors.toList());

    tradeMatches(buyTrades.get(0), 8,10);
    tradeMatches(sellTrades.get(0), 8,10);



  }

  @Test
  void modifiedOrderKeepsBookPosition() throws Exception {

    var testOrder = new TestOrder(Side.Buy, 10, 10);
    var testOrder1 = new TestOrder(Side.Buy, 10, 10);
    var testOrder2 = new TestOrder(Side.Buy, 10, 9);

    var testOrder3 = new TestOrder(Side.Sell, 8, 12);
    var testOrder4 = new TestOrder(Side.Sell, 10, 14);
    var testOrder5 = new TestOrder(Side.Sell, 10, 15);

    String idOfOrderToModify = addOrder(testOrder);
    addOrder(testOrder1);
    addOrder(testOrder2);
    addOrder(testOrder3);
    addOrder(testOrder4);
    addOrder(testOrder5);

    orderBook.modifyOrder(idOfOrderToModify, new BigDecimal(10), 8);

    var testOrder6 = new TestOrder(Side.Sell, 6, 10);
    addOrder(testOrder6);


    ComparisonBook testBook = new ComparisonBook();
    testBook.addOrder( new TestOrder(Side.Buy, 2, 10));
    testBook.addOrder( testOrder1);
    testBook.addOrder( testOrder2);
    testBook.addOrder( testOrder3);
    testBook.addOrder( testOrder4);
    testBook.addOrder( testOrder5);


    testBook.testEquals(mdEntries);

    assertEquals(2, trades.size());

    List<Trade> buyTrades = trades.stream().filter(t->t.getOrderSide() == Side.Buy).collect(
        Collectors.toList());

    List<Trade> sellTrades = trades.stream().filter(t->t.getOrderSide() == Side.Sell).collect(
        Collectors.toList());

    tradeMatches(buyTrades.get(0), 6,10);
    tradeMatches(sellTrades.get(0), 6,10);

    assertEquals(idOfOrderToModify,  buyTrades.get(0).getOrderId());

  }

  @Test
  void deleteOrder() {

    var testOrder = new TestOrder(Side.Buy, 10, 10);
    var testOrder1 = new TestOrder(Side.Buy, 10, 9);

    var testOrder2 = new TestOrder(Side.Sell, 10, 12);
    var testOrder3 = new TestOrder(Side.Sell, 10, 14);
    var testOrder4 = new TestOrder(Side.Sell, 10, 15);

    addOrder(testOrder);
    addOrder(testOrder1);
    addOrder(testOrder2);
    String orderId = addOrder(testOrder3);
    addOrder(testOrder4);


    try {
      orderBook.deleteOrder(orderId);
    } catch (OrderDeletionException e) {
      fail(e.getMessage());
    }

    ComparisonBook testBook = new ComparisonBook();
    testBook.addOrder( testOrder);
    testBook.addOrder( testOrder1);
    testBook.addOrder( testOrder2);
    testBook.addOrder( testOrder4);

    testBook.testEquals(mdEntries);

  }

  String addOrder( TestOrder order) {
    return orderBook.addOrder(order.side, order.qty, order.price, order.clOrderId);
  }

}
