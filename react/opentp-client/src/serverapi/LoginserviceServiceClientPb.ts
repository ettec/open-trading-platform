/**
 * @fileoverview gRPC-Web generated client stub for loginservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import {
  LoginParams,
  Token} from './loginservice_pb';

export class LoginServiceClient {
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

  methodInfoLogin = new grpcWeb.AbstractClientBase.MethodInfo(
    Token,
    (request: LoginParams) => {
      return request.serializeBinary();
    },
    Token.deserializeBinary
  );

  login(
    request: LoginParams,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: Token) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/loginservice.LoginService/Login',
      request,
      metadata || {},
      this.methodInfoLogin,
      callback);
  }

}

