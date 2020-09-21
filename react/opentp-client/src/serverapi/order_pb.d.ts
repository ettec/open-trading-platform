import * as jspb from 'google-protobuf'

import * as modelcommon_pb from './modelcommon_pb';


export class Ref extends jspb.Message {
  getVersion(): number;
  setVersion(value: number): Ref;

  getId(): string;
  setId(value: string): Ref;

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
  setVersion(value: number): Order;

  getId(): string;
  setId(value: string): Order;

  getSide(): Side;
  setSide(value: Side): Order;

  getQuantity(): modelcommon_pb.Decimal64 | undefined;
  setQuantity(value?: modelcommon_pb.Decimal64): Order;
  hasQuantity(): boolean;
  clearQuantity(): Order;

  getPrice(): modelcommon_pb.Decimal64 | undefined;
  setPrice(value?: modelcommon_pb.Decimal64): Order;
  hasPrice(): boolean;
  clearPrice(): Order;

  getListingid(): number;
  setListingid(value: number): Order;

  getRemainingquantity(): modelcommon_pb.Decimal64 | undefined;
  setRemainingquantity(value?: modelcommon_pb.Decimal64): Order;
  hasRemainingquantity(): boolean;
  clearRemainingquantity(): Order;

  getTradedquantity(): modelcommon_pb.Decimal64 | undefined;
  setTradedquantity(value?: modelcommon_pb.Decimal64): Order;
  hasTradedquantity(): boolean;
  clearTradedquantity(): Order;

  getAvgtradeprice(): modelcommon_pb.Decimal64 | undefined;
  setAvgtradeprice(value?: modelcommon_pb.Decimal64): Order;
  hasAvgtradeprice(): boolean;
  clearAvgtradeprice(): Order;

  getStatus(): OrderStatus;
  setStatus(value: OrderStatus): Order;

  getTargetstatus(): OrderStatus;
  setTargetstatus(value: OrderStatus): Order;

  getCreated(): modelcommon_pb.Timestamp | undefined;
  setCreated(value?: modelcommon_pb.Timestamp): Order;
  hasCreated(): boolean;
  clearCreated(): Order;

  getOwnerid(): string;
  setOwnerid(value: string): Order;

  getOriginatorid(): string;
  setOriginatorid(value: string): Order;

  getOriginatorref(): string;
  setOriginatorref(value: string): Order;

  getLastexecquantity(): modelcommon_pb.Decimal64 | undefined;
  setLastexecquantity(value?: modelcommon_pb.Decimal64): Order;
  hasLastexecquantity(): boolean;
  clearLastexecquantity(): Order;

  getLastexecprice(): modelcommon_pb.Decimal64 | undefined;
  setLastexecprice(value?: modelcommon_pb.Decimal64): Order;
  hasLastexecprice(): boolean;
  clearLastexecprice(): Order;

  getLastexecid(): string;
  setLastexecid(value: string): Order;

  getExposedquantity(): modelcommon_pb.Decimal64 | undefined;
  setExposedquantity(value?: modelcommon_pb.Decimal64): Order;
  hasExposedquantity(): boolean;
  clearExposedquantity(): Order;

  getErrormessage(): string;
  setErrormessage(value: string): Order;

  getChildordersrefsList(): Array<Ref>;
  setChildordersrefsList(value: Array<Ref>): Order;
  clearChildordersrefsList(): Order;
  addChildordersrefs(value?: Ref, index?: number): Ref;

  getRootoriginatorid(): string;
  setRootoriginatorid(value: string): Order;

  getRootoriginatorref(): string;
  setRootoriginatorref(value: string): Order;

  getExecparametersjson(): string;
  setExecparametersjson(value: string): Order;

  getDestination(): string;
  setDestination(value: string): Order;

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
    rootoriginatorid: string,
    rootoriginatorref: string,
    execparametersjson: string,
    destination: string,
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
