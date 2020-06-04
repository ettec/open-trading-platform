/**
 * @fileoverview gRPC-Web generated client stub for clientconfigservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import * as modelcommon_pb from './modelcommon_pb';

import {
  Config,
  GetConfigParameters,
  StoreConfigParams} from './clientconfigservice_pb';

export class ClientConfigServiceClient {
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

  methodInfoGetClientConfig = new grpcWeb.AbstractClientBase.MethodInfo(
    Config,
    (request: GetConfigParameters) => {
      return request.serializeBinary();
    },
    Config.deserializeBinary
  );

  getClientConfig(
    request: GetConfigParameters,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: Config) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/clientconfigservice.ClientConfigService/GetClientConfig',
      request,
      metadata || {},
      this.methodInfoGetClientConfig,
      callback);
  }

  methodInfoStoreClientConfig = new grpcWeb.AbstractClientBase.MethodInfo(
    modelcommon_pb.Empty,
    (request: StoreConfigParams) => {
      return request.serializeBinary();
    },
    modelcommon_pb.Empty.deserializeBinary
  );

  storeClientConfig(
    request: StoreConfigParams,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: modelcommon_pb.Empty) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/clientconfigservice.ClientConfigService/StoreClientConfig',
      request,
      metadata || {},
      this.methodInfoStoreClientConfig,
      callback);
  }

}

