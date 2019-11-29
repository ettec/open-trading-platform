/**
 * @fileoverview gRPC-Web generated client stub for executionvenue
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import * as order_pb from './order_pb';
import * as common_pb from './common_pb';

import {
  CreateAndRouteOrderParams,
  OrderId} from './execution-venue_pb';

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
    common_pb.Empty,
    (request: OrderId) => {
      return request.serializeBinary();
    },
    common_pb.Empty.deserializeBinary
  );

  cancelOrder(
    request: OrderId,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: common_pb.Empty) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/executionvenue.ExecutionVenue/CancelOrder',
      request,
      metadata || {},
      this.methodInfoCancelOrder,
      callback);
  }

}

