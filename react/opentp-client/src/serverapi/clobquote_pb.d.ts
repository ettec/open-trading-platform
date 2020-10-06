import * as jspb from 'google-protobuf'

import * as modelcommon_pb from './modelcommon_pb';


export class ClobLine extends jspb.Message {
  getSize(): modelcommon_pb.Decimal64 | undefined;
  setSize(value?: modelcommon_pb.Decimal64): ClobLine;
  hasSize(): boolean;
  clearSize(): ClobLine;

  getPrice(): modelcommon_pb.Decimal64 | undefined;
  setPrice(value?: modelcommon_pb.Decimal64): ClobLine;
  hasPrice(): boolean;
  clearPrice(): ClobLine;

  getEntryid(): string;
  setEntryid(value: string): ClobLine;

  getListingid(): number;
  setListingid(value: number): ClobLine;

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
    listingid: number,
  }
}

export class ClobQuote extends jspb.Message {
  getListingid(): number;
  setListingid(value: number): ClobQuote;

  getBidsList(): Array<ClobLine>;
  setBidsList(value: Array<ClobLine>): ClobQuote;
  clearBidsList(): ClobQuote;
  addBids(value?: ClobLine, index?: number): ClobLine;

  getOffersList(): Array<ClobLine>;
  setOffersList(value: Array<ClobLine>): ClobQuote;
  clearOffersList(): ClobQuote;
  addOffers(value?: ClobLine, index?: number): ClobLine;

  getStreaminterrupted(): boolean;
  setStreaminterrupted(value: boolean): ClobQuote;

  getStreamstatusmsg(): string;
  setStreamstatusmsg(value: string): ClobQuote;

  getLastprice(): modelcommon_pb.Decimal64 | undefined;
  setLastprice(value?: modelcommon_pb.Decimal64): ClobQuote;
  hasLastprice(): boolean;
  clearLastprice(): ClobQuote;

  getLastquantity(): modelcommon_pb.Decimal64 | undefined;
  setLastquantity(value?: modelcommon_pb.Decimal64): ClobQuote;
  hasLastquantity(): boolean;
  clearLastquantity(): ClobQuote;

  getTradedvolume(): modelcommon_pb.Decimal64 | undefined;
  setTradedvolume(value?: modelcommon_pb.Decimal64): ClobQuote;
  hasTradedvolume(): boolean;
  clearTradedvolume(): ClobQuote;

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
    streaminterrupted: boolean,
    streamstatusmsg: string,
    lastprice?: modelcommon_pb.Decimal64.AsObject,
    lastquantity?: modelcommon_pb.Decimal64.AsObject,
    tradedvolume?: modelcommon_pb.Decimal64.AsObject,
  }
}

