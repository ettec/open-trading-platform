import {  Error } from "grpc-web";
import Login from "../components/Login";
import { logError } from "../logging/Logging";
import { Listing } from "../serverapi/listing_pb";
import { StaticDataServiceClient } from "../serverapi/Static-data-serviceServiceClientPb";
import { ListingId } from "../serverapi/static-data-service_pb";



export interface ListingService {

  GetListingImmediate(listingId: number) : Listing | null
  
  GetListing(listingId: number, listener: (
    response: Listing) => void): void
}


/**
 * Use this to subscribe to quotes to avoid multiple server side subscriptions to the same listing
 */
export default class ListingServiceImpl implements ListingService {

  staticDataService = new StaticDataServiceClient(Login.grpcContext.serviceUrl, null, null)

  idToListeners: Map<number, Array<(response: Listing) => void>> = new Map()
  listingIdToListing: Map<number, Listing> = new Map()

  GetListingImmediate(listingId: number) : Listing | null {
    if (listingId <= 0) {
      return null
    }

    let listing = this.listingIdToListing.get(listingId)
    if (listing) {
      return listing
    }

    let listeners = this.idToListeners.get(listingId)
    if (!listeners) {
      this.GetListing(listingId, (response : Listing) => void {
        // Do nothing
      })
    }

    return null
  }

  

  GetListing(listingId: number, listener: (
    response: Listing) => void) {

    if (listingId <= 0) {
      return
    }

    let listing = this.listingIdToListing.get(listingId)
    if (listing) {
      listener(listing)
      return;
    }

    let listeners = this.idToListeners.get(listingId)
    if (!listeners) {
      listeners = new Array<(response: Listing) => void>();
      this.idToListeners.set(listingId, listeners)

      let listingParam = new ListingId()
      listingParam.setListingid(listingId)
      this.staticDataService.getListing(listingParam, Login.grpcContext.grpcMetaData, (err: Error, listing: Listing) => {
        if (err) {
          logError("get listing for id " + listingId + " failed:" + err)
        } else {
          this.listingIdToListing.set(listing.getId(), listing)
          let ls = this.idToListeners.get(listing.getId())
          if (ls) {
            ls.forEach(l => {
              l(listing)
            })
          }
        }
      })
    }

    listeners.push(listener)
  }




}