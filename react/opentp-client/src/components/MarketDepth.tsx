import {  Label } from "@blueprintjs/core";
import { Cell, Column, Table } from "@blueprintjs/table";
import * as grpcWeb from 'grpc-web';
import React from 'react';
import v4 from 'uuid';
import { Listing } from "../serverapi/listing_pb";
import {  MdsSubscribeRequest } from '../serverapi/market-data-service_pb';
import { StaticDataServiceClient } from "../serverapi/Static-data-serviceServiceClientPb";
import { QuoteListener, QuoteService } from "../services/QuoteService";
import { toNumber } from "../util/decimal64Conversion";
import { ListingContext } from "./Container";
import Login from "./Login";
import './OrderBlotter.css';
import { getListingShortName } from "../common/modelutilities";
import { ClobQuote } from "../serverapi/clobquote_pb";

interface MarketDepthProps {
  quoteService : QuoteService,
  listingContext : ListingContext
}

interface MarketDepthState {
  listing?: Listing,
  quote?: ClobQuote,
}

export default class MarketDepth extends React.Component<MarketDepthProps, MarketDepthState> implements QuoteListener {

  stream?: grpcWeb.ClientReadableStream<ClobQuote>;

  id: string;

  quoteService: QuoteService;

  count: number;

  staticDataService = new StaticDataServiceClient(Login.grpcContext.serviceUrl, null, null)

  constructor(props: MarketDepthProps) {
    super(props);

    this.quoteService = props.quoteService

    this.state = {};

    this.id = v4();

    this.count = 0;

    var subscription = new MdsSubscribeRequest()
    subscription.setSubscriberid(this.id)

    this.props.listingContext.addListener((listing:Listing)=> {

      if( this.state && this.state.listing ){
        if( this.state.listing === listing) {
          return
        }

        this.quoteService.UnsubscribeFromQuote(this.state.listing.getId(), this)  
      }

      let state: MarketDepthState = {
        ...this.state,...{
          listing: listing,
          quote: undefined
        }
      }

      this.setState(state)
      this.setState(state)
      this.quoteService.SubscribeToQuote(listing.getId(), this)
    })

  }

  onQuote(receivedQuote: ClobQuote): void {
    let state: MarketDepthState = {
      ...this.state,...{
        quote: receivedQuote,
      }
    }

    // A bug in the table implementation means state has to be set twice to update the table
    this.setState(state);
    this.setState(state);
  }


  public render() {
      return (
        <div className="bp3-dark">
          <Label>{this.getListingLabel()}</Label>
          <Table enableRowResizing={false} numRows={10} className="bp3-dark">
            <Column name="Bid Size" cellRenderer={this.renderBidSize} />
            <Column name="Bid Px" cellRenderer={this.renderBidPrice} />
            <Column name="Ask Px" cellRenderer={this.renderAskPrice} />
            <Column name="Ask Size" cellRenderer={this.renderAskSize} />
          </Table>
        </div>);
  }

  private getListingLabel(): string  {
    if( this.state && this.state.listing) {
      return getListingShortName(this.state.listing)
    }

    return "(No Selection) "
  }

  private renderBidSize = (row: number) => {

    if( !this.state || !this.state.quote) {
      return (<Cell></Cell>)
    }



    let depth = this.state.quote.getBidsList()

    if (row < depth.length) {
      return (<Cell>{toNumber(depth[row].getSize())}</Cell>)
    } else {
      return (<Cell></Cell>)
    }
  }

  private renderAskSize = (row: number) => {
    if( !this.state || !this.state.quote) {
      return (<Cell></Cell>)
    }

    let depth = this.state.quote.getOffersList()

    if (row < depth.length) {
      return (<Cell>{toNumber(depth[row].getSize())}</Cell>)
    } else {
      return (<Cell></Cell>)
    }
  }

  private renderBidPrice = (row: number) => {
    if( !this.state || !this.state.quote) {
      return (<Cell></Cell>)
    }

    let depth = this.state.quote.getBidsList()

    if (row < depth.length) {
      let line = depth[row]
      return (<Cell>{toNumber(line.getPrice())}</Cell>)
    } else {
      return (<Cell></Cell>)
    }
  }

  private renderAskPrice = (row: number) => {
    if( !this.state || !this.state.quote) {
      return (<Cell></Cell>)
    }

    let depth = this.state.quote.getOffersList()

    if (row < depth.length) {
      return (<Cell>{toNumber(depth[row].getPrice())}</Cell>)
    } else {
      return (<Cell></Cell>)
    }
  }

}