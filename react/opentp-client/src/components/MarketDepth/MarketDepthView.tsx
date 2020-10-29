import { toNumber } from "../../common/decimal64Conversion"
import { ClobQuote } from "../../serverapi/clobquote_pb"
import { Side } from "../../serverapi/order_pb"
import { ListingService } from "../../services/ListingService"

export class MarketDepthView {

    lines: Array<DepthLine>
    listingService: ListingService
    updateListener: () => void
  
    constructor(numLines: number, listingService: ListingService, updateListener: () => void) {
      this.lines = new Array<DepthLine>()
      this.listingService = listingService
      this.updateListener = updateListener
      for (let i = 0; i < numLines; i++) {
        this.lines.push(new DepthLine(i))
      }
  
    }
  
    getDepthAtIdx(idx: number, side: Side): { price: number | undefined, quantity: number | undefined } {
  
      let price: number | undefined;
      let quantity: number | undefined;
  
  
      for (let i = 0; i <= idx; i++) {
        let linePrice: number | undefined;
        let lineSize: number | undefined;
  
        if (side === Side.BUY) {
          linePrice = this.lines[i].bidPrice
          lineSize = this.lines[i].bidSize
        } else {
          linePrice = this.lines[i].askPrice
          lineSize = this.lines[i].askSize
        }
  
        if (lineSize && linePrice) {
          price = linePrice
          if (quantity) {
            quantity += lineSize
          } else {
            quantity = lineSize
          }
  
        } else {
          break
        }
  
      }
  
  
      return { price, quantity }
  
    }
  
    setQuote(quote: ClobQuote) {
      for (let i = 0; i < this.lines.length; i++) {
        let line = this.lines[i]
        line.bidMic = undefined
        line.bidSize = undefined
        line.bidPrice = undefined
        line.askPrice = undefined
        line.askSize = undefined
        line.askMic = undefined
  
        let depthList = quote.getOffersList()
  
        if (i < depthList.length) {
          let depth = depthList[i]
          let listing = this.listingService.GetListingImmediate(depth.getListingid())
          if (listing) {
            line.askMic = listing.getMarket()?.getMic()
          }
          line.askPrice = toNumber(depth.getPrice())
          line.askSize = toNumber(depth.getSize())
        }
  
        depthList = quote.getBidsList()
  
        if (i < depthList.length) {
          let depth = depthList[i]
          let listing = this.listingService.GetListingImmediate(depth.getListingid())
          if (listing) {
            line.bidMic = listing.getMarket()?.getMic()
          }
          line.bidPrice = toNumber(depth.getPrice())
          line.bidSize = toNumber(depth.getSize())
        }
  
      }
  
      this.updateListener()
    }
  
  }
  
 export class DepthLine {
  
    idx: number
    bidMic?: string
    bidSize?: number
    bidPrice?: number
    askPrice?: number
    askSize?: number
    askMic?: string
  
    constructor(idx: number) {
      this.idx = idx
    }
  }
  