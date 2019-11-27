import * as jspb from "google-protobuf"

import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as ordermodel_pb from './ordermodel_pb';
import * as common_pb from './common_pb';

export class CreateAndRouteOrderParams extends jspb.Message {
  getSide(): ordermodel_pb.Side;
  setSide(value: ordermodel_pb.Side): void;

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

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAndRouteOrderParams.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAndRouteOrderParams): CreateAndRouteOrderParams.AsObject;
  static serializeBinaryToWriter(message: CreateAndRouteOrderParams, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAndRouteOrderParams;
  static deserializeBinaryFromReader(message: CreateAndRouteOrderParams, reader: jspb.BinaryReader): CreateAndRouteOrderParams;
}

export namespace CreateAndRouteOrderParams {
  export type AsObject = {
    side: ordermodel_pb.Side,
    quantity?: common_pb.Decimal64.AsObject,
    price?: common_pb.Decimal64.AsObject,
    listingid: string,
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

