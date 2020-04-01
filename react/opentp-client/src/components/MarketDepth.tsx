import {  Label } from "@blueprintjs/core";
import { Cell, Column, Table } from "@blueprintjs/table";
import * as grpcWeb from 'grpc-web';
import React from 'react';
import { Listing } from "../serverapi/listing_pb";
import { StaticDataServiceClient } from "../serverapi/Static-data-serviceServiceClientPb";
import { QuoteListener, QuoteService } from "../services/QuoteService";
import { toNumber } from "../util/decimal64Conversion";
import { ListingContext } from "./Container";
import Login from "./Login";
import './TableView/TableCommon.css';
import { getListingShortName } from "../common/modelutilities";
import { ClobQuote } from "../serverapi/clobquote_pb";
import TableViewConfig, { getColumnState, getColIdsInOrder, reorderColumnData } from "./TableView/TableLayout";
import { TabNode, Actions, Model } from "flexlayout-react";

interface MarketDepthProps {
  node: TabNode,
  model: Model,
  quoteService : QuoteService,
  listingContext : ListingContext
}

interface MarketDepthState {
  listing?: Listing,
  quote?: ClobQuote,
  columns: Array<JSX.Element>
  columnWidths: Array<number>
}

export default class MarketDepth extends React.Component<MarketDepthProps, MarketDepthState> implements QuoteListener {

  stream?: grpcWeb.ClientReadableStream<ClobQuote>;

  quoteService: QuoteService;

  staticDataService = new StaticDataServiceClient(Login.grpcContext.serviceUrl, null, null)

  constructor(props: MarketDepthProps) {
    super(props);

    this.quoteService = props.quoteService

    let columns = [<Column key="bidSize" id="bidSize" name="Bid Qty" cellRenderer={this.renderBidSize} />,
    <Column  key="bidPx" id="bidPx" name="Bid Px" cellRenderer={this.renderBidPrice} />,
    <Column key="askPx" id="askPx" name="Ask Px" cellRenderer={this.renderAskPrice} />,
    <Column key="askSize" id="askSize" name="Ask Qty" cellRenderer={this.renderAskSize} />]  

    let config = this.props.node.getConfig()

    let { defaultCols, defaultColWidths } = getColumnState(columns, config);

    this.props.node.setEventListener("save", (p) => {
      let cols = this.state.columns
      let colOrderIds = getColIdsInOrder(cols);

      let persistentConfig: TableViewConfig = {
        columnWidths: this.state.columnWidths,
        columnOrder: colOrderIds,
      }


      this.props.model.doAction(Actions.updateNodeAttributes(props.node.getId(), { config: persistentConfig }))
    });


    this.state = {
      columns: defaultCols,
      columnWidths: defaultColWidths
    };  

    this.props.listingContext.addListener((listing:Listing)=> {

      if( this.state && this.state.listing ){
        if( this.state.listing === listing) {
          return
        }

        this.quoteService.UnsubscribeFromQuote(this.state.listing.getId(), this)  
      }

   
      let quote = this.quoteService.SubscribeToQuote(listing, this)

      let state: MarketDepthState = {
        ...this.state,...{
          listing: listing,
          quote: quote
        }
      }

       // A bug in the table implementation means state has to be set twice to update the table
       this.setState(state)
       this.setState(state)

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
          <Table enableRowResizing={false} numRows={10} className="bp3-dark"enableColumnReordering={true}
          onColumnsReordered={this.onColumnsReordered} enableColumnResizing={true} onColumnWidthChanged={this.columnResized} 
          columnWidths={this.state.columnWidths}>
          {this.state.columns}
        </Table>
      </div>
    );
  }

  columnResized = (index: number, size: number) => {
    let newColWidths = this.state.columnWidths.slice();
    newColWidths[index] = size
    let blotterState: MarketDepthState = {
      ...this.state, ...{
        columnWidths: newColWidths
      }
    }

    this.setState(blotterState)

  }

  onColumnsReordered = (oldIndex: number, newIndex: number, length: number) => {

    let newCols = reorderColumnData(oldIndex, newIndex, length, this.state.columns)
    let newColWidths = reorderColumnData(oldIndex, newIndex, length, this.state.columnWidths)

    let blotterState = {
      ...this.state, ...{
        columns: newCols,
        columnWidths: newColWidths
      }
    }

    this.setState(blotterState)
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