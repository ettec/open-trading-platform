import * as jspb from "google-protobuf"

import * as listing_pb from './listing_pb';
import * as order_pb from './order_pb';
import * as common_pb from './common_pb';

export class CreateAndRouteOrderParams extends jspb.Message {
  getSide(): order_pb.Side;
  setSide(value: order_pb.Side): void;

  getQuantity(): common_pb.Decimal64 | undefined;
  setQuantity(value?: common_pb.Decimal64): void;
  hasQuantity(): boolean;
  clearQuantity(): void;

  getPrice(): common_pb.Decimal64 | undefined;
  setPrice(value?: common_pb.Decimal64): void;
  hasPrice(): boolean;
  clearPrice(): void;

  getListing(): listing_pb.Listing | undefined;
  setListing(value?: listing_pb.Listing): void;
  hasListing(): boolean;
  clearListing(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateAndRouteOrderParams.AsObject;
  static toObject(includeInstance: boolean, msg: CreateAndRouteOrderParams): CreateAndRouteOrderParams.AsObject;
  static serializeBinaryToWriter(message: CreateAndRouteOrderParams, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateAndRouteOrderParams;
  static deserializeBinaryFromReader(message: CreateAndRouteOrderParams, reader: jspb.BinaryReader): CreateAndRouteOrderParams;
}

export namespace CreateAndRouteOrderParams {
  export type AsObject = {
    side: order_pb.Side,
    quantity?: common_pb.Decimal64.AsObject,
    price?: common_pb.Decimal64.AsObject,
    listing?: listing_pb.Listing.AsObject,
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

