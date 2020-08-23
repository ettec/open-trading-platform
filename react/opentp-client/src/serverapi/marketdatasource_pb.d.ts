import * as jspb from 'google-protobuf'

import * as clobquote_pb from './clobquote_pb';
import * as modelcommon_pb from './modelcommon_pb';


export class SubscribeRequest extends jspb.Message {
  getListingid(): number;
  setListingid(value: number): SubscribeRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscribeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SubscribeRequest): SubscribeRequest.AsObject;
  static serializeBinaryToWriter(message: SubscribeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubscribeRequest;
  static deserializeBinaryFromReader(message: SubscribeRequest, reader: jspb.BinaryReader): SubscribeRequest;
}

export namespace SubscribeRequest {
  export type AsObject = {
    listingid: number,
  }
}

