import { toNumber } from '../../util/decimal64Conversion';
import { Order, Side, OrderStatus } from '../../serverapi/order_pb';
import { Listing } from '../../serverapi/listing_pb';

export class OrderView {

    version: number;
    id: string;
    side: string;
    quantity?: number;
    price?: number;
    listingId: number;
    remainingQuantity?: number;
    exposedQuantity?: number;
    tradedQuantity?: number;
    avgTradePrice?: number;
    status: string;
    targetStatus: string;
    private order: Order;
    listing?: Listing;
    created?: Date;
    placedWith: string;
  
    constructor(order: Order) {
      this.id = ""
      this.version = 0;
      this.side = "";
      this.listingId = 0;
      this.status = "";
      this.targetStatus = "";
      this.order = order
      this.created = undefined
      this.placedWith = "";
      this.setOrder(order)
    }
  
    private setOrder(order: Order) {
      this.order = order
  
  
      this.version = order.getVersion()
      this.id = order.getId()
  
      switch (order.getSide()) {
        case Side.BUY:
          this.side = "Buy"
          break;
        case Side.SELL:
          this.side = "Sell"
          break;
        default:
          this.side = "Unknown Side: " + order.getSide()
      }
  
      this.quantity = toNumber(order.getQuantity())
      this.price = toNumber(order.getPrice())
      this.listingId = order.getListingid()
      this.remainingQuantity = toNumber(order.getRemainingquantity())
      this.exposedQuantity = toNumber(order.getExposedquantity())
      this.tradedQuantity = toNumber(order.getTradedquantity())
      this.avgTradePrice = toNumber(order.getAvgtradeprice())
      this.status = this.getStatusString(order.getStatus())
      this.targetStatus = this.getStatusString(order.getTargetstatus())
      let created = order.getCreated()
      if( created ) {
        this.created = new Date(created.getSeconds() * 1000)
      }

      this.placedWith = order.getOwnerid()
      
    }
  
    getOrder(): Order {
      return this.order
    }
  
    getStatusString(status: OrderStatus) {
  
      switch (status) {
        case OrderStatus.CANCELLED:
          return "Cancelled"
        case OrderStatus.FILLED:
          return "Filled"
        case OrderStatus.LIVE:
          return "Live"
        case OrderStatus.NONE:
          return "None"
      }
  
    }
  
    getSymbol(): string | undefined {
      return this.listing?.getInstrument()?.getDisplaysymbol()
    }
  
    getMic(): string | undefined {
      return this.listing?.getMarket()?.getMic()
    }
  
    getCountryCode(): string | undefined {
      return this.listing?.getMarket()?.getCountrycode()
    }
  
  }