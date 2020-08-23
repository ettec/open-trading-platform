import { Error, ClientReadableStream } from 'grpc-web';
import Login from "../components/Login";
import { logError } from "../logging/Logging";
import { ClobQuote } from "../serverapi/clobquote_pb";
import { Listing } from "../serverapi/listing_pb";
import { MarketDataServiceClient } from "../serverapi/Market-data-serviceServiceClientPb";
import { MdsConnectRequest, MdsSubscribeRequest } from "../serverapi/market-data-service_pb";
import { Empty } from "../serverapi/modelcommon_pb";
import { ListingService } from "./ListingService";
import Stream from "./impl/Stream";



export interface QuoteService {
  SubscribeToQuote(listing: Listing, listener: QuoteListener): ClobQuote | undefined
  UnsubscribeFromQuote(listingId: number, listener: QuoteListener): void
}

export interface QuoteListener {
  onQuote(quote: ClobQuote): void
}



/**
 * Use this to subscribe to quotes to avoid multiple server side subscriptions to the same quote
 */
export default class QuoteServiceImpl implements QuoteService {

  marketDataService = new MarketDataServiceClient(Login.grpcContext.serviceUrl, null, null)

  stream?: Stream<ClobQuote>;

  listingService: ListingService
  idToListeners: Map<number, Array<QuoteListener>> = new Map()
  listingIdToQuote: Map<number, ClobQuote> = new Map()
  nextResubscribeInterval = 1000


  constructor(listingService: ListingService) {
    this.listingService = listingService

    this.stream = new Stream(() : ClientReadableStream<any>  =>  {
      var subscription = new MdsConnectRequest()
      subscription.setSubscriberid(Login.grpcContext.appInstanceId)

      let result = this.marketDataService.connect(subscription, Login.grpcContext.grpcMetaData)

      let keys = this.listingIdToQuote.keys()
      for (var listingId of keys) {
        this.listingService.GetListing(listingId, (listing: Listing)=> {
          let subscription = new MdsSubscribeRequest()
          subscription.setListing(listing)
          subscription.setSubscriberid(Login.grpcContext.appInstanceId)
          this.marketDataService.subscribe(subscription, Login.grpcContext.grpcMetaData, (err: Error,
            response: Empty) => {
            if (err) {
              logError("market data subscription failed:" + err)
            }
          })
        })


      }

      return result
    }, (quote: ClobQuote) => {
      this.listingIdToQuote.set(quote.getListingid(), quote);
      let listeners = this.idToListeners.get(quote.getListingid());
      if (listeners) {
        listeners.forEach(l => {
          l.onQuote(quote);
        });
      }
    }, "quote stream")


  }


  SubscribeToQuote(listing: Listing, listener: QuoteListener): ClobQuote | undefined {
    let listeners = this.idToListeners.get(listing.getId())
    if (!listeners) {
      listeners = new Array<QuoteListener>();
      this.idToListeners.set(listing.getId(), listeners)

      let subscription = new MdsSubscribeRequest()
      subscription.setListing(listing)
      subscription.setSubscriberid(Login.grpcContext.appInstanceId)
      this.marketDataService.subscribe(subscription, Login.grpcContext.grpcMetaData, (err: Error,
        response: Empty) => {
        if (err) {
          logError("market data subscription failed:" + err)
        }
      })

    }

    listeners.push(listener)

    return this.listingIdToQuote.get(listing.getId())
  }

  UnsubscribeFromQuote(listingId: number, listener: QuoteListener) {
    let listeners = this.idToListeners.get(listingId)
    if (listeners) {
      const index = listeners.indexOf(listener);
      if (index > -1) {
        listeners.splice(index, 1);
      }
    }


  }


}