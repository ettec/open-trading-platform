/**
 * @fileoverview gRPC-Web generated client stub for staticdataservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


import * as grpcWeb from 'grpc-web';

import * as listing_pb from './listing_pb';

import {
  ExactMatchParameters,
  ListingId,
  ListingIds,
  Listings,
  MatchParameters} from './staticdataservice_pb';

export class StaticDataServiceClient {
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

  methodInfoGetListingsWithSameInstrument = new grpcWeb.AbstractClientBase.MethodInfo(
    Listings,
    (request: ListingId) => {
      return request.serializeBinary();
    },
    Listings.deserializeBinary
  );

  getListingsWithSameInstrument(
    request: ListingId,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: Listings) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/staticdataservice.StaticDataService/GetListingsWithSameInstrument',
      request,
      metadata || {},
      this.methodInfoGetListingsWithSameInstrument,
      callback);
  }

  methodInfoGetListingMatching = new grpcWeb.AbstractClientBase.MethodInfo(
    listing_pb.Listing,
    (request: ExactMatchParameters) => {
      return request.serializeBinary();
    },
    listing_pb.Listing.deserializeBinary
  );

  getListingMatching(
    request: ExactMatchParameters,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: listing_pb.Listing) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/staticdataservice.StaticDataService/GetListingMatching',
      request,
      metadata || {},
      this.methodInfoGetListingMatching,
      callback);
  }

  methodInfoGetListingsMatching = new grpcWeb.AbstractClientBase.MethodInfo(
    Listings,
    (request: MatchParameters) => {
      return request.serializeBinary();
    },
    Listings.deserializeBinary
  );

  getListingsMatching(
    request: MatchParameters,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: Listings) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/staticdataservice.StaticDataService/GetListingsMatching',
      request,
      metadata || {},
      this.methodInfoGetListingsMatching,
      callback);
  }

  methodInfoGetListing = new grpcWeb.AbstractClientBase.MethodInfo(
    listing_pb.Listing,
    (request: ListingId) => {
      return request.serializeBinary();
    },
    listing_pb.Listing.deserializeBinary
  );

  getListing(
    request: ListingId,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: listing_pb.Listing) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/staticdataservice.StaticDataService/GetListing',
      request,
      metadata || {},
      this.methodInfoGetListing,
      callback);
  }

  methodInfoGetListings = new grpcWeb.AbstractClientBase.MethodInfo(
    Listings,
    (request: ListingIds) => {
      return request.serializeBinary();
    },
    Listings.deserializeBinary
  );

  getListings(
    request: ListingIds,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: Listings) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/staticdataservice.StaticDataService/GetListings',
      request,
      metadata || {},
      this.methodInfoGetListings,
      callback);
  }

}

