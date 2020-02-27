/**
 * @fileoverview gRPC-Web generated client stub for viewservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import * as modelcommon_pb from './modelcommon_pb';
import * as order_pb from './order_pb';

import {SubscribeToOrders} from './view-service_pb';

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

  methodInfoSubscribe = new grpcWeb.AbstractClientBase.MethodInfo(
    order_pb.Order,
    (request: SubscribeToOrders) => {
      return request.serializeBinary();
    },
    order_pb.Order.deserializeBinary
  );

  subscribe(
    request: SubscribeToOrders,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/viewservice.ViewService/Subscribe',
      request,
      metadata || {},
      this.methodInfoSubscribe);
  }

}

