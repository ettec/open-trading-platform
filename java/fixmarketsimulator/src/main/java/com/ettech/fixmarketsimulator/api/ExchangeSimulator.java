package com.ettech.fixmarketsimulator.api;

import com.ettech.fixmarketsimulator.App;
import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.exchange.OrderBook;


import com.ettech.fixmarketsimulator.exchange.Side;
import java.math.BigDecimal;
import java.util.Arrays;
import javax.inject.Inject;
import javax.ws.rs.*;
import javax.ws.rs.core.MediaType;

@Path("/api/exchange-simulator")
public class ExchangeSimulator {

    Exchange exchange;

    public ExchangeSimulator() {
        this.exchange = App.exchangeSimulatorInjector.getInstance(Exchange.class);
    }

    @GET
    @Path("book")
    @Produces(MediaType.TEXT_PLAIN)
    public String getBook(@QueryParam("symbol") String symbol) {

        OrderBook orderBook = exchange.getOrderBook(symbol);

        var buys = orderBook.getBuyOrders();
        var sells = orderBook.getSellOrders();

        var result = "Bids\n";
        for(var bid : buys) {
            result+=bid.getRemainingQty() + "@" + bid.getPrice() + "\t" + bid.getClOrdId() +"\n";
        }
        result += "\nAsks\n";

        for(var sell : sells) {
            result += sell.getRemainingQty() + "@" + sell.getPrice() + "\t" + sell.getClOrdId() + "\n";
        }

        return result;

    }

    @PUT
    @Path("order")
    @Produces(MediaType.TEXT_PLAIN)
    public String addOrder(@QueryParam("symbol") String symbol,
                           @QueryParam("Quantity") int qty,
                           @QueryParam("Price") BigDecimal price,
                           @QueryParam("Side") boolean sideBoolean) {

        OrderBook orderBook = exchange.getOrderBook(symbol);


        Side side = sideBoolean ? Side.Buy : Side.Sell;

        return orderBook.addOrder(side, qty, price, "");

    }

}
