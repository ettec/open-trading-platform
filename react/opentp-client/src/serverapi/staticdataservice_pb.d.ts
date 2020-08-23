import * as jspb from 'google-protobuf'

import * as listing_pb from './listing_pb';


export class ListingId extends jspb.Message {
  getListingid(): number;
  setListingid(value: number): ListingId;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListingId.AsObject;
  static toObject(includeInstance: boolean, msg: ListingId): ListingId.AsObject;
  static serializeBinaryToWriter(message: ListingId, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListingId;
  static deserializeBinaryFromReader(message: ListingId, reader: jspb.BinaryReader): ListingId;
}

export namespace ListingId {
  export type AsObject = {
    listingid: number,
  }
}

export class ListingIds extends jspb.Message {
  getListingidsList(): Array<number>;
  setListingidsList(value: Array<number>): ListingIds;
  clearListingidsList(): ListingIds;
  addListingids(value: number, index?: number): ListingIds;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListingIds.AsObject;
  static toObject(includeInstance: boolean, msg: ListingIds): ListingIds.AsObject;
  static serializeBinaryToWriter(message: ListingIds, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListingIds;
  static deserializeBinaryFromReader(message: ListingIds, reader: jspb.BinaryReader): ListingIds;
}

export namespace ListingIds {
  export type AsObject = {
    listingidsList: Array<number>,
  }
}

export class Listings extends jspb.Message {
  getListingsList(): Array<listing_pb.Listing>;
  setListingsList(value: Array<listing_pb.Listing>): Listings;
  clearListingsList(): Listings;
  addListings(value?: listing_pb.Listing, index?: number): listing_pb.Listing;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Listings.AsObject;
  static toObject(includeInstance: boolean, msg: Listings): Listings.AsObject;
  static serializeBinaryToWriter(message: Listings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Listings;
  static deserializeBinaryFromReader(message: Listings, reader: jspb.BinaryReader): Listings;
}

export namespace Listings {
  export type AsObject = {
    listingsList: Array<listing_pb.Listing.AsObject>,
  }
}

export class MatchParameters extends jspb.Message {
  getSymbolmatch(): string;
  setSymbolmatch(value: string): MatchParameters;

  getNamematch(): string;
  setNamematch(value: string): MatchParameters;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchParameters.AsObject;
  static toObject(includeInstance: boolean, msg: MatchParameters): MatchParameters.AsObject;
  static serializeBinaryToWriter(message: MatchParameters, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchParameters;
  static deserializeBinaryFromReader(message: MatchParameters, reader: jspb.BinaryReader): MatchParameters;
}

export namespace MatchParameters {
  export type AsObject = {
    symbolmatch: string,
    namematch: string,
  }
}

export class ExactMatchParameters extends jspb.Message {
  getSymbol(): string;
  setSymbol(value: string): ExactMatchParameters;

  getMic(): string;
  setMic(value: string): ExactMatchParameters;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExactMatchParameters.AsObject;
  static toObject(includeInstance: boolean, msg: ExactMatchParameters): ExactMatchParameters.AsObject;
  static serializeBinaryToWriter(message: ExactMatchParameters, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExactMatchParameters;
  static deserializeBinaryFromReader(message: ExactMatchParameters, reader: jspb.BinaryReader): ExactMatchParameters;
}

export namespace ExactMatchParameters {
  export type AsObject = {
    symbol: string,
    mic: string,
  }
}

