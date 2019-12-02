import * as jspb from "google-protobuf"

import * as order_pb from './order_pb';
import * as common_pb from './common_pb';

export class SubscribeToOrders extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscribeToOrders.AsObject;
  static toObject(includeInstance: boolean, msg: SubscribeToOrders): SubscribeToOrders.AsObject;
  static serializeBinaryToWriter(message: SubscribeToOrders, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubscribeToOrders;
  static deserializeBinaryFromReader(message: SubscribeToOrders, reader: jspb.BinaryReader): SubscribeToOrders;
}

export namespace SubscribeToOrders {
  export type AsObject = {
  }
}

