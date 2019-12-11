/**
 * @fileoverview gRPC-Web generated client stub for marketdataservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import * as common_pb from './common_pb';

import {
  AddSubscriptionResponse,
  Quote,
  SubscribeRequest,
  Subscription} from './market-data-service_pb';

export class MarketDataServiceClient {
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
        '/marketdataservice.MarketDataService/AddSubscription',
      request,
      metadata || {},
      this.methodInfoAddSubscription,
      callback);
  }

  methodInfoSubscribe = new grpcWeb.AbstractClientBase.MethodInfo(
    Quote,
    (request: SubscribeRequest) => {
      return request.serializeBinary();
    },
    Quote.deserializeBinary
  );

  subscribe(
    request: SubscribeRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/marketdataservice.MarketDataService/Subscribe',
      request,
      metadata || {},
      this.methodInfoSubscribe);
  }

}

