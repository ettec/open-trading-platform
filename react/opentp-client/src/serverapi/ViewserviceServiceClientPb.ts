/**
 * @fileoverview gRPC-Web generated client stub for viewservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import * as modelcommon_pb from './modelcommon_pb';
import * as order_pb from './order_pb';

import {
  GetOrderHistoryArgs,
  Orders,
  SubscribeToOrdersWithRootOriginatorIdArgs} from './viewservice_pb';

export class ViewServiceClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: string; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; }) {
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
    (request: SubscribeToOrdersWithRootOriginatorIdArgs) => {
      return request.serializeBinary();
    },
    order_pb.Order.deserializeBinary
  );

  subscribeToOrdersWithRootOriginatorId(
    request: SubscribeToOrdersWithRootOriginatorIdArgs,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/viewservice.ViewService/SubscribeToOrdersWithRootOriginatorId',
      request,
      metadata || {},
      this.methodInfoSubscribeToOrdersWithRootOriginatorId);
  }

  methodInfoGetOrderHistory = new grpcWeb.AbstractClientBase.MethodInfo(
    Orders,
    (request: GetOrderHistoryArgs) => {
      return request.serializeBinary();
    },
    Orders.deserializeBinary
  );

  getOrderHistory(
    request: GetOrderHistoryArgs,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: Orders) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/viewservice.ViewService/GetOrderHistory',
      request,
      metadata || {},
      this.methodInfoGetOrderHistory,
      callback);
  }

}

