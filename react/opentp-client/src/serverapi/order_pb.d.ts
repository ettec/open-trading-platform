import * as jspb from "google-protobuf"

import * as modelcommon_pb from './modelcommon_pb';

export class Ref extends jspb.Message {
  getVersion(): number;
  setVersion(value: number): void;

  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Ref.AsObject;
  static toObject(includeInstance: boolean, msg: Ref): Ref.AsObject;
  static serializeBinaryToWriter(message: Ref, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Ref;
  static deserializeBinaryFromReader(message: Ref, reader: jspb.BinaryReader): Ref;
}

export namespace Ref {
  export type AsObject = {
    version: number,
    id: string,
  }
}

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

  getOwnerid(): string;
  setOwnerid(value: string): void;

  getOriginatorid(): string;
  setOriginatorid(value: string): void;

  getOriginatorref(): string;
  setOriginatorref(value: string): void;

  getLastexecquantity(): modelcommon_pb.Decimal64 | undefined;
  setLastexecquantity(value?: modelcommon_pb.Decimal64): void;
  hasLastexecquantity(): boolean;
  clearLastexecquantity(): void;

  getLastexecprice(): modelcommon_pb.Decimal64 | undefined;
  setLastexecprice(value?: modelcommon_pb.Decimal64): void;
  hasLastexecprice(): boolean;
  clearLastexecprice(): void;

  getLastexecid(): string;
  setLastexecid(value: string): void;

  getExposedquantity(): modelcommon_pb.Decimal64 | undefined;
  setExposedquantity(value?: modelcommon_pb.Decimal64): void;
  hasExposedquantity(): boolean;
  clearExposedquantity(): void;

  getErrormessage(): string;
  setErrormessage(value: string): void;

  getChildordersrefsList(): Array<Ref>;
  setChildordersrefsList(value: Array<Ref>): void;
  clearChildordersrefsList(): void;
  addChildordersrefs(value?: Ref, index?: number): Ref;

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
    ownerid: string,
    originatorid: string,
    originatorref: string,
    lastexecquantity?: modelcommon_pb.Decimal64.AsObject,
    lastexecprice?: modelcommon_pb.Decimal64.AsObject,
    lastexecid: string,
    exposedquantity?: modelcommon_pb.Decimal64.AsObject,
    errormessage: string,
    childordersrefsList: Array<Ref.AsObject>,
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
