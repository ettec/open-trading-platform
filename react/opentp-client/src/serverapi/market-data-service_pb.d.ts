import * as jspb from "google-protobuf"

import * as common_pb from './common_pb';

export class Subscription extends jspb.Message {
  getSubscriberid(): string;
  setSubscriberid(value: string): void;

  getListingid(): number;
  setListingid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Subscription.AsObject;
  static toObject(includeInstance: boolean, msg: Subscription): Subscription.AsObject;
  static serializeBinaryToWriter(message: Subscription, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Subscription;
  static deserializeBinaryFromReader(message: Subscription, reader: jspb.BinaryReader): Subscription;
}

export namespace Subscription {
  export type AsObject = {
    subscriberid: string,
    listingid: number,
  }
}

export class AddSubscriptionResponse extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddSubscriptionResponse): AddSubscriptionResponse.AsObject;
  static serializeBinaryToWriter(message: AddSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSubscriptionResponse;
  static deserializeBinaryFromReader(message: AddSubscriptionResponse, reader: jspb.BinaryReader): AddSubscriptionResponse;
}

export namespace AddSubscriptionResponse {
  export type AsObject = {
    message: string,
  }
}

export class SubscribeRequest extends jspb.Message {
  getSubscriberid(): string;
  setSubscriberid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscribeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SubscribeRequest): SubscribeRequest.AsObject;
  static serializeBinaryToWriter(message: SubscribeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubscribeRequest;
  static deserializeBinaryFromReader(message: SubscribeRequest, reader: jspb.BinaryReader): SubscribeRequest;
}

export namespace SubscribeRequest {
  export type AsObject = {
    subscriberid: string,
  }
}

export class DepthLine extends jspb.Message {
  getBidsize(): common_pb.Decimal64 | undefined;
  setBidsize(value?: common_pb.Decimal64): void;
  hasBidsize(): boolean;
  clearBidsize(): void;

  getBidprice(): common_pb.Decimal64 | undefined;
  setBidprice(value?: common_pb.Decimal64): void;
  hasBidprice(): boolean;
  clearBidprice(): void;

  getAsksize(): common_pb.Decimal64 | undefined;
  setAsksize(value?: common_pb.Decimal64): void;
  hasAsksize(): boolean;
  clearAsksize(): void;

  getAskprice(): common_pb.Decimal64 | undefined;
  setAskprice(value?: common_pb.Decimal64): void;
  hasAskprice(): boolean;
  clearAskprice(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DepthLine.AsObject;
  static toObject(includeInstance: boolean, msg: DepthLine): DepthLine.AsObject;
  static serializeBinaryToWriter(message: DepthLine, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DepthLine;
  static deserializeBinaryFromReader(message: DepthLine, reader: jspb.BinaryReader): DepthLine;
}

export namespace DepthLine {
  export type AsObject = {
    bidsize?: common_pb.Decimal64.AsObject,
    bidprice?: common_pb.Decimal64.AsObject,
    asksize?: common_pb.Decimal64.AsObject,
    askprice?: common_pb.Decimal64.AsObject,
  }
}

export class Quote extends jspb.Message {
  getListingid(): number;
  setListingid(value: number): void;

  getDepthList(): Array<DepthLine>;
  setDepthList(value: Array<DepthLine>): void;
  clearDepthList(): void;
  addDepth(value?: DepthLine, index?: number): DepthLine;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Quote.AsObject;
  static toObject(includeInstance: boolean, msg: Quote): Quote.AsObject;
  static serializeBinaryToWriter(message: Quote, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Quote;
  static deserializeBinaryFromReader(message: Quote, reader: jspb.BinaryReader): Quote;
}

export namespace Quote {
  export type AsObject = {
    listingid: number,
    depthList: Array<DepthLine.AsObject>,
  }
}

