/**
 * @fileoverview gRPC-Web generated client stub for executionvenue
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as modelcommon_pb from './modelcommon_pb';
import * as executionvenue_pb from './executionvenue_pb';


export class ExecutionVenueClient {
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

  methodInfoCreateAndRouteOrder = new grpcWeb.AbstractClientBase.MethodInfo(
    executionvenue_pb.OrderId,
    (request: executionvenue_pb.CreateAndRouteOrderParams) => {
      return request.serializeBinary();
    },
    executionvenue_pb.OrderId.deserializeBinary
  );

  createAndRouteOrder(
    request: executionvenue_pb.CreateAndRouteOrderParams,
    metadata: grpcWeb.Metadata | null): Promise<executionvenue_pb.OrderId>;

  createAndRouteOrder(
    request: executionvenue_pb.CreateAndRouteOrderParams,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: executionvenue_pb.OrderId) => void): grpcWeb.ClientReadableStream<executionvenue_pb.OrderId>;

  createAndRouteOrder(
    request: executionvenue_pb.CreateAndRouteOrderParams,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: executionvenue_pb.OrderId) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/executionvenue.ExecutionVenue/CreateAndRouteOrder',
        request,
        metadata || {},
        this.methodInfoCreateAndRouteOrder,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/executionvenue.ExecutionVenue/CreateAndRouteOrder',
    request,
    metadata || {},
    this.methodInfoCreateAndRouteOrder);
  }

  methodInfoCancelOrder = new grpcWeb.AbstractClientBase.MethodInfo(
    modelcommon_pb.Empty,
    (request: executionvenue_pb.CancelOrderParams) => {
      return request.serializeBinary();
    },
    modelcommon_pb.Empty.deserializeBinary
  );

  cancelOrder(
    request: executionvenue_pb.CancelOrderParams,
    metadata: grpcWeb.Metadata | null): Promise<modelcommon_pb.Empty>;

  cancelOrder(
    request: executionvenue_pb.CancelOrderParams,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void): grpcWeb.ClientReadableStream<modelcommon_pb.Empty>;

  cancelOrder(
    request: executionvenue_pb.CancelOrderParams,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/executionvenue.ExecutionVenue/CancelOrder',
        request,
        metadata || {},
        this.methodInfoCancelOrder,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/executionvenue.ExecutionVenue/CancelOrder',
    request,
    metadata || {},
    this.methodInfoCancelOrder);
  }

  methodInfoModifyOrder = new grpcWeb.AbstractClientBase.MethodInfo(
    modelcommon_pb.Empty,
    (request: executionvenue_pb.ModifyOrderParams) => {
      return request.serializeBinary();
    },
    modelcommon_pb.Empty.deserializeBinary
  );

  modifyOrder(
    request: executionvenue_pb.ModifyOrderParams,
    metadata: grpcWeb.Metadata | null): Promise<modelcommon_pb.Empty>;

  modifyOrder(
    request: executionvenue_pb.ModifyOrderParams,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void): grpcWeb.ClientReadableStream<modelcommon_pb.Empty>;

  modifyOrder(
    request: executionvenue_pb.ModifyOrderParams,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/executionvenue.ExecutionVenue/ModifyOrder',
        request,
        metadata || {},
        this.methodInfoModifyOrder,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/executionvenue.ExecutionVenue/ModifyOrder',
    request,
    metadata || {},
    this.methodInfoModifyOrder);
  }

  methodInfoGetExecutionParametersMetaData = new grpcWeb.AbstractClientBase.MethodInfo(
    executionvenue_pb.ExecParamsMetaDataJson,
    (request: modelcommon_pb.Empty) => {
      return request.serializeBinary();
    },
    executionvenue_pb.ExecParamsMetaDataJson.deserializeBinary
  );

  getExecutionParametersMetaData(
    request: modelcommon_pb.Empty,
    metadata: grpcWeb.Metadata | null): Promise<executionvenue_pb.ExecParamsMetaDataJson>;

  getExecutionParametersMetaData(
    request: modelcommon_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: executionvenue_pb.ExecParamsMetaDataJson) => void): grpcWeb.ClientReadableStream<executionvenue_pb.ExecParamsMetaDataJson>;

  getExecutionParametersMetaData(
    request: modelcommon_pb.Empty,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: executionvenue_pb.ExecParamsMetaDataJson) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/executionvenue.ExecutionVenue/GetExecutionParametersMetaData',
        request,
        metadata || {},
        this.methodInfoGetExecutionParametersMetaData,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/executionvenue.ExecutionVenue/GetExecutionParametersMetaData',
    request,
    metadata || {},
    this.methodInfoGetExecutionParametersMetaData);
  }

}

