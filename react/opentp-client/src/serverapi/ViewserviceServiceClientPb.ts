/**
 * @fileoverview gRPC-Web generated client stub for viewservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as order_pb from './order_pb';
import * as viewservice_pb from './viewservice_pb';


export class ViewServiceClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: any; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'text';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname;
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodInfoSubscribeToOrdersWithRootOriginatorId = new grpcWeb.AbstractClientBase.MethodInfo(
    order_pb.Order,
    (request: viewservice_pb.SubscribeToOrdersWithRootOriginatorIdArgs) => {
      return request.serializeBinary();
    },
    order_pb.Order.deserializeBinary
  );

  subscribeToOrdersWithRootOriginatorId(
    request: viewservice_pb.SubscribeToOrdersWithRootOriginatorIdArgs,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/viewservice.ViewService/SubscribeToOrdersWithRootOriginatorId',
      request,
      metadata || {},
      this.methodInfoSubscribeToOrdersWithRootOriginatorId);
  }

  methodInfoGetOrderHistory = new grpcWeb.AbstractClientBase.MethodInfo(
    viewservice_pb.OrderHistory,
    (request: viewservice_pb.GetOrderHistoryArgs) => {
      return request.serializeBinary();
    },
    viewservice_pb.OrderHistory.deserializeBinary
  );

  getOrderHistory(
    request: viewservice_pb.GetOrderHistoryArgs,
    metadata: grpcWeb.Metadata | null): Promise<viewservice_pb.OrderHistory>;

  getOrderHistory(
    request: viewservice_pb.GetOrderHistoryArgs,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: viewservice_pb.OrderHistory) => void): grpcWeb.ClientReadableStream<viewservice_pb.OrderHistory>;

  getOrderHistory(
    request: viewservice_pb.GetOrderHistoryArgs,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: viewservice_pb.OrderHistory) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/viewservice.ViewService/GetOrderHistory',
        request,
        metadata || {},
        this.methodInfoGetOrderHistory,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/viewservice.ViewService/GetOrderHistory',
    request,
    metadata || {},
    this.methodInfoGetOrderHistory);
  }

}

