import * as jspb from "google-protobuf"

import * as instrument_pb from './instrument_pb';
import * as market_pb from './market_pb';

export class Listing extends jspb.Message {
  getVersion(): number;
  setVersion(value: number): void;

  getId(): number;
  setId(value: number): void;

  getMarket(): market_pb.Market | undefined;
  setMarket(value?: market_pb.Market): void;
  hasMarket(): boolean;
  clearMarket(): void;

  getInstrument(): instrument_pb.Instrument | undefined;
  setInstrument(value?: instrument_pb.Instrument): void;
  hasInstrument(): boolean;
  clearInstrument(): void;

  getMarketsymbol(): string;
  setMarketsymbol(value: string): void;

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
  }
}

