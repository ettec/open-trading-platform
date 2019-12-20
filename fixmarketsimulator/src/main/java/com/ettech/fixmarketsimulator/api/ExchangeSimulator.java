package com.ettech.fixmarketsimulator.api;

import com.ettech.fixmarketsimulator.App;
import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.exchange.OrderBook;


import com.ettech.fixmarketsimulator.exchange.Side;
import java.math.BigDecimal;
import javax.inject.Inject;
import javax.ws.rs.*;
import javax.ws.rs.core.MediaType;

@Path("/api/exchange-simulator")
public class ExchangeSimulator {

    Exchange exchange;

    public ExchangeSimulator() {
        this.exchange = App.exchangeSimulatorInjector.getInstance(Exchange.class);
    }

    @PUT
    @Path("order")
    @Produces(MediaType.TEXT_PLAIN)
    public String addOrder(@QueryParam("instrumentId") String instrumentId,
                           @QueryParam("Quantity") int qty,
                           @QueryParam("Price") BigDecimal price,
                           @QueryParam("Side") boolean sideBoolean) {

        OrderBook orderBook = exchange.getOrderBook(instrumentId);


        Side side = sideBoolean ? Side.Buy : Side.Sell;

        return orderBook.addOrder(side, qty, price, "");

    }

}
