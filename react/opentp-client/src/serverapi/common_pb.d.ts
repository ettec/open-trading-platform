import * as jspb from "google-protobuf"

export class Empty extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Empty.AsObject;
  static toObject(includeInstance: boolean, msg: Empty): Empty.AsObject;
  static serializeBinaryToWriter(message: Empty, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Empty;
  static deserializeBinaryFromReader(message: Empty, reader: jspb.BinaryReader): Empty;
}

export namespace Empty {
  export type AsObject = {
  }
}

export class Decimal64 extends jspb.Message {
  getMantissa(): number;
  setMantissa(value: number): void;

  getExponent(): number;
  setExponent(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Decimal64.AsObject;
  static toObject(includeInstance: boolean, msg: Decimal64): Decimal64.AsObject;
  static serializeBinaryToWriter(message: Decimal64, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Decimal64;
  static deserializeBinaryFromReader(message: Decimal64, reader: jspb.BinaryReader): Decimal64;
}

export namespace Decimal64 {
  export type AsObject = {
    mantissa: number,
    exponent: number,
  }
}
