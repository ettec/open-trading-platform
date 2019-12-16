import { Listing } from "../serverapi/listing_pb"


export function getListingLabel(listing:Listing ): string  {
    
      let i = listing.getInstrument() 
      let m = listing.getMarket() 
      if( i && m ){
        return i.getDisplaysymbol() + " - " + m.getMic()
      } else {
        return "Listing:" + listing.getId() + " missing instrument or market"
      }
    
  }