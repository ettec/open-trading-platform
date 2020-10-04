/**
 * @fileoverview gRPC-Web generated client stub for orderdataservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as order_pb from './order_pb';
import * as orderdataservice_pb from './orderdataservice_pb';


export class OrderDataServiceClient {
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
    (request: orderdataservice_pb.SubscribeToOrdersWithRootOriginatorIdArgs) => {
      return request.serializeBinary();
    },
    order_pb.Order.deserializeBinary
  );

  subscribeToOrdersWithRootOriginatorId(
    request: orderdataservice_pb.SubscribeToOrdersWithRootOriginatorIdArgs,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/orderdataservice.OrderDataService/SubscribeToOrdersWithRootOriginatorId',
      request,
      metadata || {},
      this.methodInfoSubscribeToOrdersWithRootOriginatorId);
  }

  methodInfoGetOrderHistory = new grpcWeb.AbstractClientBase.MethodInfo(
    orderdataservice_pb.OrderHistory,
    (request: orderdataservice_pb.GetOrderHistoryArgs) => {
      return request.serializeBinary();
    },
    orderdataservice_pb.OrderHistory.deserializeBinary
  );

  getOrderHistory(
    request: orderdataservice_pb.GetOrderHistoryArgs,
    metadata: grpcWeb.Metadata | null): Promise<orderdataservice_pb.OrderHistory>;

  getOrderHistory(
    request: orderdataservice_pb.GetOrderHistoryArgs,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: orderdataservice_pb.OrderHistory) => void): grpcWeb.ClientReadableStream<orderdataservice_pb.OrderHistory>;

  getOrderHistory(
    request: orderdataservice_pb.GetOrderHistoryArgs,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: orderdataservice_pb.OrderHistory) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/orderdataservice.OrderDataService/GetOrderHistory',
        request,
        metadata || {},
        this.methodInfoGetOrderHistory,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/orderdataservice.OrderDataService/GetOrderHistory',
    request,
    metadata || {},
    this.methodInfoGetOrderHistory);
  }

}

