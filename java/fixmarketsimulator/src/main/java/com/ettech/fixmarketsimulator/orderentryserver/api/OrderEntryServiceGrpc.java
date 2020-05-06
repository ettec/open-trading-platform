package com.ettech.fixmarketsimulator.orderentryserver.api;

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
    comments = "Source: orderentryapi.proto")
public final class OrderEntryServiceGrpc {

  private OrderEntryServiceGrpc() {}

  public static final String SERVICE_NAME = "orderentryapi.OrderEntryService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.NewOrderParams,
      com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId> getSubmitNewOrderMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SubmitNewOrder",
      requestType = com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.NewOrderParams.class,
      responseType = com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.NewOrderParams,
      com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId> getSubmitNewOrderMethod() {
    io.grpc.MethodDescriptor<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.NewOrderParams, com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId> getSubmitNewOrderMethod;
    if ((getSubmitNewOrderMethod = OrderEntryServiceGrpc.getSubmitNewOrderMethod) == null) {
      synchronized (OrderEntryServiceGrpc.class) {
        if ((getSubmitNewOrderMethod = OrderEntryServiceGrpc.getSubmitNewOrderMethod) == null) {
          OrderEntryServiceGrpc.getSubmitNewOrderMethod = getSubmitNewOrderMethod =
              io.grpc.MethodDescriptor.<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.NewOrderParams, com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SubmitNewOrder"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.NewOrderParams.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId.getDefaultInstance()))
              .setSchemaDescriptor(new OrderEntryServiceMethodDescriptorSupplier("SubmitNewOrder"))
              .build();
        }
      }
    }
    return getSubmitNewOrderMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId,
      com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.Empty> getCancelOrderMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CancelOrder",
      requestType = com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId.class,
      responseType = com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.Empty.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId,
      com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.Empty> getCancelOrderMethod() {
    io.grpc.MethodDescriptor<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId, com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.Empty> getCancelOrderMethod;
    if ((getCancelOrderMethod = OrderEntryServiceGrpc.getCancelOrderMethod) == null) {
      synchronized (OrderEntryServiceGrpc.class) {
        if ((getCancelOrderMethod = OrderEntryServiceGrpc.getCancelOrderMethod) == null) {
          OrderEntryServiceGrpc.getCancelOrderMethod = getCancelOrderMethod =
              io.grpc.MethodDescriptor.<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId, com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.Empty>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CancelOrder"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.Empty.getDefaultInstance()))
              .setSchemaDescriptor(new OrderEntryServiceMethodDescriptorSupplier("CancelOrder"))
              .build();
        }
      }
    }
    return getCancelOrderMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static OrderEntryServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<OrderEntryServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<OrderEntryServiceStub>() {
        @java.lang.Override
        public OrderEntryServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new OrderEntryServiceStub(channel, callOptions);
        }
      };
    return OrderEntryServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static OrderEntryServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<OrderEntryServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<OrderEntryServiceBlockingStub>() {
        @java.lang.Override
        public OrderEntryServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new OrderEntryServiceBlockingStub(channel, callOptions);
        }
      };
    return OrderEntryServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static OrderEntryServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<OrderEntryServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<OrderEntryServiceFutureStub>() {
        @java.lang.Override
        public OrderEntryServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new OrderEntryServiceFutureStub(channel, callOptions);
        }
      };
    return OrderEntryServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class OrderEntryServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void submitNewOrder(com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.NewOrderParams request,
        io.grpc.stub.StreamObserver<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId> responseObserver) {
      asyncUnimplementedUnaryCall(getSubmitNewOrderMethod(), responseObserver);
    }

    /**
     */
    public void cancelOrder(com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId request,
        io.grpc.stub.StreamObserver<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.Empty> responseObserver) {
      asyncUnimplementedUnaryCall(getCancelOrderMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getSubmitNewOrderMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.NewOrderParams,
                com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId>(
                  this, METHODID_SUBMIT_NEW_ORDER)))
          .addMethod(
            getCancelOrderMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId,
                com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.Empty>(
                  this, METHODID_CANCEL_ORDER)))
          .build();
    }
  }

  /**
   */
  public static final class OrderEntryServiceStub extends io.grpc.stub.AbstractAsyncStub<OrderEntryServiceStub> {
    private OrderEntryServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected OrderEntryServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new OrderEntryServiceStub(channel, callOptions);
    }

    /**
     */
    public void submitNewOrder(com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.NewOrderParams request,
        io.grpc.stub.StreamObserver<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getSubmitNewOrderMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void cancelOrder(com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId request,
        io.grpc.stub.StreamObserver<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.Empty> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getCancelOrderMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class OrderEntryServiceBlockingStub extends io.grpc.stub.AbstractBlockingStub<OrderEntryServiceBlockingStub> {
    private OrderEntryServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected OrderEntryServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new OrderEntryServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId submitNewOrder(com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.NewOrderParams request) {
      return blockingUnaryCall(
          getChannel(), getSubmitNewOrderMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.Empty cancelOrder(com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId request) {
      return blockingUnaryCall(
          getChannel(), getCancelOrderMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class OrderEntryServiceFutureStub extends io.grpc.stub.AbstractFutureStub<OrderEntryServiceFutureStub> {
    private OrderEntryServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected OrderEntryServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new OrderEntryServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId> submitNewOrder(
        com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.NewOrderParams request) {
      return futureUnaryCall(
          getChannel().newCall(getSubmitNewOrderMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.Empty> cancelOrder(
        com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId request) {
      return futureUnaryCall(
          getChannel().newCall(getCancelOrderMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_SUBMIT_NEW_ORDER = 0;
  private static final int METHODID_CANCEL_ORDER = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final OrderEntryServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(OrderEntryServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_SUBMIT_NEW_ORDER:
          serviceImpl.submitNewOrder((com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.NewOrderParams) request,
              (io.grpc.stub.StreamObserver<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId>) responseObserver);
          break;
        case METHODID_CANCEL_ORDER:
          serviceImpl.cancelOrder((com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.OrderId) request,
              (io.grpc.stub.StreamObserver<com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.Empty>) responseObserver);
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

  private static abstract class OrderEntryServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    OrderEntryServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.ettech.fixmarketsimulator.orderentryserver.api.Orderentryapi.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("OrderEntryService");
    }
  }

  private static final class OrderEntryServiceFileDescriptorSupplier
      extends OrderEntryServiceBaseDescriptorSupplier {
    OrderEntryServiceFileDescriptorSupplier() {}
  }

  private static final class OrderEntryServiceMethodDescriptorSupplier
      extends OrderEntryServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    OrderEntryServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (OrderEntryServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new OrderEntryServiceFileDescriptorSupplier())
              .addMethod(getSubmitNewOrderMethod())
              .addMethod(getCancelOrderMethod())
              .build();
        }
      }
    }
    return result;
  }
}
