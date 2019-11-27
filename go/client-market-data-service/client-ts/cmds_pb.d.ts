import * as jspb from "google-protobuf"

export class Subscription extends jspb.Message {
  getSubscriberid(): string;
  setSubscriberid(value: string): void;

  getListingid(): string;
  setListingid(value: string): void;

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
    listingid: string,
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

export class BookLine extends jspb.Message {
  getBidsize(): string;
  setBidsize(value: string): void;

  getBidprice(): string;
  setBidprice(value: string): void;

  getAsksize(): string;
  setAsksize(value: string): void;

  getAskprice(): string;
  setAskprice(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BookLine.AsObject;
  static toObject(includeInstance: boolean, msg: BookLine): BookLine.AsObject;
  static serializeBinaryToWriter(message: BookLine, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BookLine;
  static deserializeBinaryFromReader(message: BookLine, reader: jspb.BinaryReader): BookLine;
}

export namespace BookLine {
  export type AsObject = {
    bidsize: string,
    bidprice: string,
    asksize: string,
    askprice: string,
  }
}

export class Book extends jspb.Message {
  getListingid(): string;
  setListingid(value: string): void;

  getDepthList(): Array<BookLine>;
  setDepthList(value: Array<BookLine>): void;
  clearDepthList(): void;
  addDepth(value?: BookLine, index?: number): BookLine;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Book.AsObject;
  static toObject(includeInstance: boolean, msg: Book): Book.AsObject;
  static serializeBinaryToWriter(message: Book, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Book;
  static deserializeBinaryFromReader(message: Book, reader: jspb.BinaryReader): Book;
}

export namespace Book {
  export type AsObject = {
    listingid: string,
    depthList: Array<BookLine.AsObject>,
  }
}

