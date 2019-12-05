import { Button, InputGroup } from "@blueprintjs/core";
import React from 'react';


import v4 from 'uuid';
import './OrderBlotter.css';
import * as grpcWeb from 'grpc-web'
import { ClientMarketDataServiceClient } from '../serverapi/CmdsServiceClientPb'
import { SubscribeRequest, Book, BookLine, Subscription } from '../serverapi/cmds_pb'
import { LocalBookLine } from '../model/Model'
import { GrpcConsumer } from "./GrpcContextProvider/GrpcContextProvider";
import Login from "./Login";
import { Table, Column, Cell } from "@blueprintjs/table";



interface MarketDepthState {
  symbol?: string,
  book: Book
}

export default class MarketDepth extends React.Component<{}, MarketDepthState> {

  stream?: grpcWeb.ClientReadableStream<Book>;

  id: string;

  marketDataService = new ClientMarketDataServiceClient(Login.grpcContext.serviceUrl, null, null)

  count: number;



  constructor() {
    super({});
    this.setState({ symbol: '' });

    this.id = v4();

    this.count = 0;

    var subscription = new SubscribeRequest()
    subscription.setSubscriberid(this.id)

    this.stream = this.marketDataService.subscribe(subscription, Login.grpcContext.grpcMetaData)

    this.stream.on('data', (response: Book) => {

      let blotterState: MarketDepthState = {
        ...this.state, ... {
          book: response,
        }
      }

      // A bug in the table implementation means state has to be set twice to update the table
      this.setState(blotterState);
      this.setState(blotterState);
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
    });

    this.handleSymbolChange = this.handleSymbolChange.bind(this);
    this.onSubscribe = this.onSubscribe.bind(this);

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
    }

  }




  public render() {




    if (this.state && this.state.book) {
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
    let depth = this.state.book.getDepthList()

    if (row < depth.length) {
      return (<Cell>{depth[row].getBidsize()}</Cell>)
    } else {
      return (<Cell></Cell>)
    }
  }

  private renderAskSize = (row: number) => {
    let depth = this.state.book.getDepthList()

    if (row < depth.length) {
      return (<Cell>{depth[row].getAsksize()}</Cell>)
    } else {
      return (<Cell></Cell>)
    }
  }

  private renderBidPrice = (row: number) => {
    let depth = this.state.book.getDepthList()

    if (row < depth.length) {
      let line = depth[row]
      return (<Cell>{line.getBidprice()}</Cell>)
    } else {
      return (<Cell></Cell>)
    }
  }

  private renderAskPrice = (row: number) => {
    let depth = this.state.book.getDepthList()

    if (row < depth.length) {
      return (<Cell>{depth[row].getAskprice()}</Cell>)
    } else {
      return (<Cell></Cell>)
    }
  }



}