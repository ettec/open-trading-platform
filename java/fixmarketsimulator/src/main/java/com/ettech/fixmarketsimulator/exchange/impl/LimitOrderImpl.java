package com.ettech.fixmarketsimulator.exchange.impl;

import com.ettech.fixmarketsimulator.exchange.Side;
import java.math.BigDecimal;
import java.util.Objects;

public class LimitOrderImpl implements com.ettech.fixmarketsimulator.exchange.Order {

    public LimitOrderImpl(double quantity, BigDecimal price, String clOrdId, Side side,
        String orderId) {
        this.quantity = quantity;
        this.remainingQty = quantity;
        this.price = price;
        this.clOrdId = clOrdId;
        this.side = side;
        this.orderId = orderId;
    }

    private double quantity;
    private double remainingQty;
    private BigDecimal price;
    private String clOrdId;
    private Side side;
    private String orderId;

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        LimitOrderImpl that = (LimitOrderImpl) o;
        return orderId.equals(that.orderId);
    }

    @Override
    public int hashCode() {
        return Objects.hash(orderId);
    }

    @Override
    public double getQuantity() {
        return quantity;
    }

    public void setQuantity(double quantity) {
        this.quantity = quantity;
    }

    @Override
    public double getRemainingQty() {
        return remainingQty;
    }

    public void setRemainingQty(double remainingQty) {
        this.remainingQty = remainingQty;
    }

    @Override
    public BigDecimal getPrice() {
        return price;
    }

    public void setPrice(BigDecimal price) {
        this.price = price;
    }

    @Override
    public String getClOrdId() {
        return clOrdId;
    }

    public void setClOrdId(String clOrdId) {
        this.clOrdId = clOrdId;
    }

    @Override
    public Side getSide() {
        return side;
    }

    public void setSide(Side side) {
        this.side = side;
    }

    @Override
    public String getOrderId() {
        return orderId;
    }

    public void setOrderId(String orderId) {
        this.orderId = orderId;
    }
}
