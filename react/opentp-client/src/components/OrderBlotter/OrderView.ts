import { toNumber } from '../../common/decimal64Conversion';
import { Destinations } from '../../common/destinations';
import { roundToTick } from '../../common/modelutilities';
import { getStrategyDisplayName } from '../../common/strategydescriptions';
import { Listing } from '../../serverapi/listing_pb';
import { Order, OrderStatus, Side } from '../../serverapi/order_pb';
import { ListingService } from '../../services/ListingService';
import { VwapParameters } from '../OrderTicket/Strategies/VwapParams/VwapParamsPanel';


export class OrdersView {

  orders: Map<String, Order> = new Map()
  view?: OrderView[]
  listingSvc: ListingService
  updateListener: () => void
  requestedListingIds: Set<number> = new Set<number>()

  constructor(listingSvc: ListingService, updateListener: () => void) {
    this.listingSvc = listingSvc
    this.updateListener = updateListener
  }

  clear(): void {
    this.orders = new Map()
    this.view = undefined
  }

  addOrUpdateOrder(order: Order): void {
    this.orders.set(order.getId(), order)
    this.view = undefined
  }

  getOrders(): OrderView[] {

    if (!this.view) {

      let result = new Array<Order>()

      for (let order of this.orders.values()) {
        result.push(order)
      }

      this.view = result.map((o: Order) => {
        let v = new OrderView(o)
        let listing = this.listingSvc.GetListingImmediate(v.listingId)
        if (!listing) {

          if (!this.requestedListingIds.has(v.listingId)) {
            this.requestedListingIds.add(v.listingId)
            this.listingSvc.GetListing(v.listingId, (listing: Listing) => {
              v.setListing(listing)
              this.view = undefined
              this.updateListener()
            })
          }
        } else {
          v.setListing(listing)
        }

        return v
      })
    }

    return this.view
  }
}

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
  private listing?: Listing;
  created?: string;
  destination: string;
  owner: string;
  createdBy: string;
  errorMsg: string;
  parameters: string;

  symbol?: string;
  mic?: string;
  countryCode?: string;




  constructor(order: Order) {
    this.id = ""
    this.version = 0;
    this.side = "";
    this.listingId = 0;
    this.status = "";
    this.targetStatus = "";
    this.order = order
    this.created = undefined
    this.destination = "";
    this.owner = "";
    this.errorMsg = "";
    this.createdBy = "";
    this.parameters = "";
    this.setOrder(order)
  }

  public getListing(): Listing | undefined {
    return this.listing
  }

  public setListing(listing: Listing) {
    this.listing = listing
    this.symbol = listing?.getInstrument()?.getDisplaysymbol()
    this.mic = listing?.getMarket()?.getMic()
    this.countryCode = listing?.getMarket()?.getCountrycode()
    this.updateAvgPrice() 
  }

  private updateAvgPrice() {
    if (this.order && this.listing) {
      let avgPrice = toNumber(this.order.getAvgtradeprice())

      if (avgPrice) {
        this.avgTradePrice = roundToTick(avgPrice, this.listing)
      }
    }
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
    if (created) {
      let date = new Date(created.getSeconds() * 1000)
      this.created = date.toLocaleTimeString()
    }

    this.destination = order.getDestination()
    this.owner = order.getOwnerid()
    this.errorMsg = order.getErrormessage()
    this.createdBy = order.getRootoriginatorref()
    this.parameters = this.getParametersDisplayString(order)
    this.updateAvgPrice() 
  }

  getOrder(): Order {
    return this.order
  }

  getParametersDisplayString(order: Order): string {
    if (order.getDestination() !== "" && order.getExecparametersjson() !== "") {


      switch (order.getDestination()) {
        case Destinations.VWAP:
          // order.getDestination() + ":" + order.getExecparametersjson()
          let p = VwapParameters.fromJsonString(order.getExecparametersjson()) as VwapParameters
          return p.toDisplayString()
      }
    }


    return ""
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

  getDestination(): string | undefined {
    if (this.listing?.getMarket()?.getMic()) {
      if (this.destination === this.listing?.getMarket()?.getMic()) {
        return this.listing?.getMarket()?.getName()
      }
    }

    return getStrategyDisplayName(this.destination)
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