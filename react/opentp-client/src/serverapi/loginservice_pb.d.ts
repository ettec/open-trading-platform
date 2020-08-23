import * as jspb from 'google-protobuf'



export class LoginParams extends jspb.Message {
  getUser(): string;
  setUser(value: string): LoginParams;

  getPassword(): string;
  setPassword(value: string): LoginParams;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginParams.AsObject;
  static toObject(includeInstance: boolean, msg: LoginParams): LoginParams.AsObject;
  static serializeBinaryToWriter(message: LoginParams, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginParams;
  static deserializeBinaryFromReader(message: LoginParams, reader: jspb.BinaryReader): LoginParams;
}

export namespace LoginParams {
  export type AsObject = {
    user: string,
    password: string,
  }
}

export class Token extends jspb.Message {
  getToken(): string;
  setToken(value: string): Token;

  getDesk(): string;
  setDesk(value: string): Token;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Token.AsObject;
  static toObject(includeInstance: boolean, msg: Token): Token.AsObject;
  static serializeBinaryToWriter(message: Token, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Token;
  static deserializeBinaryFromReader(message: Token, reader: jspb.BinaryReader): Token;
}

export namespace Token {
  export type AsObject = {
    token: string,
    desk: string,
  }
}

