import Login from "../components/Login";


import { MarketDataServiceClient } from "../serverapi/Market-data-serviceServiceClientPb";
import { SubscribeRequest, Quote, Subscription, AddSubscriptionResponse } from "../serverapi/market-data-service_pb";
import { ClientReadableStream, Status, Error } from "grpc-web";
import { logError, logDebug, logGrpcError } from "../logging/Logging";

export interface QuoteService {
  SubscribeToQuote(listingId: number, listener: QuoteListener): Quote | undefined
  UnsubscribeFromQuote(listingId: number, listener: QuoteListener ) : void
}

export interface QuoteListener {
  onQuote(quote: Quote): void
}
/**
 * Use this to subscribe to quotes to avoid multiple server side subscriptions to the same quote
 */
export default class QuoteServiceImpl implements QuoteService {

  marketDataService = new MarketDataServiceClient(Login.grpcContext.serviceUrl, null, null)

  stream?: ClientReadableStream<Quote>;

  idToListeners: Map<number, Array<QuoteListener>> = new Map()
  listingIdToQuote: Map<number, Quote> = new Map()

  constructor() {

    var subscription = new SubscribeRequest()
    subscription.setSubscriberid(Login.grpcContext.appInstanceId)

    this.stream = this.marketDataService.subscribe(subscription, Login.grpcContext.grpcMetaData)

    this.stream.on('data', (quote: Quote) => {

      this.listingIdToQuote.set(quote.getListingid(), quote)
      let listeners = this.idToListeners.get(quote.getListingid())
      if(  listeners ) {
        listeners.forEach(l => {
          l.onQuote(quote)
        })

      } else {
        logError("recieved update for quote for which there are no listeners, quote id:" + quote.getListingid())
      }


    });
    this.stream.on('status', (status: Status) => {
      if (status.metadata) {
        logDebug("market data service subscribe call metadata:" + status.metadata);
      }
    });
    this.stream.on('error', (err: Error) => {
      logGrpcError("market data subscription failed", err)
    });
    this.stream.on('end', () => {
      logDebug('stream end signal received')
    });

  }

  SubscribeToQuote(listingId: number, listener: QuoteListener): Quote | undefined {
    let listeners = this.idToListeners.get(listingId)
    if (!listeners) {
      listeners = new Array<QuoteListener>();
      this.idToListeners.set(listingId, listeners)

      let subscription = new Subscription()
      subscription.setListingid(listingId)
      subscription.setSubscriberid(Login.grpcContext.appInstanceId)
      this.marketDataService.addSubscription(subscription, Login.grpcContext.grpcMetaData, (err : Error, response : AddSubscriptionResponse)=> {
        if( err ) {
          logError("market data subscription failed:" + err)
        }
      })

    }

    listeners.push(listener)

    return this.listingIdToQuote.get(listingId)
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