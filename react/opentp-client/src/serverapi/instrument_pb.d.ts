import * as jspb from 'google-protobuf'



export class Instrument extends jspb.Message {
  getVersion(): number;
  setVersion(value: number): Instrument;

  getId(): number;
  setId(value: number): Instrument;

  getName(): string;
  setName(value: string): Instrument;

  getDisplaysymbol(): string;
  setDisplaysymbol(value: string): Instrument;

  getEnabled(): boolean;
  setEnabled(value: boolean): Instrument;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Instrument.AsObject;
  static toObject(includeInstance: boolean, msg: Instrument): Instrument.AsObject;
  static serializeBinaryToWriter(message: Instrument, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Instrument;
  static deserializeBinaryFromReader(message: Instrument, reader: jspb.BinaryReader): Instrument;
}

export namespace Instrument {
  export type AsObject = {
    version: number,
    id: number,
    name: string,
    displaysymbol: string,
    enabled: boolean,
  }
}

