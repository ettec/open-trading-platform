import * as jspb from 'google-protobuf'

import * as instrument_pb from './instrument_pb';
import * as market_pb from './market_pb';
import * as modelcommon_pb from './modelcommon_pb';


export class Listing extends jspb.Message {
  getVersion(): number;
  setVersion(value: number): Listing;

  getId(): number;
  setId(value: number): Listing;

  getMarket(): market_pb.Market | undefined;
  setMarket(value?: market_pb.Market): Listing;
  hasMarket(): boolean;
  clearMarket(): Listing;

  getInstrument(): instrument_pb.Instrument | undefined;
  setInstrument(value?: instrument_pb.Instrument): Listing;
  hasInstrument(): boolean;
  clearInstrument(): Listing;

  getMarketsymbol(): string;
  setMarketsymbol(value: string): Listing;

  getTicksize(): TickSizeTable | undefined;
  setTicksize(value?: TickSizeTable): Listing;
  hasTicksize(): boolean;
  clearTicksize(): Listing;

  getSizeincrement(): modelcommon_pb.Decimal64 | undefined;
  setSizeincrement(value?: modelcommon_pb.Decimal64): Listing;
  hasSizeincrement(): boolean;
  clearSizeincrement(): Listing;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Listing.AsObject;
  static toObject(includeInstance: boolean, msg: Listing): Listing.AsObject;
  static serializeBinaryToWriter(message: Listing, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Listing;
  static deserializeBinaryFromReader(message: Listing, reader: jspb.BinaryReader): Listing;
}

export namespace Listing {
  export type AsObject = {
    version: number,
    id: number,
    market?: market_pb.Market.AsObject,
    instrument?: instrument_pb.Instrument.AsObject,
    marketsymbol: string,
    ticksize?: TickSizeTable.AsObject,
    sizeincrement?: modelcommon_pb.Decimal64.AsObject,
  }
}

export class TickSizeTable extends jspb.Message {
  getEntriesList(): Array<TickSizeEntry>;
  setEntriesList(value: Array<TickSizeEntry>): TickSizeTable;
  clearEntriesList(): TickSizeTable;
  addEntries(value?: TickSizeEntry, index?: number): TickSizeEntry;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TickSizeTable.AsObject;
  static toObject(includeInstance: boolean, msg: TickSizeTable): TickSizeTable.AsObject;
  static serializeBinaryToWriter(message: TickSizeTable, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TickSizeTable;
  static deserializeBinaryFromReader(message: TickSizeTable, reader: jspb.BinaryReader): TickSizeTable;
}

export namespace TickSizeTable {
  export type AsObject = {
    entriesList: Array<TickSizeEntry.AsObject>,
  }
}

export class TickSizeEntry extends jspb.Message {
  getLowerpricebound(): modelcommon_pb.Decimal64 | undefined;
  setLowerpricebound(value?: modelcommon_pb.Decimal64): TickSizeEntry;
  hasLowerpricebound(): boolean;
  clearLowerpricebound(): TickSizeEntry;

  getUpperpricebound(): modelcommon_pb.Decimal64 | undefined;
  setUpperpricebound(value?: modelcommon_pb.Decimal64): TickSizeEntry;
  hasUpperpricebound(): boolean;
  clearUpperpricebound(): TickSizeEntry;

  getTicksize(): modelcommon_pb.Decimal64 | undefined;
  setTicksize(value?: modelcommon_pb.Decimal64): TickSizeEntry;
  hasTicksize(): boolean;
  clearTicksize(): TickSizeEntry;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TickSizeEntry.AsObject;
  static toObject(includeInstance: boolean, msg: TickSizeEntry): TickSizeEntry.AsObject;
  static serializeBinaryToWriter(message: TickSizeEntry, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TickSizeEntry;
  static deserializeBinaryFromReader(message: TickSizeEntry, reader: jspb.BinaryReader): TickSizeEntry;
}

export namespace TickSizeEntry {
  export type AsObject = {
    lowerpricebound?: modelcommon_pb.Decimal64.AsObject,
    upperpricebound?: modelcommon_pb.Decimal64.AsObject,
    ticksize?: modelcommon_pb.Decimal64.AsObject,
  }
}

