import { Book } from "../serverapi/cmds_pb";
import { ClientMarketDataServiceClient } from "../serverapi/CmdsServiceClientPb";
import Login from "../components/Login";
import { Listing } from "../serverapi/listing_pb";




export interface QuoteService {

    SubscribeToQuote(listing : Listing, listener : QuoteListener ) : void

}

export interface QuoteListener {
    onQuote( quote : Book) : void
}


export default class QuoteServiceImpl implements QuoteService{

    marketDataService = new ClientMarketDataServiceClient(Login.grpcContext.serviceUrl, null, null)

    idToListeners: Map<number, QuoteListener[]> = new Map()

    constructor(clientId : string) {
        this.marketDataService.subscribe
    }


    SubscribeToQuote(listing : Listing, listener : QuoteListener ) : void {
        let listeners  = this.idToListeners.get(listing.getId())
        if( !listeners) {
            listeners =  new Array<QuoteListener>();
            this.idToListeners.set(listing.getId(), listeners )
            this.marketDataService.

        }

        listeners.push(listener)
    }

}