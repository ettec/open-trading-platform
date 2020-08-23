/**
 * @fileoverview gRPC-Web generated client stub for marketdataservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as modelcommon_pb from './modelcommon_pb';
import * as clobquote_pb from './clobquote_pb';
import * as market$data$service_pb from './market-data-service_pb';


export class MarketDataServiceClient {
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

  methodInfoSubscribe = new grpcWeb.AbstractClientBase.MethodInfo(
    modelcommon_pb.Empty,
    (request: market$data$service_pb.MdsSubscribeRequest) => {
      return request.serializeBinary();
    },
    modelcommon_pb.Empty.deserializeBinary
  );

  subscribe(
    request: market$data$service_pb.MdsSubscribeRequest,
    metadata: grpcWeb.Metadata | null): Promise<modelcommon_pb.Empty>;

  subscribe(
    request: market$data$service_pb.MdsSubscribeRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void): grpcWeb.ClientReadableStream<modelcommon_pb.Empty>;

  subscribe(
    request: market$data$service_pb.MdsSubscribeRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/marketdataservice.MarketDataService/Subscribe',
        request,
        metadata || {},
        this.methodInfoSubscribe,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/marketdataservice.MarketDataService/Subscribe',
    request,
    metadata || {},
    this.methodInfoSubscribe);
  }

  methodInfoConnect = new grpcWeb.AbstractClientBase.MethodInfo(
    clobquote_pb.ClobQuote,
    (request: market$data$service_pb.MdsConnectRequest) => {
      return request.serializeBinary();
    },
    clobquote_pb.ClobQuote.deserializeBinary
  );

  connect(
    request: market$data$service_pb.MdsConnectRequest,
    metadata?: grpcWeb.Metadata) {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/marketdataservice.MarketDataService/Connect',
      request,
      metadata || {},
      this.methodInfoConnect);
  }

}

