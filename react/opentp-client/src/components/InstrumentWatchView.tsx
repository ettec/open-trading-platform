import { TabNode, Model, Actions } from "flexlayout-react";
import React from 'react';

import v4 from 'uuid';

import InstrumentSearchBar from "./InstrumentSearchBar";
import './OrderBlotter.css';
import { Listing } from "../serverapi/listing_pb";
import Login from "./Login";
import { StaticDataServiceClient } from "../serverapi/Static-data-serviceServiceClientPb";
import { ListingIds, Listings } from "../serverapi/static-data-service_pb";
import { Error } from "grpc-web";
import { QuoteService, QuoteListener } from "../services/QuoteService";
import { Quote } from "../serverapi/market-data-service_pb";
import { logDebug } from "../logging/Logging";





interface InstrumentWatchState {
  watches: ListingWatchLine[]
}

interface InstrumentWatchProps {
  node: TabNode,
  model: Model,
  quoteService: QuoteService
}

interface PersistentConfig {
  listingIds: number[]
}



export default class InstrumentWatchView extends React.Component<InstrumentWatchProps, InstrumentWatchState> implements QuoteListener {

  staticDataService = new StaticDataServiceClient(Login.grpcContext.serviceUrl, null, null)
  quoteService: QuoteService

  watchMap: Map<number, ListingWatchLine> = new Map()

  constructor(props: InstrumentWatchProps) {
    super(props);

    this.quoteService = props.quoteService

    let initialState: InstrumentWatchState = {
      watches: Array.from(this.watchMap.values())
    }

    this.state = initialState;

    this.addListing = this.addListing.bind(this);

    this.props.node.setEventListener("save", (p) => {
      let persistentConfig: PersistentConfig = { listingIds: Array.from(this.watchMap.keys()) }
      this.props.model.doAction(Actions.updateNodeAttributes(props.node.getId(), { config: persistentConfig }))
    });

    if (this.props.node.getConfig() && this.props.node.getConfig()) {
      let persistentConfig: PersistentConfig = this.props.node.getConfig();
      this.addListingLineByIds(persistentConfig.listingIds)
    }


  }

  addListing(listing?: Listing) {

    if (listing) {

      if (this.watchMap.has(listing.getId())) {
        return;
      }

      this.addListingLine(listing);
    }

  }

  private addListingLineByIds(listingIds: number[]) {
    let ids = new ListingIds()
    ids.setListingidsList(listingIds)

    this.staticDataService.getListings(ids, Login.grpcContext.grpcMetaData,
      (err: Error, listings: Listings) => {

        listings.getListingsList().forEach((listing) => {
          this.addListingLine(listing)
        })
      })
  }

  private addListingLine(listing: Listing) {

    if (!this.watchMap.get(listing.getId())) {
      let line = new ListingWatchLine()
      line.listing = listing

      let listingWatchLine = new ListingWatchLine()
      listingWatchLine.listing = listing


      this.watchMap.set(listing.getId(), listingWatchLine);
      this.setState({
        watches: Array.from(this.watchMap.values())
      });

      this.quoteService.SubscribeToQuote(listing.getId(), this)


    }


  }

  onQuote(quote: Quote): void {

    let line = this.watchMap.get(quote.getListingid())
    if( line ) {

      line.quote = quote
      this.setState({
        watches: Array.from(this.watchMap.values())
      });

    } else {
      logDebug("received quote update for non-existent line, quote:" + quote)
    }

  }

  public render() {
    var watches: ListingWatchLine[];
    if (this.state) {
      watches = Object.assign([], this.state.watches);
    } else {
      watches = []
    }

    const clonedWatches = watches;


    return (

      <div className="bp3-dark">

        <InstrumentSearchBar add={this.addListing} />



      </div>


    );
  }

}

class ListingWatchLine {
  listing?: Listing;
  quote?: Quote;
}