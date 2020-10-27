import { toNumber } from '../../common/decimal64Conversion';
import { Destinations } from '../../common/destinations';
import { roundToTick } from '../../common/modelutilities';
import { Listing } from '../../serverapi/listing_pb';
import { Order, OrderStatus, Side } from '../../serverapi/order_pb';
import { ListingService } from '../../services/ListingService';
import { VwapParameters } from '../OrderTicket/Strategies/VwapParams/VwapParamsPanel';


export class OrdersView {

  viewIdx: Map<String, number> = new Map()
  listingSvc: ListingService
  updateListener: (orderView: OrderView) => void
  views: Array<OrderView>
  orderViewUpdateListener: (orderView: OrderView) => void

  constructor(listingSvc: ListingService, updateListener: (orderView: OrderView) => void) {
    this.listingSvc = listingSvc
    this.updateListener = updateListener
    this.views = new Array<OrderView>()

    this.orderViewUpdateListener = (orderView: OrderView) => {
      this.updateListener(orderView)
    }
  }

  updateView(order: Order) {

    let idx = this.viewIdx.get(order.getId())
    var view : OrderView;
    if (idx) {
      view = this.views[idx] 
      view.setOrder(order)
    } else {
      view = new OrderView(order, this.listingSvc.GetListing, this.orderViewUpdateListener)
      this.viewIdx.set(view.id, this.views.length)
      this.views.push(view)
    }

     this.updateListener(view) 
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

  updateListener: (orderView: OrderView) => void




  constructor(order: Order, getListing: (listingId: number, callback: (listing: Listing) => void) => void, updateListener: (orderView: OrderView) => void) {
    this.updateListener = updateListener

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

    this.setListing = this.setListing.bind(this)
    
      getListing(this.listingId, (listing: Listing) => {
        this.setListing(listing)
        this.updateListener(this)
      })
  
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

  setOrder(order: Order) {
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

  private getParametersDisplayString(order: Order): string {
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

  private getStatusString(status: OrderStatus) {

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

}