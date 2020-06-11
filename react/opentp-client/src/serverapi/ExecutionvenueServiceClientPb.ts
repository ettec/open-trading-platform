/**
 * @fileoverview gRPC-Web generated client stub for executionvenue
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import * as listing_pb from './listing_pb';
import * as order_pb from './order_pb';
import * as modelcommon_pb from './modelcommon_pb';

import {
  CancelOrderParams,
  CreateAndRouteOrderParams,
  ModifyOrderParams,
  OrderId} from './executionvenue_pb';

export class ExecutionVenueClient {
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

  methodInfoCreateAndRouteOrder = new grpcWeb.AbstractClientBase.MethodInfo(
    OrderId,
    (request: CreateAndRouteOrderParams) => {
      return request.serializeBinary();
    },
    OrderId.deserializeBinary
  );

  createAndRouteOrder(
    request: CreateAndRouteOrderParams,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: OrderId) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/executionvenue.ExecutionVenue/CreateAndRouteOrder',
      request,
      metadata || {},
      this.methodInfoCreateAndRouteOrder,
      callback);
  }

  methodInfoCancelOrder = new grpcWeb.AbstractClientBase.MethodInfo(
    modelcommon_pb.Empty,
    (request: CancelOrderParams) => {
      return request.serializeBinary();
    },
    modelcommon_pb.Empty.deserializeBinary
  );

  cancelOrder(
    request: CancelOrderParams,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/executionvenue.ExecutionVenue/CancelOrder',
      request,
      metadata || {},
      this.methodInfoCancelOrder,
      callback);
  }

  methodInfoModifyOrder = new grpcWeb.AbstractClientBase.MethodInfo(
    modelcommon_pb.Empty,
    (request: ModifyOrderParams) => {
      return request.serializeBinary();
    },
    modelcommon_pb.Empty.deserializeBinary
  );

  modifyOrder(
    request: ModifyOrderParams,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/executionvenue.ExecutionVenue/ModifyOrder',
      request,
      metadata || {},
      this.methodInfoModifyOrder,
      callback);
  }

}

