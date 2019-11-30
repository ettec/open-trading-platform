import * as jspb from "google-protobuf"

import * as common_pb from './common_pb';

export class Order extends jspb.Message {
  getVersion(): number;
  setVersion(value: number): void;

  getId(): string;
  setId(value: string): void;

  getSide(): Side;
  setSide(value: Side): void;

  getQuantity(): common_pb.Decimal64 | undefined;
  setQuantity(value?: common_pb.Decimal64): void;
  hasQuantity(): boolean;
  clearQuantity(): void;

  getPrice(): common_pb.Decimal64 | undefined;
  setPrice(value?: common_pb.Decimal64): void;
  hasPrice(): boolean;
  clearPrice(): void;

  getListingid(): string;
  setListingid(value: string): void;

  getRemainingquantity(): common_pb.Decimal64 | undefined;
  setRemainingquantity(value?: common_pb.Decimal64): void;
  hasRemainingquantity(): boolean;
  clearRemainingquantity(): void;

  getTradedquantity(): common_pb.Decimal64 | undefined;
  setTradedquantity(value?: common_pb.Decimal64): void;
  hasTradedquantity(): boolean;
  clearTradedquantity(): void;

  getAvgtradeprice(): common_pb.Decimal64 | undefined;
  setAvgtradeprice(value?: common_pb.Decimal64): void;
  hasAvgtradeprice(): boolean;
  clearAvgtradeprice(): void;

  getStatus(): OrderStatus;
  setStatus(value: OrderStatus): void;

  getTargetstatus(): OrderStatus;
  setTargetstatus(value: OrderStatus): void;

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
    quantity?: common_pb.Decimal64.AsObject,
    price?: common_pb.Decimal64.AsObject,
    listingid: string,
    remainingquantity?: common_pb.Decimal64.AsObject,
    tradedquantity?: common_pb.Decimal64.AsObject,
    avgtradeprice?: common_pb.Decimal64.AsObject,
    status: OrderStatus,
    targetstatus: OrderStatus,
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