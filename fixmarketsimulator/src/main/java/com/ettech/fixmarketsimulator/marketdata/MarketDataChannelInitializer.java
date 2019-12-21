package com.ettech.fixmarketsimulator.marketdata;

import com.ettech.fixmarketsimulator.exchange.Exchange;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.ChannelPipeline;
import io.netty.channel.socket.SocketChannel;
import io.netty.handler.codec.protobuf.ProtobufDecoder;
import io.netty.handler.codec.protobuf.ProtobufEncoder;
import io.netty.handler.codec.protobuf.ProtobufVarint32FrameDecoder;
import io.netty.handler.codec.protobuf.ProtobufVarint32LengthFieldPrepender;
import org.fixprotocol.components.MarketData;

public class MarketDataChannelInitializer extends ChannelInitializer<SocketChannel> {

  Exchange exchange;

  MarketDataChannelInitializer(Exchange exchange) {
    this.exchange = exchange;
  }

  @Override
  protected void initChannel(SocketChannel ch) throws Exception {
    ChannelPipeline p = ch.pipeline();
    p.addLast(new ProtobufVarint32FrameDecoder());
    p.addLast(new ProtobufDecoder(MarketData.MarketDataRequest.getDefaultInstance()));

    p.addLast(new ProtobufVarint32LengthFieldPrepender());
    p.addLast(new ProtobufEncoder());

    p.addLast(new MarketDataServerHandler(exchange));

  }

}
