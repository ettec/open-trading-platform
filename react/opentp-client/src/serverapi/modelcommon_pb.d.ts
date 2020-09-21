import * as jspb from 'google-protobuf'



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
  setMantissa(value: number): Decimal64;

  getExponent(): number;
  setExponent(value: number): Decimal64;

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

export class Timestamp extends jspb.Message {
  getSeconds(): number;
  setSeconds(value: number): Timestamp;

  getNanoseconds(): number;
  setNanoseconds(value: number): Timestamp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Timestamp.AsObject;
  static toObject(includeInstance: boolean, msg: Timestamp): Timestamp.AsObject;
  static serializeBinaryToWriter(message: Timestamp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Timestamp;
  static deserializeBinaryFromReader(message: Timestamp, reader: jspb.BinaryReader): Timestamp;
}

export namespace Timestamp {
  export type AsObject = {
    seconds: number,
    nanoseconds: number,
  }
}

