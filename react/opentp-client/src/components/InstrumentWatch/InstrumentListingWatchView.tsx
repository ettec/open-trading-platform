import { toNumber } from "../../common/decimal64Conversion";
import { ClobQuote } from "../../serverapi/clobquote_pb";
import { Listing } from "../../serverapi/listing_pb";
import { ListingService } from "../../services/ListingService";
import { QuoteListener, QuoteService } from "../../services/QuoteService";


export enum WatchEventType {
  Add,
  Update,
  Remove
}


export class InstrumentListingWatchesView implements QuoteListener {


  listingSvc: ListingService
  quoteSvc: QuoteService
  updateListener: (orderView: ListingWatchView[], eventType: WatchEventType) => void
  views: Map<number, ListingWatchView> = new Map()

  constructor(listingSvc: ListingService, quoteSvc: QuoteService,
    updateListener: (watch: ListingWatchView[], eventType: WatchEventType) => void) {

    this.listingSvc = listingSvc
    this.quoteSvc = quoteSvc
    this.updateListener = updateListener
  }

  removeListings(listingIds: number[]) {

    let removed = new Array<ListingWatchView>()

    for (let listingId of listingIds) {

      let view = this.views.get(listingId)
      if (view) {
        this.views.delete(listingId)
        this.quoteSvc.UnsubscribeFromQuote(listingId, this)
        removed.push(view)

      }

    }

    this.updateListener(removed, WatchEventType.Remove)
  }


  addListing(listingId: number) {



    if (!this.views.has(listingId)) {
      let view = new ListingWatchView(listingId)
      this.views.set(listingId, view)
      this.updateListener([view], WatchEventType.Add)

      this.listingSvc.GetListing(listingId, (listing: Listing) => {
        view.setListing(listing)
        this.updateListener([view], WatchEventType.Update)
      })

      this.quoteSvc.SubscribeToQuote(listingId, this)

    }

  }

  onQuote(quote: ClobQuote): void {
    let view = this.views.get(quote.getListingid())
    if (view) {
      view.setQuote(quote)
      this.updateListener([view], WatchEventType.Update)
    }
  }



}

export interface PriceUpdate {
  price?: string;
  direction?: number;
}

export class ListingWatchView {


  listing?: Listing;
  quote?: ClobQuote;

  listingId: number;
  lastBidPrice?: number;
  lastAskPrice?: number;
  currentBidPrice?: number;
  currentAskPrice?: number;
  lastLastPrice?: number;
  currentLastPrice?: number;
  lastTradedQty?: number;
  tradedVolume?: number;

  symbol?: string;
  name?: string;
  mic?: string;
  countryCode?: string;

  bidPrice?: PriceUpdate
  askPrice?: PriceUpdate
  lastPrice?: PriceUpdate

  bidSize?: string
  askSize?: string
  lastSize?: string


  constructor(listingId: number) {
    this.listingId = listingId
  }

  getListing(): Listing | undefined {
    return this.listing
  }

  setListing(listing: Listing) {
    this.listing = listing
    this.symbol = listing?.getInstrument()?.getDisplaysymbol()
    this.name = listing?.getInstrument()?.getName()
    this.mic = listing?.getMarket()?.getMic()
    this.countryCode = listing?.getMarket()?.getCountrycode()
  }

  setQuote(quote: ClobQuote) {
    this.quote = quote

    if (this.quote.getBidsList().length >= 1) {
      let depth = this.quote.getBidsList()[0]
      let price = toNumber(depth.getPrice())
      if (price !== this.currentBidPrice) {
        this.lastBidPrice = this.currentBidPrice
        this.currentBidPrice = price
      }


      let direction = undefined
      if (this.currentBidPrice && this.lastBidPrice) {
        direction = this.currentBidPrice - this.lastBidPrice
      }

      if (price) {
        this.bidPrice = { price: price.toString(), direction: direction }
      }


      let sz = toNumber(depth.getSize())
      if (sz) {
        this.bidSize = sz.toString()
      }

    } else {
      this.bidPrice = undefined
      this.bidSize = undefined
    }

    if (this.quote.getOffersList().length >= 1) {
      let depth = this.quote.getOffersList()[0]
      let price = toNumber(depth.getPrice())
      if (price !== this.currentAskPrice) {
        this.lastAskPrice = this.currentAskPrice
        this.currentAskPrice = price
      }


      let direction = undefined
      if (this.currentAskPrice && this.lastAskPrice) {
        direction = this.currentAskPrice - this.lastAskPrice
      }

      if (price) {
        this.askPrice = { price: price.toString(), direction: direction }
      }


      let sz = toNumber(depth.getSize())
      if (sz) {
        this.askSize = sz.toString()
      }

    } else {
      this.askPrice = undefined
      this.askSize = undefined
    }


    if (this.quote.getLastprice()) {

      let price = toNumber(this.quote.getLastprice())
      if (price !== this.currentLastPrice) {
        this.lastLastPrice = this.currentLastPrice
        this.currentLastPrice = price
      }


      if (price !== this.currentLastPrice) {
        this.lastLastPrice = this.currentLastPrice
        this.currentLastPrice = price
      }

      let direction = undefined
      if (this.currentLastPrice && this.lastLastPrice) {

        direction = this.currentLastPrice - this.lastLastPrice
      }


      if (price) {
        this.lastPrice = { price: price.toString(), direction: direction }
      }

      let sz = toNumber(this.quote.getLastquantity())
      if (sz) {
        this.lastSize = sz.toString()
      }

    } else {
      this.lastPrice = undefined
      this.lastSize = undefined
    }


    if (this.quote.getTradedvolume()) {
      this.tradedVolume = toNumber(this.quote.getTradedvolume())
    }

  }
}