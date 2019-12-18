import { ClientReadableStream, Error } from "grpc-web";
import Login from "../components/Login";
import { logError } from "../logging/Logging";
import { Listing } from "../serverapi/listing_pb";
import { Quote } from "../serverapi/market-data-service_pb";
import { StaticDataServiceClient } from "../serverapi/Static-data-serviceServiceClientPb";
import { ListingId } from "../serverapi/static-data-service_pb";



export interface ListingService {
  GetListing(listingId: number, listener: ListingListener): void
}

export interface ListingListener {
  onListing(listing: Listing ): void
}
/**
 * Use this to subscribe to quotes to avoid multiple server side subscriptions to the same quote
 */
export default class ListingServiceImpl implements ListingService {

  staticDataService = new StaticDataServiceClient(Login.grpcContext.serviceUrl, null, null)

  stream?: ClientReadableStream<Quote>;

  idToListeners: Map<number, Array<ListingListener>> = new Map()
  listingIdToListing: Map<number, Listing> = new Map()

  constructor() {
  }

  GetListing(listingId: number, listener: ListingListener) {
    let listing = this.listingIdToListing.get(listingId)
    if( listing ) {
        listener.onListing(listing)
        return;
    }

    let listeners = this.idToListeners.get(listingId)
    if (!listeners) {
      listeners = new Array<ListingListener>();
      this.idToListeners.set(listingId, listeners)

      let listingParam = new ListingId()
      listingParam.setListingid(listingId)
      this.staticDataService.getListing(listingParam, Login.grpcContext.grpcMetaData, (err: Error, listing: Listing)=> {
            if(err) {
                logError("get listing for id " + listingId + " failed:" + err)
            } else {
                this.listingIdToListing.set(listing.getId(), listing)
                let ls = this.idToListeners.get(listing.getId())
                if( ls ) {
                    ls.forEach(l=> {
                        l.onListing(listing)
                    })
                }
            }
      })
    }
  }

  


}