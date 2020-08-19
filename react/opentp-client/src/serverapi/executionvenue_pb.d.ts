import * as jspb from "google-protobuf"

import * as order_pb from './order_pb';
import * as modelcommon_pb from './modelcommon_pb';

export class CreateAndRouteOrderParams extends jspb.Message {
  getOrderside(): order_pb.Side;
  setOrderside(value: order_pb.Side): void;

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

  getDestination(): string;
  setDestination(value: string): void;

  getOriginatorid(): string;
  setOriginatorid(value: string): void;

  getOriginatorref(): string;
  setOriginatorref(value: string): void;

  getRootoriginatorid(): string;
  setRootoriginatorid(value: string): void;

  getRootoriginatorref(): string;
  setRootoriginatorref(value: string): void;

  getExecparametersjson(): string;
  setExecparametersjson(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAndRouteOrderParams.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAndRouteOrderParams): CreateAndRouteOrderParams.AsObject;
  static serializeBinaryToWriter(message: CreateAndRouteOrderParams, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAndRouteOrderParams;
  static deserializeBinaryFromReader(message: CreateAndRouteOrderParams, reader: jspb.BinaryReader): CreateAndRouteOrderParams;
}

export namespace CreateAndRouteOrderParams {
  export type AsObject = {
    orderside: order_pb.Side,
    quantity?: modelcommon_pb.Decimal64.AsObject,
    price?: modelcommon_pb.Decimal64.AsObject,
    listingid: number,
    destination: string,
    originatorid: string,
    originatorref: string,
    rootoriginatorid: string,
    rootoriginatorref: string,
    execparametersjson: string,
  }
}

export class OrderId extends jspb.Message {
  getOrderid(): string;
  setOrderid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OrderId.AsObject;
  static toObject(includeInstance: boolean, msg: OrderId): OrderId.AsObject;
  static serializeBinaryToWriter(message: OrderId, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OrderId;
  static deserializeBinaryFromReader(message: OrderId, reader: jspb.BinaryReader): OrderId;
}

export namespace OrderId {
  export type AsObject = {
    orderid: string,
  }
}

export class ExecParamsMetaDataJson extends jspb.Message {
  getJson(): string;
  setJson(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExecParamsMetaDataJson.AsObject;
  static toObject(includeInstance: boolean, msg: ExecParamsMetaDataJson): ExecParamsMetaDataJson.AsObject;
  static serializeBinaryToWriter(message: ExecParamsMetaDataJson, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExecParamsMetaDataJson;
  static deserializeBinaryFromReader(message: ExecParamsMetaDataJson, reader: jspb.BinaryReader): ExecParamsMetaDataJson;
}

export namespace ExecParamsMetaDataJson {
  export type AsObject = {
    json: string,
  }
}

export class CancelOrderParams extends jspb.Message {
  getOrderid(): string;
  setOrderid(value: string): void;

  getListingid(): number;
  setListingid(value: number): void;

  getOwnerid(): string;
  setOwnerid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CancelOrderParams.AsObject;
  static toObject(includeInstance: boolean, msg: CancelOrderParams): CancelOrderParams.AsObject;
  static serializeBinaryToWriter(message: CancelOrderParams, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CancelOrderParams;
  static deserializeBinaryFromReader(message: CancelOrderParams, reader: jspb.BinaryReader): CancelOrderParams;
}

export namespace CancelOrderParams {
  export type AsObject = {
    orderid: string,
    listingid: number,
    ownerid: string,
  }
}

export class ModifyOrderParams extends jspb.Message {
  getOrderid(): string;
  setOrderid(value: string): void;

  getListingid(): number;
  setListingid(value: number): void;

  getOwnerid(): string;
  setOwnerid(value: string): void;

  getQuantity(): modelcommon_pb.Decimal64 | undefined;
  setQuantity(value?: modelcommon_pb.Decimal64): void;
  hasQuantity(): boolean;
  clearQuantity(): void;

  getPrice(): modelcommon_pb.Decimal64 | undefined;
  setPrice(value?: modelcommon_pb.Decimal64): void;
  hasPrice(): boolean;
  clearPrice(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ModifyOrderParams.AsObject;
  static toObject(includeInstance: boolean, msg: ModifyOrderParams): ModifyOrderParams.AsObject;
  static serializeBinaryToWriter(message: ModifyOrderParams, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ModifyOrderParams;
  static deserializeBinaryFromReader(message: ModifyOrderParams, reader: jspb.BinaryReader): ModifyOrderParams;
}

export namespace ModifyOrderParams {
  export type AsObject = {
    orderid: string,
    listingid: number,
    ownerid: string,
    quantity?: modelcommon_pb.Decimal64.AsObject,
    price?: modelcommon_pb.Decimal64.AsObject,
  }
}

