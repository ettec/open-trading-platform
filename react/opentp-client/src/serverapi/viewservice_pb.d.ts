import * as jspb from "google-protobuf"

import * as modelcommon_pb from './modelcommon_pb';
import * as order_pb from './order_pb';

export class SubscribeToOrdersWithRootOriginatorIdArgs extends jspb.Message {
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

export class OrderUpdate extends jspb.Message {
  getOrder(): order_pb.Order | undefined;
  setOrder(value?: order_pb.Order): void;
  hasOrder(): boolean;
  clearOrder(): void;

  getTime(): modelcommon_pb.Timestamp | undefined;
  setTime(value?: modelcommon_pb.Timestamp): void;
  hasTime(): boolean;
  clearTime(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrderUpdate.AsObject;
  static toObject(includeInstance: boolean, msg: OrderUpdate): OrderUpdate.AsObject;
  static serializeBinaryToWriter(message: OrderUpdate, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrderUpdate;
  static deserializeBinaryFromReader(message: OrderUpdate, reader: jspb.BinaryReader): OrderUpdate;
}

export namespace OrderUpdate {
  export type AsObject = {
    order?: order_pb.Order.AsObject,
    time?: modelcommon_pb.Timestamp.AsObject,
  }
}

export class OrderHistory extends jspb.Message {
  getUpdatesList(): Array<OrderUpdate>;
  setUpdatesList(value: Array<OrderUpdate>): void;
  clearUpdatesList(): void;
  addUpdates(value?: OrderUpdate, index?: number): OrderUpdate;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrderHistory.AsObject;
  static toObject(includeInstance: boolean, msg: OrderHistory): OrderHistory.AsObject;
  static serializeBinaryToWriter(message: OrderHistory, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrderHistory;
  static deserializeBinaryFromReader(message: OrderHistory, reader: jspb.BinaryReader): OrderHistory;
}

export namespace OrderHistory {
  export type AsObject = {
    updatesList: Array<OrderUpdate.AsObject>,
  }
}

