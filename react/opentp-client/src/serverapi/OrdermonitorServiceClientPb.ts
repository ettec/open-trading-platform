/**
 * @fileoverview gRPC-Web generated client stub for ordermonitor
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import * as modelcommon_pb from './modelcommon_pb';

import {CancelAllOrdersForOriginatorIdParams} from './ordermonitor_pb';

export class OrderMonitorClient {
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

  methodInfoCancelAllOrdersForOriginatorId = new grpcWeb.AbstractClientBase.MethodInfo(
    modelcommon_pb.Empty,
    (request: CancelAllOrdersForOriginatorIdParams) => {
      return request.serializeBinary();
    },
    modelcommon_pb.Empty.deserializeBinary
  );

  cancelAllOrdersForOriginatorId(
    request: CancelAllOrdersForOriginatorIdParams,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/ordermonitor.OrderMonitor/CancelAllOrdersForOriginatorId',
      request,
      metadata || {},
      this.methodInfoCancelAllOrdersForOriginatorId,
      callback);
  }

}

