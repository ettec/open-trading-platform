

/*
export interface SearchDisplayInstrument {
  id:number,
  name:string,
  symbol:string
}

export interface Instrument {

  id:number,
  canon : {
    name:string,
    type:string,
    symbol:string
  },
  symbols : {
    IEX : string
  }
  
} */

export interface InstrumentWatchLine {
  id:number,
  name:string,
  symbol:string,
  bidSize:string,
  bidPrice:string,
  askPrice:string,
  askSize:string
}

export interface LocalBookLine {
  bidSize:string,
  bidPrice:string,
  askPrice:string,
  askSize:string
}

/*

export enum OrderStatus {
    Live,
    Cancelled,
    Filled,
    None
  }
  
  
  
  export enum Side {
    Buy, Sell
  }
  
  export interface Order {
    id:string,
    qty:number,
    price: number,
    instrumentId:string,
    orderStatus: OrderStatus,
    targetStatus: OrderStatus,
    side: Side
    tradedQty: number,
    remainingQty: number,
    errorMsg: string
  } */