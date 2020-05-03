import Login from "../components/Login";

import { Error } from 'grpc-web';
import { MarketDataServiceClient } from "../serverapi/Market-data-serviceServiceClientPb";
import {  MdsSubscribeRequest, MdsConnectRequest } from "../serverapi/market-data-service_pb";
import { ClientReadableStream, Status } from "grpc-web";
import { logError, logDebug, logGrpcError } from "../logging/Logging";
import { ClobQuote } from "../serverapi/clobquote_pb";
import { Empty } from "../serverapi/modelcommon_pb";
import { Listing } from "../serverapi/listing_pb";
import { ListingService } from "./ListingService";

export interface QuoteService {
  SubscribeToQuote(listing: Listing, listener: QuoteListener): ClobQuote | undefined
  UnsubscribeFromQuote(listingId: number, listener: QuoteListener ) : void
}

export interface QuoteListener {
  onQuote(quote: ClobQuote): void
}
/**
 * Use this to subscribe to quotes to avoid multiple server side subscriptions to the same quote
 */
export default class QuoteServiceImpl implements QuoteService {

  marketDataService = new MarketDataServiceClient(Login.grpcContext.serviceUrl, null, null)

  stream?: ClientReadableStream<ClobQuote>;

  listingService: ListingService
  idToListeners: Map<number, Array<QuoteListener>> = new Map()
  listingIdToQuote: Map<number, ClobQuote> = new Map()
  nextResubscribeInterval = 1000


  constructor(listingService: ListingService ) {
    this.listingService = listingService
    this.subscribeToQuoteStream();

    this.subscribeToQuoteStream = this.subscribeToQuoteStream.bind(this);

  }

  private subscribeToQuoteStream() {

    console.log("subscribing to quote stream")
    var subscription = new MdsConnectRequest()
    subscription.setSubscriberid(Login.grpcContext.appInstanceId)

    this.stream = this.marketDataService.connect(subscription, Login.grpcContext.grpcMetaData);

    this.stream.on('data', (quote: ClobQuote) => {
      this.nextResubscribeInterval = 1000
      this.listingIdToQuote.set(quote.getListingid(), quote);
      let listeners = this.idToListeners.get(quote.getListingid());
      if (listeners) {
        listeners.forEach(l => {
          l.onQuote(quote);
        });
      }
    });
    this.stream.on('status', (status: Status) => {
      console.log("market data stream status:" + status.details)
      if (status.metadata) {
        logDebug("market data service subscribe call metadata:" + status.metadata);
      }
    });
    this.stream.on('error', (err: Error) => {

      
      let nextInterval = this.getNextResubscribeInterval()
      logGrpcError("market data subscription failed, resubscribing in " + this.nextResubscribeInterval + "ms", err);
      setTimeout(this.subscribeToQuoteStream, nextInterval)
    });
    this.stream.on('end', () => {
      let nextInterval = this.getNextResubscribeInterval()
      logDebug("market data stream end signal received, resubscribing in " + this.nextResubscribeInterval + "ms");
      setTimeout(this.subscribeToQuoteStream, nextInterval)
    });


    let keys = this.listingIdToQuote.keys()
    for (var listingId of keys) {
      let listing = this.listingService.GetListingImmediate(listingId)

      if (listing ) {
        let subscription = new MdsSubscribeRequest()
        subscription.setListing(listing)
        subscription.setSubscriberid(Login.grpcContext.appInstanceId)
        this.marketDataService.subscribe(subscription, Login.grpcContext.grpcMetaData, (err: Error,
          response: Empty)=> {
          if( err ) {
            logError("market data subscription failed:" + err)
          }
        })
      }

  
    }

  }


  private getNextResubscribeInterval(): number {
    this.nextResubscribeInterval = this.nextResubscribeInterval * 2
    if (this.nextResubscribeInterval > 30000) {
      this.nextResubscribeInterval = 30000
    }

    return this.nextResubscribeInterval
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
        response: Empty)=> {
        if( err ) {
          logError("market data subscription failed:" + err)
        }
      })

    }

    listeners.push(listener)

    return this.listingIdToQuote.get(listing.getId())
  }

  UnsubscribeFromQuote(listingId: number, listener: QuoteListener ) {
    let listeners = this.idToListeners.get(listingId)
    if (listeners) {
        const index = listeners.indexOf(listener);
        if (index > -1) {
          listeners.splice(index, 1);
        }
    }

    
  }


}