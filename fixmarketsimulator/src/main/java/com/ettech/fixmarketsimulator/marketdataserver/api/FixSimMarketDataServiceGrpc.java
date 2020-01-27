package com.ettech.fixmarketsimulator.marketdataserver.api;

import static io.grpc.MethodDescriptor.generateFullMethodName;
import static io.grpc.stub.ClientCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ClientCalls.asyncClientStreamingCall;
import static io.grpc.stub.ClientCalls.asyncServerStreamingCall;
import static io.grpc.stub.ClientCalls.asyncUnaryCall;
import static io.grpc.stub.ClientCalls.blockingServerStreamingCall;
import static io.grpc.stub.ClientCalls.blockingUnaryCall;
import static io.grpc.stub.ClientCalls.futureUnaryCall;
import static io.grpc.stub.ServerCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ServerCalls.asyncClientStreamingCall;
import static io.grpc.stub.ServerCalls.asyncServerStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnaryCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.26.0)",
    comments = "Source: marketdataserver.proto")
public final class FixSimMarketDataServiceGrpc {

  private FixSimMarketDataServiceGrpc() {}

  public static final String SERVICE_NAME = "marketdataservice.FixSimMarketDataService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<org.fixprotocol.components.MarketData.MarketDataRequest,
      com.google.protobuf.Empty> getSubscribeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Subscribe",
      requestType = org.fixprotocol.components.MarketData.MarketDataRequest.class,
      responseType = com.google.protobuf.Empty.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.fixprotocol.components.MarketData.MarketDataRequest,
      com.google.protobuf.Empty> getSubscribeMethod() {
    io.grpc.MethodDescriptor<org.fixprotocol.components.MarketData.MarketDataRequest, com.google.protobuf.Empty> getSubscribeMethod;
    if ((getSubscribeMethod = FixSimMarketDataServiceGrpc.getSubscribeMethod) == null) {
      synchronized (FixSimMarketDataServiceGrpc.class) {
        if ((getSubscribeMethod = FixSimMarketDataServiceGrpc.getSubscribeMethod) == null) {
          FixSimMarketDataServiceGrpc.getSubscribeMethod = getSubscribeMethod =
              io.grpc.MethodDescriptor.<org.fixprotocol.components.MarketData.MarketDataRequest, com.google.protobuf.Empty>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Subscribe"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.fixprotocol.components.MarketData.MarketDataRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setSchemaDescriptor(new FixSimMarketDataServiceMethodDescriptorSupplier("Subscribe"))
              .build();
        }
      }
    }
    return getSubscribeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver.SubscribeRequest,
      org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh> getConnectMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Connect",
      requestType = com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver.SubscribeRequest.class,
      responseType = org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver.SubscribeRequest,
      org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh> getConnectMethod() {
    io.grpc.MethodDescriptor<com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver.SubscribeRequest, org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh> getConnectMethod;
    if ((getConnectMethod = FixSimMarketDataServiceGrpc.getConnectMethod) == null) {
      synchronized (FixSimMarketDataServiceGrpc.class) {
        if ((getConnectMethod = FixSimMarketDataServiceGrpc.getConnectMethod) == null) {
          FixSimMarketDataServiceGrpc.getConnectMethod = getConnectMethod =
              io.grpc.MethodDescriptor.<com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver.SubscribeRequest, org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Connect"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver.SubscribeRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh.getDefaultInstance()))
              .setSchemaDescriptor(new FixSimMarketDataServiceMethodDescriptorSupplier("Connect"))
              .build();
        }
      }
    }
    return getConnectMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static FixSimMarketDataServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FixSimMarketDataServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FixSimMarketDataServiceStub>() {
        @java.lang.Override
        public FixSimMarketDataServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FixSimMarketDataServiceStub(channel, callOptions);
        }
      };
    return FixSimMarketDataServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static FixSimMarketDataServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FixSimMarketDataServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FixSimMarketDataServiceBlockingStub>() {
        @java.lang.Override
        public FixSimMarketDataServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FixSimMarketDataServiceBlockingStub(channel, callOptions);
        }
      };
    return FixSimMarketDataServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static FixSimMarketDataServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FixSimMarketDataServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FixSimMarketDataServiceFutureStub>() {
        @java.lang.Override
        public FixSimMarketDataServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FixSimMarketDataServiceFutureStub(channel, callOptions);
        }
      };
    return FixSimMarketDataServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class FixSimMarketDataServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void subscribe(org.fixprotocol.components.MarketData.MarketDataRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      asyncUnimplementedUnaryCall(getSubscribeMethod(), responseObserver);
    }

    /**
     */
    public void connect(com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver.SubscribeRequest request,
        io.grpc.stub.StreamObserver<org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh> responseObserver) {
      asyncUnimplementedUnaryCall(getConnectMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getSubscribeMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.fixprotocol.components.MarketData.MarketDataRequest,
                com.google.protobuf.Empty>(
                  this, METHODID_SUBSCRIBE)))
          .addMethod(
            getConnectMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver.SubscribeRequest,
                org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh>(
                  this, METHODID_CONNECT)))
          .build();
    }
  }

  /**
   */
  public static final class FixSimMarketDataServiceStub extends io.grpc.stub.AbstractAsyncStub<FixSimMarketDataServiceStub> {
    private FixSimMarketDataServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FixSimMarketDataServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FixSimMarketDataServiceStub(channel, callOptions);
    }

    /**
     */
    public void subscribe(org.fixprotocol.components.MarketData.MarketDataRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getSubscribeMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void connect(com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver.SubscribeRequest request,
        io.grpc.stub.StreamObserver<org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getConnectMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class FixSimMarketDataServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<FixSimMarketDataServiceBlockingStub> {
    private FixSimMarketDataServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FixSimMarketDataServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FixSimMarketDataServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.google.protobuf.Empty subscribe(org.fixprotocol.components.MarketData.MarketDataRequest request) {
      return blockingUnaryCall(
          getChannel(), getSubscribeMethod(), getCallOptions(), request);
    }

    /**
     */
    public java.util.Iterator<org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh> connect(
        com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver.SubscribeRequest request) {
      return blockingServerStreamingCall(
          getChannel(), getConnectMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class FixSimMarketDataServiceFutureStub extends io.grpc.stub.AbstractFutureStub<FixSimMarketDataServiceFutureStub> {
    private FixSimMarketDataServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FixSimMarketDataServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FixSimMarketDataServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.google.protobuf.Empty> subscribe(
        org.fixprotocol.components.MarketData.MarketDataRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getSubscribeMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_SUBSCRIBE = 0;
  private static final int METHODID_CONNECT = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final FixSimMarketDataServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(FixSimMarketDataServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_SUBSCRIBE:
          serviceImpl.subscribe((org.fixprotocol.components.MarketData.MarketDataRequest) request,
              (io.grpc.stub.StreamObserver<com.google.protobuf.Empty>) responseObserver);
          break;
        case METHODID_CONNECT:
          serviceImpl.connect((com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver.SubscribeRequest) request,
              (io.grpc.stub.StreamObserver<org.fixprotocol.components.MarketData.MarketDataIncrementalRefresh>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class FixSimMarketDataServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    FixSimMarketDataServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.ettech.fixmarketsimulator.marketdataserver.api.Marketdataserver.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("FixSimMarketDataService");
    }
  }

  private static final class FixSimMarketDataServiceFileDescriptorSupplier
      extends FixSimMarketDataServiceBaseDescriptorSupplier {
    FixSimMarketDataServiceFileDescriptorSupplier() {}
  }

  private static final class FixSimMarketDataServiceMethodDescriptorSupplier
      extends FixSimMarketDataServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    FixSimMarketDataServiceMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (FixSimMarketDataServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new FixSimMarketDataServiceFileDescriptorSupplier())
              .addMethod(getSubscribeMethod())
              .addMethod(getConnectMethod())
              .build();
        }
      }
    }
    return result;
  }
}
