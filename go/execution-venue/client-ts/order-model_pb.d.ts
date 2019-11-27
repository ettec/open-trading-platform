import * as jspb from "google-protobuf"

import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';

export class Order extends jspb.Message {
  getSide(): Side;
  setSide(value: Side): void;

  getQuantity(): number;
  setQuantity(value: number): void;

  getPrice(): number;
  setPrice(value: number): void;

  getListingid(): string;
  setListingid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Order.AsObject;
  static toObject(includeInstance: boolean, msg: Order): Order.AsObject;
  static serializeBinaryToWriter(message: Order, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Order;
  static deserializeBinaryFromReader(message: Order, reader: jspb.BinaryReader): Order;
}

export namespace Order {
  export type AsObject = {
    side: Side,
    quantity: number,
    price: number,
    listingid: string,
  }
}

export enum Side { 
  BUY = 0,
  SELL = 1,
}
