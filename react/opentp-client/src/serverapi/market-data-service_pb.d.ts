import * as jspb from 'google-protobuf'

import * as listing_pb from './listing_pb';
import * as modelcommon_pb from './modelcommon_pb';
import * as clobquote_pb from './clobquote_pb';


export class MdsConnectRequest extends jspb.Message {
  getSubscriberid(): string;
  setSubscriberid(value: string): MdsConnectRequest;

  getMaxquotepersecond(): number;
  setMaxquotepersecond(value: number): MdsConnectRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MdsConnectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: MdsConnectRequest): MdsConnectRequest.AsObject;
  static serializeBinaryToWriter(message: MdsConnectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MdsConnectRequest;
  static deserializeBinaryFromReader(message: MdsConnectRequest, reader: jspb.BinaryReader): MdsConnectRequest;
}

export namespace MdsConnectRequest {
  export type AsObject = {
    subscriberid: string,
    maxquotepersecond: number,
  }
}

export class MdsSubscribeRequest extends jspb.Message {
  getSubscriberid(): string;
  setSubscriberid(value: string): MdsSubscribeRequest;

  getListing(): listing_pb.Listing | undefined;
  setListing(value?: listing_pb.Listing): MdsSubscribeRequest;
  hasListing(): boolean;
  clearListing(): MdsSubscribeRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MdsSubscribeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: MdsSubscribeRequest): MdsSubscribeRequest.AsObject;
  static serializeBinaryToWriter(message: MdsSubscribeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MdsSubscribeRequest;
  static deserializeBinaryFromReader(message: MdsSubscribeRequest, reader: jspb.BinaryReader): MdsSubscribeRequest;
}

export namespace MdsSubscribeRequest {
  export type AsObject = {
    subscriberid: string,
    listing?: listing_pb.Listing.AsObject,
  }
}

