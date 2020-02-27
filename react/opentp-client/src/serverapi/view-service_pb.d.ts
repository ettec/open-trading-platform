import * as jspb from "google-protobuf"

import * as modelcommon_pb from './modelcommon_pb';
import * as order_pb from './order_pb';

export class SubscribeToOrders extends jspb.Message {
  getAfter(): modelcommon_pb.Timestamp | undefined;
  setAfter(value?: modelcommon_pb.Timestamp): void;
  hasAfter(): boolean;
  clearAfter(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscribeToOrders.AsObject;
  static toObject(includeInstance: boolean, msg: SubscribeToOrders): SubscribeToOrders.AsObject;
  static serializeBinaryToWriter(message: SubscribeToOrders, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubscribeToOrders;
  static deserializeBinaryFromReader(message: SubscribeToOrders, reader: jspb.BinaryReader): SubscribeToOrders;
}

export namespace SubscribeToOrders {
  export type AsObject = {
    after?: modelcommon_pb.Timestamp.AsObject,
  }
}

