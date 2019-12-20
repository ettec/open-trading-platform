package com.ettech.fixmarketsimulator.marketdata;

import com.ettech.fixmarketsimulator.exchange.Exchange;
import com.google.inject.Inject;
import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;
import java.io.Closeable;
import java.io.IOException;

public class MarketDataServer implements Closeable {

  int port = 9888;

  EventLoopGroup serverGroup;
  EventLoopGroup workerGroup;

  @Inject
  public MarketDataServer(Exchange exchange) throws Exception {

    serverGroup = new NioEventLoopGroup(1);
    workerGroup = new NioEventLoopGroup();


      ServerBootstrap bootStrap = new ServerBootstrap();
      bootStrap.group(serverGroup, workerGroup)
          .channel(NioServerSocketChannel.class)
          .handler(new LoggingHandler(LogLevel.INFO))
          .childHandler(new MarketDataChannelInitializer(exchange));

      // Bind to port
      bootStrap.bind(port);




  }

  @Override
  public void close() throws IOException {
    serverGroup.shutdownGracefully();
    workerGroup.shutdownGracefully();
  }



}
