/**
 * @fileoverview gRPC-Web generated client stub for staticdataservice
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as listing_pb from './listing_pb';
import * as staticdataservice_pb from './staticdataservice_pb';


export class StaticDataServiceClient {
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

  methodInfoGetListingsWithSameInstrument = new grpcWeb.AbstractClientBase.MethodInfo(
    staticdataservice_pb.Listings,
    (request: staticdataservice_pb.ListingId) => {
      return request.serializeBinary();
    },
    staticdataservice_pb.Listings.deserializeBinary
  );

  getListingsWithSameInstrument(
    request: staticdataservice_pb.ListingId,
    metadata: grpcWeb.Metadata | null): Promise<staticdataservice_pb.Listings>;

  getListingsWithSameInstrument(
    request: staticdataservice_pb.ListingId,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: staticdataservice_pb.Listings) => void): grpcWeb.ClientReadableStream<staticdataservice_pb.Listings>;

  getListingsWithSameInstrument(
    request: staticdataservice_pb.ListingId,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: staticdataservice_pb.Listings) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/staticdataservice.StaticDataService/GetListingsWithSameInstrument',
        request,
        metadata || {},
        this.methodInfoGetListingsWithSameInstrument,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/staticdataservice.StaticDataService/GetListingsWithSameInstrument',
    request,
    metadata || {},
    this.methodInfoGetListingsWithSameInstrument);
  }

  methodInfoGetListingMatching = new grpcWeb.AbstractClientBase.MethodInfo(
    listing_pb.Listing,
    (request: staticdataservice_pb.ExactMatchParameters) => {
      return request.serializeBinary();
    },
    listing_pb.Listing.deserializeBinary
  );

  getListingMatching(
    request: staticdataservice_pb.ExactMatchParameters,
    metadata: grpcWeb.Metadata | null): Promise<listing_pb.Listing>;

  getListingMatching(
    request: staticdataservice_pb.ExactMatchParameters,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: listing_pb.Listing) => void): grpcWeb.ClientReadableStream<listing_pb.Listing>;

  getListingMatching(
    request: staticdataservice_pb.ExactMatchParameters,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: listing_pb.Listing) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/staticdataservice.StaticDataService/GetListingMatching',
        request,
        metadata || {},
        this.methodInfoGetListingMatching,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/staticdataservice.StaticDataService/GetListingMatching',
    request,
    metadata || {},
    this.methodInfoGetListingMatching);
  }

  methodInfoGetListingsMatching = new grpcWeb.AbstractClientBase.MethodInfo(
    staticdataservice_pb.Listings,
    (request: staticdataservice_pb.MatchParameters) => {
      return request.serializeBinary();
    },
    staticdataservice_pb.Listings.deserializeBinary
  );

  getListingsMatching(
    request: staticdataservice_pb.MatchParameters,
    metadata: grpcWeb.Metadata | null): Promise<staticdataservice_pb.Listings>;

  getListingsMatching(
    request: staticdataservice_pb.MatchParameters,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: staticdataservice_pb.Listings) => void): grpcWeb.ClientReadableStream<staticdataservice_pb.Listings>;

  getListingsMatching(
    request: staticdataservice_pb.MatchParameters,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: staticdataservice_pb.Listings) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/staticdataservice.StaticDataService/GetListingsMatching',
        request,
        metadata || {},
        this.methodInfoGetListingsMatching,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/staticdataservice.StaticDataService/GetListingsMatching',
    request,
    metadata || {},
    this.methodInfoGetListingsMatching);
  }

  methodInfoGetListing = new grpcWeb.AbstractClientBase.MethodInfo(
    listing_pb.Listing,
    (request: staticdataservice_pb.ListingId) => {
      return request.serializeBinary();
    },
    listing_pb.Listing.deserializeBinary
  );

  getListing(
    request: staticdataservice_pb.ListingId,
    metadata: grpcWeb.Metadata | null): Promise<listing_pb.Listing>;

  getListing(
    request: staticdataservice_pb.ListingId,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: listing_pb.Listing) => void): grpcWeb.ClientReadableStream<listing_pb.Listing>;

  getListing(
    request: staticdataservice_pb.ListingId,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: listing_pb.Listing) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/staticdataservice.StaticDataService/GetListing',
        request,
        metadata || {},
        this.methodInfoGetListing,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/staticdataservice.StaticDataService/GetListing',
    request,
    metadata || {},
    this.methodInfoGetListing);
  }

  methodInfoGetListings = new grpcWeb.AbstractClientBase.MethodInfo(
    staticdataservice_pb.Listings,
    (request: staticdataservice_pb.ListingIds) => {
      return request.serializeBinary();
    },
    staticdataservice_pb.Listings.deserializeBinary
  );

  getListings(
    request: staticdataservice_pb.ListingIds,
    metadata: grpcWeb.Metadata | null): Promise<staticdataservice_pb.Listings>;

  getListings(
    request: staticdataservice_pb.ListingIds,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: staticdataservice_pb.Listings) => void): grpcWeb.ClientReadableStream<staticdataservice_pb.Listings>;

  getListings(
    request: staticdataservice_pb.ListingIds,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.Error,
               response: staticdataservice_pb.Listings) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/staticdataservice.StaticDataService/GetListings',
        request,
        metadata || {},
        this.methodInfoGetListings,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/staticdataservice.StaticDataService/GetListings',
    request,
    metadata || {},
    this.methodInfoGetListings);
  }

}

