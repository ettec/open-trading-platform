import * as jspb from 'google-protobuf'



export class Market extends jspb.Message {
  getVersion(): number;
  setVersion(value: number): Market;

  getId(): number;
  setId(value: number): Market;

  getName(): string;
  setName(value: string): Market;

  getMic(): string;
  setMic(value: string): Market;

  getCountrycode(): string;
  setCountrycode(value: string): Market;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Market.AsObject;
  static toObject(includeInstance: boolean, msg: Market): Market.AsObject;
  static serializeBinaryToWriter(message: Market, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Market;
  static deserializeBinaryFromReader(message: Market, reader: jspb.BinaryReader): Market;
}

export namespace Market {
  export type AsObject = {
    version: number,
    id: number,
    name: string,
    mic: string,
    countrycode: string,
  }
}

