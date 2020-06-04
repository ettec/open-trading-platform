import * as jspb from "google-protobuf"

import * as modelcommon_pb from './modelcommon_pb';

export class GetConfigParameters extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetConfigParameters.AsObject;
  static toObject(includeInstance: boolean, msg: GetConfigParameters): GetConfigParameters.AsObject;
  static serializeBinaryToWriter(message: GetConfigParameters, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetConfigParameters;
  static deserializeBinaryFromReader(message: GetConfigParameters, reader: jspb.BinaryReader): GetConfigParameters;
}

export namespace GetConfigParameters {
  export type AsObject = {
    userid: string,
  }
}

export class Config extends jspb.Message {
  getConfig(): string;
  setConfig(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Config.AsObject;
  static toObject(includeInstance: boolean, msg: Config): Config.AsObject;
  static serializeBinaryToWriter(message: Config, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Config;
  static deserializeBinaryFromReader(message: Config, reader: jspb.BinaryReader): Config;
}

export namespace Config {
  export type AsObject = {
    config: string,
  }
}

export class StoreConfigParams extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getConfig(): string;
  setConfig(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StoreConfigParams.AsObject;
  static toObject(includeInstance: boolean, msg: StoreConfigParams): StoreConfigParams.AsObject;
  static serializeBinaryToWriter(message: StoreConfigParams, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StoreConfigParams;
  static deserializeBinaryFromReader(message: StoreConfigParams, reader: jspb.BinaryReader): StoreConfigParams;
}

export namespace StoreConfigParams {
  export type AsObject = {
    userid: string,
    config: string,
  }
}

export class Void extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Void.AsObject;
  static toObject(includeInstance: boolean, msg: Void): Void.AsObject;
  static serializeBinaryToWriter(message: Void, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Void;
  static deserializeBinaryFromReader(message: Void, reader: jspb.BinaryReader): Void;
}

export namespace Void {
  export type AsObject = {
  }
}

