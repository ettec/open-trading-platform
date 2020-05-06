import * as jspb from "google-protobuf"

import * as modelcommon_pb from './modelcommon_pb';
import * as order_pb from './order_pb';

export class SubscribeToOrdersWithRootOriginatorIdArgs extends jspb.Message {
  getAfter(): modelcommon_pb.Timestamp | undefined;
  setAfter(value?: modelcommon_pb.Timestamp): void;
  hasAfter(): boolean;
  clearAfter(): void;

  getRootoriginatorid(): string;
  setRootoriginatorid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubscribeToOrdersWithRootOriginatorIdArgs.AsObject;
  static toObject(includeInstance: boolean, msg: SubscribeToOrdersWithRootOriginatorIdArgs): SubscribeToOrdersWithRootOriginatorIdArgs.AsObject;
  static serializeBinaryToWriter(message: SubscribeToOrdersWithRootOriginatorIdArgs, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubscribeToOrdersWithRootOriginatorIdArgs;
  static deserializeBinaryFromReader(message: SubscribeToOrdersWithRootOriginatorIdArgs, reader: jspb.BinaryReader): SubscribeToOrdersWithRootOriginatorIdArgs;
}

export namespace SubscribeToOrdersWithRootOriginatorIdArgs {
  export type AsObject = {
    after?: modelcommon_pb.Timestamp.AsObject,
    rootoriginatorid: string,
  }
}

export class GetOrderHistoryArgs extends jspb.Message {
  getOrderid(): string;
  setOrderid(value: string): void;

  getToversion(): number;
  setToversion(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrderHistoryArgs.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrderHistoryArgs): GetOrderHistoryArgs.AsObject;
  static serializeBinaryToWriter(message: GetOrderHistoryArgs, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrderHistoryArgs;
  static deserializeBinaryFromReader(message: GetOrderHistoryArgs, reader: jspb.BinaryReader): GetOrderHistoryArgs;
}

export namespace GetOrderHistoryArgs {
  export type AsObject = {
    orderid: string,
    toversion: number,
  }
}

export class Orders extends jspb.Message {
  getOrdersList(): Array<order_pb.Order>;
  setOrdersList(value: Array<order_pb.Order>): void;
  clearOrdersList(): void;
  addOrders(value?: order_pb.Order, index?: number): order_pb.Order;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Orders.AsObject;
  static toObject(includeInstance: boolean, msg: Orders): Orders.AsObject;
  static serializeBinaryToWriter(message: Orders, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Orders;
  static deserializeBinaryFromReader(message: Orders, reader: jspb.BinaryReader): Orders;
}

export namespace Orders {
  export type AsObject = {
    ordersList: Array<order_pb.Order.AsObject>,
  }
}

