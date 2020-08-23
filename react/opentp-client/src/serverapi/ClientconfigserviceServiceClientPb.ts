/**
 * @fileoverview gRPC-Web generated client stub for clientconfigservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as modelcommon_pb from './modelcommon_pb';
import * as clientconfigservice_pb from './clientconfigservice_pb';


export class ClientConfigServiceClient {
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

  methodInfoGetClientConfig = new grpcWeb.AbstractClientBase.MethodInfo(
    clientconfigservice_pb.Config,
    (request: clientconfigservice_pb.GetConfigParameters) => {
      return request.serializeBinary();
    },
    clientconfigservice_pb.Config.deserializeBinary
  );

  getClientConfig(
    request: clientconfigservice_pb.GetConfigParameters,
    metadata: grpcWeb.Metadata | null): Promise<clientconfigservice_pb.Config>;

  getClientConfig(
    request: clientconfigservice_pb.GetConfigParameters,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: clientconfigservice_pb.Config) => void): grpcWeb.ClientReadableStream<clientconfigservice_pb.Config>;

  getClientConfig(
    request: clientconfigservice_pb.GetConfigParameters,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: clientconfigservice_pb.Config) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/clientconfigservice.ClientConfigService/GetClientConfig',
        request,
        metadata || {},
        this.methodInfoGetClientConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/clientconfigservice.ClientConfigService/GetClientConfig',
    request,
    metadata || {},
    this.methodInfoGetClientConfig);
  }

  methodInfoStoreClientConfig = new grpcWeb.AbstractClientBase.MethodInfo(
    modelcommon_pb.Empty,
    (request: clientconfigservice_pb.StoreConfigParams) => {
      return request.serializeBinary();
    },
    modelcommon_pb.Empty.deserializeBinary
  );

  storeClientConfig(
    request: clientconfigservice_pb.StoreConfigParams,
    metadata: grpcWeb.Metadata | null): Promise<modelcommon_pb.Empty>;

  storeClientConfig(
    request: clientconfigservice_pb.StoreConfigParams,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void): grpcWeb.ClientReadableStream<modelcommon_pb.Empty>;

  storeClientConfig(
    request: clientconfigservice_pb.StoreConfigParams,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/clientconfigservice.ClientConfigService/StoreClientConfig',
        request,
        metadata || {},
        this.methodInfoStoreClientConfig,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/clientconfigservice.ClientConfigService/StoreClientConfig',
    request,
    metadata || {},
    this.methodInfoStoreClientConfig);
  }

}

