package com.ettech.fixmarketsimulator.marketdata;

import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.ettech.fixmarketsimulator.exchange.OrderBook;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.SimpleChannelInboundHandler;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import javax.inject.Inject;
import org.fixprotocol.components.InstrmtMDReqGrp;
import org.fixprotocol.components.MarketData.MarketDataRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class MarketDataServerHandler extends SimpleChannelInboundHandler<MarketDataRequest>   {

  Logger log = LoggerFactory.getLogger(MarketDataServerHandler.class);


  List<MarketDataSubscription> subscriptions = new ArrayList<>();
  Set<String> symbols = new HashSet<>();

  Exchange exchange;

  MarketDataServerHandler(Exchange exchange) {
    this.exchange = exchange;
  }


  @Override
  protected void channelRead0(ChannelHandlerContext ctx, MarketDataRequest msg)
      throws Exception {

    log.info("Received subscription request", msg);

    for( InstrmtMDReqGrp mdReqGrp : msg.getInstrmtMdReqGrpList() ) {
      String symbol = mdReqGrp.getInstrument().getSymbol();

      log.info("Received subscriprtion request for symbol {}", symbol);

      if( symbols.contains(symbol)) {
        log.warn("Ignoring subscription as already subscribed to " + symbol);
        return;
      }

      symbols.add(symbol);

      OrderBook book = exchange.getOrderBook(symbol);
      subscriptions.add(new MarketDataSubscription(ctx, book, msg.getMdReqId()));
    }
  }

  @Override
  public void channelReadComplete(ChannelHandlerContext ctx) {
    ctx.flush();
  }

  @Override
  public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) {
    log.warn("Exception caught:", cause);

    subscriptions.forEach(s->  {
      try {
        s.close();
      } catch (Throwable t) {
        log.error("Failed to close subscription", t);
      }
    });

    ctx.close();
  }


}
