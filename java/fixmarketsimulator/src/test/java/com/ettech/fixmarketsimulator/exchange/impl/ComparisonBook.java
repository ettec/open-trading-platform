package com.ettech.fixmarketsimulator.exchange.impl;

import com.ettech.fixmarketsimulator.exchange.MDEntry;
import com.ettech.fixmarketsimulator.exchange.Side;
import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;
import java.util.UUID;
import org.opentest4j.AssertionFailedError;

public class ComparisonBook {

  int addedCount = 0;

  List<TestBookOrder> bids = new ArrayList<>();
  List<TestBookOrder> offers = new ArrayList<>();


  void addOrder( TestOrder testOrder) {
    addOrder(testOrder.side, testOrder.qty, testOrder.price);
  }

  void addOrder( Side side, int qty, int price) {
    addOrder(side, qty, new BigDecimal(price));
  }

  void addOrder( Side side, int qty, BigDecimal price) {

    if( side == Side.Buy) {
      bids.add( new TestBookOrder( side, qty, price, addedCount++, UUID.randomUUID().toString()));
    } else {
      offers.add( new TestBookOrder( side, qty, price, addedCount++, UUID.randomUUID().toString()));
    }


    sortBids(bids);
    sortOffers(offers);

  }


  public void testEquals(List<MDEntry> entries) {

    List<TestBookOrder> bidsFromMdEntries = new ArrayList<>();
    List<TestBookOrder> offersFromMdEntries = new ArrayList<>();

    entries.forEach(e-> {
      List<TestBookOrder> orders = e.getSide() == Side.Buy ? bidsFromMdEntries : offersFromMdEntries;


      switch( e.getMdEntryType()) {
        case Add:
          orders.add(new TestBookOrder( e.getSide(), e.getQuantity(), e.getPrice(), addedCount++,
              e.getOrderId()));
        break;
        case Modify:
          TestBookOrder originalOrder = orders.stream().filter(o->o.orderId.equals(e.getOrderId())).findFirst().get();
          originalOrder.price = e.getPrice();
          originalOrder.qty = e.getQuantity();
          break;
        case Remove:
          orders.remove(new TestBookOrder( e.getSide(), e.getQuantity(), e.getPrice(), addedCount++,
              e.getOrderId()));
          break;

      }

    });

    sortBids(bidsFromMdEntries);
    sortOffers(offersFromMdEntries);

    areEqual(bids, bidsFromMdEntries);
    areEqual(offers, offersFromMdEntries);

  }


  private void areEqual(List<TestBookOrder> a, List<TestBookOrder> b) {
    if( a.size() != b.size() ) {
      throw new AssertionFailedError("Order books size mismatch", b.size(), a.size());
    }

    for( int i=0; i < a.size(); i++) {
      if( !testOrdersQuantityPriceAndSideEquals(a.get(i), b.get(i))) {
        throw new AssertionFailedError(String.format("Orders in books at index %s do not match. %s %s @ %s   versus %s %s @ %s  ", i,
            a.get(i).side, a.get(i).qty, a.get(i).price, b.get(i).side, b.get(i).qty, b.get(i).price ));
      }
    }

  }

  private boolean testOrdersQuantityPriceAndSideEquals(TestBookOrder a, TestBookOrder b) {
    if (b == a) {
      return true;
    }
    if (b == null || a.getClass() != b.getClass()) {
      return false;
    }
    TestBookOrder testOrder = (TestBookOrder) b;
    return a.side == b.side &&
        a.price.equals(b.price) &&
        a.qty == b.qty;
  }

  private void sortOffers(List<TestBookOrder> offersFromMdEntries) {
    offersFromMdEntries.sort((a,b)-> {

      int priceComparison = a.price.compareTo(b.price);

      if( priceComparison == 0) {
        return a.ordinal - b.ordinal;
      } else {
        return priceComparison;
      }
    }   );
  }

  private void sortBids(List<TestBookOrder> bidsFromMdEntries) {
    bidsFromMdEntries.sort((a,b)-> {

      int priceComparison = a.price.compareTo(b.price);

      if( priceComparison == 0) {
        return a.ordinal - b.ordinal;
      } else {
        return -priceComparison;
      }


    }   );
  }


  static class TestBookOrder {


    public TestBookOrder( Side side, double qty, BigDecimal price, int addedCount,
        String orderId) {

      this.side = side;
      this.qty = qty;
      this.price = price;
      this.ordinal = addedCount;
      this.orderId = orderId;
    }


    Side side;
    double qty;
    BigDecimal price;
    int ordinal;
    String orderId;


    @Override
    public boolean equals(Object o) {
      if (this == o) {
        return true;
      }
      if (o == null || getClass() != o.getClass()) {
        return false;
      }
      TestBookOrder testOrder = (TestBookOrder) o;
      return side == testOrder.side &&
          price.equals(testOrder.price) &&
          orderId.equals(testOrder.orderId);
    }

    @Override
    public int hashCode() {
      return Objects.hash(side, price, orderId);
    }
  }


}
