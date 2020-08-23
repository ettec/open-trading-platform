/**
 * @fileoverview gRPC-Web generated client stub for loginservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as loginservice_pb from './loginservice_pb';


export class LoginServiceClient {
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

  methodInfoLogin = new grpcWeb.AbstractClientBase.MethodInfo(
    loginservice_pb.Token,
    (request: loginservice_pb.LoginParams) => {
      return request.serializeBinary();
    },
    loginservice_pb.Token.deserializeBinary
  );

  login(
    request: loginservice_pb.LoginParams,
    metadata: grpcWeb.Metadata | null): Promise<loginservice_pb.Token>;

  login(
    request: loginservice_pb.LoginParams,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: loginservice_pb.Token) => void): grpcWeb.ClientReadableStream<loginservice_pb.Token>;

  login(
    request: loginservice_pb.LoginParams,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: loginservice_pb.Token) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/loginservice.LoginService/Login',
        request,
        metadata || {},
        this.methodInfoLogin,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/loginservice.LoginService/Login',
    request,
    metadata || {},
    this.methodInfoLogin);
  }

}

