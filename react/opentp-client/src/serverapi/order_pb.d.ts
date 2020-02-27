import * as jspb from "google-protobuf"

import * as modelcommon_pb from './modelcommon_pb';

export class Order extends jspb.Message {
  getVersion(): number;
  setVersion(value: number): void;

  getId(): string;
  setId(value: string): void;

  getSide(): Side;
  setSide(value: Side): void;

  getQuantity(): modelcommon_pb.Decimal64 | undefined;
  setQuantity(value?: modelcommon_pb.Decimal64): void;
  hasQuantity(): boolean;
  clearQuantity(): void;

  getPrice(): modelcommon_pb.Decimal64 | undefined;
  setPrice(value?: modelcommon_pb.Decimal64): void;
  hasPrice(): boolean;
  clearPrice(): void;

  getListingid(): number;
  setListingid(value: number): void;

  getRemainingquantity(): modelcommon_pb.Decimal64 | undefined;
  setRemainingquantity(value?: modelcommon_pb.Decimal64): void;
  hasRemainingquantity(): boolean;
  clearRemainingquantity(): void;

  getTradedquantity(): modelcommon_pb.Decimal64 | undefined;
  setTradedquantity(value?: modelcommon_pb.Decimal64): void;
  hasTradedquantity(): boolean;
  clearTradedquantity(): void;

  getAvgtradeprice(): modelcommon_pb.Decimal64 | undefined;
  setAvgtradeprice(value?: modelcommon_pb.Decimal64): void;
  hasAvgtradeprice(): boolean;
  clearAvgtradeprice(): void;

  getStatus(): OrderStatus;
  setStatus(value: OrderStatus): void;

  getTargetstatus(): OrderStatus;
  setTargetstatus(value: OrderStatus): void;

  getCreated(): modelcommon_pb.Timestamp | undefined;
  setCreated(value?: modelcommon_pb.Timestamp): void;
  hasCreated(): boolean;
  clearCreated(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Order.AsObject;
  static toObject(includeInstance: boolean, msg: Order): Order.AsObject;
  static serializeBinaryToWriter(message: Order, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Order;
  static deserializeBinaryFromReader(message: Order, reader: jspb.BinaryReader): Order;
}

export namespace Order {
  export type AsObject = {
    version: number,
    id: string,
    side: Side,
    quantity?: modelcommon_pb.Decimal64.AsObject,
    price?: modelcommon_pb.Decimal64.AsObject,
    listingid: number,
    remainingquantity?: modelcommon_pb.Decimal64.AsObject,
    tradedquantity?: modelcommon_pb.Decimal64.AsObject,
    avgtradeprice?: modelcommon_pb.Decimal64.AsObject,
    status: OrderStatus,
    targetstatus: OrderStatus,
    created?: modelcommon_pb.Timestamp.AsObject,
  }
}

export enum Side { 
  BUY = 0,
  SELL = 1,
}
export enum OrderStatus { 
  NONE = 0,
  LIVE = 1,
  FILLED = 2,
  CANCELLED = 3,
}
