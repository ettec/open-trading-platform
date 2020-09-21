import { toNumber } from '../../util/decimal64Conversion';
import { Order, Side, OrderStatus } from '../../serverapi/order_pb';
import { Listing } from '../../serverapi/listing_pb';
import { ListingService } from '../../services/ListingService';

export interface Filter {
  id(): string
  exclude(order: Order): boolean
}

export interface Sort {
  compare(a: Order, b: Order): number
}

export class OrdersView {

  orders: Map<String, Order> = new Map()
  filter?: (order: Order, index: number, array: Order[]) => boolean
  sort?: (a: Order, b: Order) => number
  view?: OrderView[]
  listingSvc: ListingService
  updateListener: () => void
  requestedListingIds: Set<number> = new Set<number>()

  constructor(listingSvc: ListingService, updateListener: () => void) {
    this.listingSvc = listingSvc
    this.updateListener = updateListener
  }

  setFilter(f?: (order: Order, index: number, array: Order[]) => boolean) {
    this.filter = f
    this.view = undefined
  }

  setSort(s?: (a: Order, b: Order) => number) {
    this.sort = s
    this.view = undefined
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

      if (this.filter) {
        result = result.filter(this.filter)
      }

      if (this.sort) {
        result.sort(this.sort)
      }

      this.view = result.map((o: Order) => {
        let v = new OrderView(o)
        v.listing = this.listingSvc.GetListingImmediate(v.listingId)
        if (!v.listing) {

          if (!this.requestedListingIds.has(v.listingId)) {
            this.requestedListingIds.add(v.listingId)
            this.listingSvc.GetListing(v.listingId, (listing: Listing) => {
              v.listing = listing
              this.view = undefined
              this.updateListener()
            })
          }
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
  listing?: Listing;
  created?: Date;
  destination: string;
  owner: string;

  errorMsg: string;

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
    if (created) {
      this.created = new Date(created.getSeconds() * 1000)
    }

    this.destination = order.getDestination()
    this.owner = order.getOwnerid()
    this.errorMsg = order.getErrormessage()

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