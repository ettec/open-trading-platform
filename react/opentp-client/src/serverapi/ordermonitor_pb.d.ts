import * as jspb from "google-protobuf"

import * as modelcommon_pb from './modelcommon_pb';

export class CancelAllOrdersForOriginatorIdParams extends jspb.Message {
  getOriginatorid(): string;
  setOriginatorid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CancelAllOrdersForOriginatorIdParams.AsObject;
  static toObject(includeInstance: boolean, msg: CancelAllOrdersForOriginatorIdParams): CancelAllOrdersForOriginatorIdParams.AsObject;
  static serializeBinaryToWriter(message: CancelAllOrdersForOriginatorIdParams, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CancelAllOrdersForOriginatorIdParams;
  static deserializeBinaryFromReader(message: CancelAllOrdersForOriginatorIdParams, reader: jspb.BinaryReader): CancelAllOrdersForOriginatorIdParams;
}

export namespace CancelAllOrdersForOriginatorIdParams {
  export type AsObject = {
    originatorid: string,
  }
}

