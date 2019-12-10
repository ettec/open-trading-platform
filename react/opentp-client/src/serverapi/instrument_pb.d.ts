import * as jspb from "google-protobuf"

export class Instrument extends jspb.Message {
  getVersion(): number;
  setVersion(value: number): void;

  getId(): number;
  setId(value: number): void;

  getName(): string;
  setName(value: string): void;

  getDisplaysymbol(): string;
  setDisplaysymbol(value: string): void;

  getEnabled(): boolean;
  setEnabled(value: boolean): void;

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

