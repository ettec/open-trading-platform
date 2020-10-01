import { Listing, TickSizeEntry} from "../serverapi/listing_pb"
import { toNumber } from "./decimal64Conversion"


export function getListingShortName(listing:Listing ): string  {
    
      let i = listing.getInstrument() 
      let m = listing.getMarket() 
      if( i && m ){
        return i.getDisplaysymbol() + " - " + m.getMic()
      } else {
        return "Listing:" + listing.getId() + " missing instrument or market"
      }
    
  }

  export function getListingLongName(listing:Listing ): string  {
    
    let i = listing.getInstrument() 
    let m = listing.getMarket() 
    if( i && m ){
      return i.getName() + " - " + m.getMic()
    } else {
      return "Listing:" + listing.getId() + " missing instrument or market"
    }
  
}

export function roundToTick(price: number, listing: Listing): number {
  if( price > 0 ) {
    let tickSize = getTickSize(price, listing)
    let numTicks = Math.round(price/tickSize)
    return numTicks * tickSize
  } else {
    return price
  }
  
}

export function getTickSize(price: number, listing: Listing): number {
  let tt = listing.getTicksize()
  if (tt) {
    let el = tt.getEntriesList()
    for (var entry of el) {
      let tickSize = tickSizeFromEntry(entry, price)
      if (tickSize) {
        return tickSize
      }
    }
  }

  return 1
}

export function tickSizeFromEntry(entry: TickSizeEntry, price: number): number | undefined {
  let lowerBound = toNumber(entry.getLowerpricebound())
  let upperBound = toNumber(entry.getUpperpricebound())

  if (lowerBound !== undefined && upperBound !== undefined) {
    if (price >= lowerBound && price <= upperBound) {
      return toNumber(entry.getTicksize())
    }
  }

  return undefined
}

