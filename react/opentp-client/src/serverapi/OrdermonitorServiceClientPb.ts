/**
 * @fileoverview gRPC-Web generated client stub for ordermonitor
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as modelcommon_pb from './modelcommon_pb';
import * as ordermonitor_pb from './ordermonitor_pb';


export class OrderMonitorClient {
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

  methodInfoCancelAllOrdersForOriginatorId = new grpcWeb.AbstractClientBase.MethodInfo(
    modelcommon_pb.Empty,
    (request: ordermonitor_pb.CancelAllOrdersForOriginatorIdParams) => {
      return request.serializeBinary();
    },
    modelcommon_pb.Empty.deserializeBinary
  );

  cancelAllOrdersForOriginatorId(
    request: ordermonitor_pb.CancelAllOrdersForOriginatorIdParams,
    metadata: grpcWeb.Metadata | null): Promise<modelcommon_pb.Empty>;

  cancelAllOrdersForOriginatorId(
    request: ordermonitor_pb.CancelAllOrdersForOriginatorIdParams,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void): grpcWeb.ClientReadableStream<modelcommon_pb.Empty>;

  cancelAllOrdersForOriginatorId(
    request: ordermonitor_pb.CancelAllOrdersForOriginatorIdParams,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/ordermonitor.OrderMonitor/CancelAllOrdersForOriginatorId',
        request,
        metadata || {},
        this.methodInfoCancelAllOrdersForOriginatorId,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/ordermonitor.OrderMonitor/CancelAllOrdersForOriginatorId',
    request,
    metadata || {},
    this.methodInfoCancelAllOrdersForOriginatorId);
  }

}

