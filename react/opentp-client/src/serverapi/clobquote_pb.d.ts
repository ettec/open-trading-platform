import * as jspb from "google-protobuf"

import * as modelcommon_pb from './modelcommon_pb';

export class ClobLine extends jspb.Message {
  getSize(): modelcommon_pb.Decimal64 | undefined;
  setSize(value?: modelcommon_pb.Decimal64): void;
  hasSize(): boolean;
  clearSize(): void;

  getPrice(): modelcommon_pb.Decimal64 | undefined;
  setPrice(value?: modelcommon_pb.Decimal64): void;
  hasPrice(): boolean;
  clearPrice(): void;

  getEntryid(): string;
  setEntryid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClobLine.AsObject;
  static toObject(includeInstance: boolean, msg: ClobLine): ClobLine.AsObject;
  static serializeBinaryToWriter(message: ClobLine, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClobLine;
  static deserializeBinaryFromReader(message: ClobLine, reader: jspb.BinaryReader): ClobLine;
}

export namespace ClobLine {
  export type AsObject = {
    size?: modelcommon_pb.Decimal64.AsObject,
    price?: modelcommon_pb.Decimal64.AsObject,
    entryid: string,
  }
}

export class ClobQuote extends jspb.Message {
  getListingid(): number;
  setListingid(value: number): void;

  getBidsList(): Array<ClobLine>;
  setBidsList(value: Array<ClobLine>): void;
  clearBidsList(): void;
  addBids(value?: ClobLine, index?: number): ClobLine;

  getOffersList(): Array<ClobLine>;
  setOffersList(value: Array<ClobLine>): void;
  clearOffersList(): void;
  addOffers(value?: ClobLine, index?: number): ClobLine;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClobQuote.AsObject;
  static toObject(includeInstance: boolean, msg: ClobQuote): ClobQuote.AsObject;
  static serializeBinaryToWriter(message: ClobQuote, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClobQuote;
  static deserializeBinaryFromReader(message: ClobQuote, reader: jspb.BinaryReader): ClobQuote;
}

export namespace ClobQuote {
  export type AsObject = {
    listingid: number,
    bidsList: Array<ClobLine.AsObject>,
    offersList: Array<ClobLine.AsObject>,
  }
}

