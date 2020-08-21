import { Menu, Button } from "@blueprintjs/core";
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
import TableView, { getConfiguredColumns, getColIdsInOrder, reorderColumnData, TableViewConfig, TableViewProperties } from "./TableView/TableView";
import { TabNode, Actions, Model } from "flexlayout-react";
import { ListingService } from "../services/ListingService";

interface MarketDepthProps extends TableViewProperties {
  node: TabNode,
  model: Model,
  quoteService: QuoteService,
  listingContext: ListingContext
  listingService: ListingService
}

interface MarketDepthState {
  listing?: Listing,
  quote?: ClobQuote,
  columns: Array<JSX.Element>
  columnWidths: Array<number>
  locked: boolean
  lockedListingId: number
}

interface MarketDepthConfig extends TableViewConfig {
  lockedListingId: number
  locked: boolean
}

export default class MarketDepth extends TableView<MarketDepthProps, MarketDepthState> implements QuoteListener {

  stream?: grpcWeb.ClientReadableStream<ClobQuote>;

  quoteService: QuoteService;
  listingService: ListingService;

  staticDataService = new StaticDataServiceClient(Login.grpcContext.serviceUrl, null, null)

  constructor(props: MarketDepthProps) {
    super(props);

    this.toggleLock = this.toggleLock.bind(this)

    this.quoteService = props.quoteService
    this.listingService = props.listingService

    let columns = this.getColumns()

    let config = this.props.node.getConfig() as MarketDepthConfig

    let [defaultCols, defaultColWidths] = getConfiguredColumns(columns, config);

    this.props.node.setEventListener("save", (p) => {
      let cols = this.state.columns
      let colOrderIds = getColIdsInOrder(cols);

      let lockedListingId = -1
      if (this.state.listing) {
        lockedListingId = this.state.listing.getId()
      }

      let persistentConfig: MarketDepthConfig = {
        columnWidths: this.state.columnWidths,
        columnOrder: colOrderIds,
        locked: this.state.locked,
        lockedListingId: lockedListingId
      }


      this.props.model.doAction(Actions.updateNodeAttributes(props.node.getId(), { config: persistentConfig }))
    });


    this.state = {
      locked: config.locked,
      lockedListingId: config.lockedListingId,
      columns: defaultCols,
      columnWidths: defaultColWidths
    };
  }

  public componentDidMount(): void {
    this.props.listingContext.addListener((listing: Listing) => {

      if (!this.state.locked) {

          this.setListing(listing)
      }
    })

    if (this.state.locked) {
        this.listingService.GetListing(this.state.lockedListingId, (response:Listing)=> {
          this.setListing(response)
      }) 
      
    }
  }

  private setListing(listing: Listing): void {
    if (this.state && this.state.listing) {
      if (this.state.listing === listing) {
        return
      }

      this.quoteService.UnsubscribeFromQuote(this.state.listing.getId(), this)
    }


    let quote = this.quoteService.SubscribeToQuote(listing, this)

    let state: MarketDepthState = {
      ...this.state, ...{
        listing: listing,
        quote: quote
      }
    }

    this.setState(state)
  }

  protected getColumns(): JSX.Element[] {
    return [
      <Column key="bidMic" id="bidMic" name="Bid Mic" cellRenderer={this.renderBidMic} />,
      <Column key="bidSize" id="bidSize" name="Bid Qty" cellRenderer={this.renderBidSize} />,
      <Column key="bidPx" id="bidPx" name="Bid Px" cellRenderer={this.renderBidPrice} />,
      <Column key="askPx" id="askPx" name="Ask Px" cellRenderer={this.renderAskPrice} />,
      <Column key="askSize" id="askSize" name="Ask Qty" cellRenderer={this.renderAskSize} />,
      <Column key="askMic" id="askMic" name="Ask Mic" cellRenderer={this.renderAskMic} />]
  }


  protected getTableName(): string {
    return "Market Depth"
  }

  onQuote(receivedQuote: ClobQuote): void {
    let state: MarketDepthState = {
      ...this.state, ...{
        quote: receivedQuote,
      }
    }

    // A bug in the table implementation means state has to be set twice to update the table
    this.setState(state);
    this.setState(state);
  }

  toggleLock() : void {
    let state: MarketDepthState = {
      ...this.state, ...{
        locked: !this.state.locked,
      }
    }
    this.setState(state);
  }

  public render() {
    return (
      <div className="bp3-dark">
        <Button icon={this.state.locked?"lock":"unlock"} onClick={this.toggleLock}>{this.getListingLabel()}</Button>
        <Table enableRowResizing={false} numRows={10} className="bp3-dark" enableColumnReordering={true}
          onColumnsReordered={this.onColumnsReordered} enableColumnResizing={true} onColumnWidthChanged={this.columnResized}
          columnWidths={this.state.columnWidths} bodyContextMenuRenderer={this.renderContextMenu} >
          {this.state.columns}
        </Table>
      </div>
    );
  }

  renderContextMenu = () => {
    return (

      <Menu >

        <Menu.Item text="Edit Visible Columns" onClick={() => this.editVisibleColumns()}  >
        </Menu.Item>
      </Menu>

    );
  };


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

  private getListingLabel(): string {
    if (this.state && this.state.listing) {
      return getListingShortName(this.state.listing)
    }

    return "(No Selection) "
  }

  private renderBidMic = (row: number) => {

    if (!this.state || !this.state.quote) {
      return (<Cell></Cell>)
    }

    let depth = this.state.quote.getBidsList()

    if (row < depth.length) {
      let listing = this.listingService.GetListingImmediate(depth[row].getListingid())
      if (listing) {
        return (<Cell>{listing.getMarket()?.getMic()}</Cell>)
      } else {
        return (<Cell></Cell>)
      }
    } else {
      return (<Cell></Cell>)
    }
  }

  private renderAskMic = (row: number) => {

    if (!this.state || !this.state.quote) {
      return (<Cell></Cell>)
    }

    let depth = this.state.quote.getOffersList()

    if (row < depth.length) {
      let listing = this.listingService.GetListingImmediate(depth[row].getListingid())
      if (listing) {
        return (<Cell>{listing.getMarket()?.getMic()}</Cell>)
      } else {
        return (<Cell></Cell>)
      }
    } else {
      return (<Cell></Cell>)
    }
  }


  private renderBidSize = (row: number) => {

    if (!this.state || !this.state.quote) {
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
    if (!this.state || !this.state.quote) {
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
    if (!this.state || !this.state.quote) {
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
    if (!this.state || !this.state.quote) {
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