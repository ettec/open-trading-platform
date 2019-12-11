import { Button, InputGroup } from "@blueprintjs/core";
import React from 'react';


import v4 from 'uuid';
import './OrderBlotter.css';
import * as grpcWeb from 'grpc-web'

import { SubscribeRequest, Quote, DepthLine, Subscription } from '../serverapi/market-data-service_pb'
import Login from "./Login";
import { Table, Column, Cell } from "@blueprintjs/table";
import  { QuoteService, QuoteListener } from "../services/QuoteService";
import { StaticDataServiceClient } from "../serverapi/Static-data-serviceServiceClientPb";
import { toNumber } from "../util/decimal64Conversion";



interface MarketDepthProps {
  quoteService : QuoteService
}

interface MarketDepthState {
  symbol?: string,
  quote: Quote
}

export default class MarketDepth extends React.Component<MarketDepthProps, MarketDepthState> implements QuoteListener {

  stream?: grpcWeb.ClientReadableStream<Quote>;

  id: string;

  quoteService: QuoteService;

  count: number;

  staticDataService = new StaticDataServiceClient(Login.grpcContext.serviceUrl, null, null)



  constructor(props: MarketDepthProps) {
    super(props);

    this.quoteService = props.quoteService

    this.setState({ symbol: '' });

    this.id = v4();

    this.count = 0;

    var subscription = new SubscribeRequest()
    subscription.setSubscriberid(this.id)


    this.quoteService.SubscribeToQuote(54123, this)

 

    
/*
    this.stream = this.marketDataService.subscribe(subscription, Login.grpcContext.grpcMetaData)

    this.stream.on('data', (response: Quote) => {

      let state: MarketDepthState = {
        ...this.state, ... {
          quote: response,
        }
      }

      // A bug in the table implementation means state has to be set twice to update the table
      this.setState(state);
      this.setState(state);
    });
    this.stream.on('status', (status: grpcWeb.Status) => {
      if (status.metadata) {
        console.log(status.metadata);
      }
    });
    this.stream.on('error', (err: grpcWeb.Error) => {
      console.log('Received error')
    });
    this.stream.on('end', () => {
      console.log('stream end signal received');
    });*/

    this.handleSymbolChange = this.handleSymbolChange.bind(this);
    this.onSubscribe = this.onSubscribe.bind(this);

  }

  onQuote(receivedQuote: Quote): void {
    let state: MarketDepthState = {
      ...this.state, ... {
        quote: receivedQuote,
      }
    }

    // A bug in the table implementation means state has to be set twice to update the table
    this.setState(state);
    this.setState(state);
  }

  handleSymbolChange(e: any) {

    if (e.target && e.target.value) {

      let newSymbol: string = e.target.value;

      let blotterState: MarketDepthState = {
        ...this.state, ... {
          symbol: newSymbol ,
        }
      }

      this.setState(state => (blotterState))
    }
  }

  onSubscribe() {
/*
    console.log("Subscribe to:" + this.state.symbol)

    if (this.state.symbol != null) {
      var subscription = new Subscription()
      subscription.setSubscriberid(this.id)
      subscription.setListingid(this.state.symbol)
      console.log("adding subscription:" + subscription)
      this.marketDataService.addSubscription(subscription, Login.grpcContext.grpcMetaData, (err, response) => {
        if (err) {
          console.log("failed to add subscription:" + err)
          return
        }

        if (response) {
          console.log("Add subscription response:" + response)
        }


      })
    } */

  }




  public render() {

    if (this.state && this.state.quote) {
      return (
        <div className="bp3-dark">
          <InputGroup onChange={this.handleSymbolChange} />
          <Button onClick={this.onSubscribe} >Subscribe</Button>

          <Table enableRowResizing={false} numRows={10} className="bp3-dark">
            <Column name="Bid Size" cellRenderer={this.renderBidSize} />
            <Column name="Bid Px" cellRenderer={this.renderBidPrice} />
            <Column name="Ask Px" cellRenderer={this.renderAskPrice} />
            <Column name="Ask Size" cellRenderer={this.renderAskSize} />
          </Table>
        </div>);

    } else {
      return (
        <div className="bp3-dark">
          <InputGroup onChange={this.handleSymbolChange} />
          <Button onClick={this.onSubscribe} >Subscribe</Button>
        </div>
      );
    }

  }

  private renderBidSize = (row: number) => {
    let depth = this.state.quote.getDepthList()

    if (row < depth.length) {
      return (<Cell>{toNumber(depth[row].getBidsize())}</Cell>)
    } else {
      return (<Cell></Cell>)
    }
  }

  private renderAskSize = (row: number) => {
    let depth = this.state.quote.getDepthList()

    if (row < depth.length) {
      return (<Cell>{toNumber(depth[row].getAsksize())}</Cell>)
    } else {
      return (<Cell></Cell>)
    }
  }

  private renderBidPrice = (row: number) => {
    let depth = this.state.quote.getDepthList()

    if (row < depth.length) {
      let line = depth[row]
      return (<Cell>{toNumber(line.getBidprice())}</Cell>)
    } else {
      return (<Cell></Cell>)
    }
  }

  private renderAskPrice = (row: number) => {
    let depth = this.state.quote.getDepthList()

    if (row < depth.length) {
      return (<Cell>{toNumber(depth[row].getAskprice())}</Cell>)
    } else {
      return (<Cell></Cell>)
    }
  }

}