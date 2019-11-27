/**
 * @fileoverview gRPC-Web generated client stub for clientmarketdataservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import {
  AddSubscriptionResponse,
  Book,
  SubscribeRequest,
  Subscription} from './cmds_pb';

export class ClientMarketDataServiceClient {
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

  methodInfoAddSubscription = new grpcWeb.AbstractClientBase.MethodInfo(
    AddSubscriptionResponse,
    (request: Subscription) => {
      return request.serializeBinary();
    },
    AddSubscriptionResponse.deserializeBinary
  );

  addSubscription(
    request: Subscription,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: AddSubscriptionResponse) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/clientmarketdataservice.ClientMarketDataService/AddSubscription',
      request,
      metadata || {},
      this.methodInfoAddSubscription,
      callback);
  }

  methodInfoSubscribe = new grpcWeb.AbstractClientBase.MethodInfo(
    Book,
    (request: SubscribeRequest) => {
      return request.serializeBinary();
    },
    Book.deserializeBinary
  );

  subscribe(
    request: SubscribeRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/clientmarketdataservice.ClientMarketDataService/Subscribe',
      request,
      metadata || {},
      this.methodInfoSubscribe);
  }

}

