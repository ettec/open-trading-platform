import { Cell, Column, IRegion, SelectionModes, Table} from "@blueprintjs/table";
import { Actions, Model, TabNode } from "flexlayout-react";
import { Error } from "grpc-web";
import React from 'react';
import { logDebug } from "../logging/Logging";
import { Listing } from "../serverapi/listing_pb";
import { Quote } from "../serverapi/market-data-service_pb";
import { StaticDataServiceClient } from "../serverapi/Static-data-serviceServiceClientPb";
import { ListingIds, Listings } from "../serverapi/static-data-service_pb";
import { QuoteListener, QuoteService } from "../services/QuoteService";
import { toNumber } from "../util/decimal64Conversion";
import { ListingContext, TicketController } from "./Container";
import InstrumentListingSearchBar from "./InstrumentListingSearchBar";
import Login from "./Login";
import { MenuItem } from "react-contextmenu";
import { Menu,   Colors,  Hotkeys, Hotkey } from '@blueprintjs/core';
import './OrderBlotter.css';
import { Side } from "../serverapi/order_pb";


interface InstrumentListingWatchState {
  watches: ListingWatch[]
}

interface InstrumentListingWatchProps {
  node: TabNode,
  model: Model,
  quoteService: QuoteService,
  listingContext: ListingContext,
  ticketController: TicketController
}

interface PersistentConfig {
  listingIds: number[]
}

export default class InstrumentListingWatch extends React.Component<InstrumentListingWatchProps, InstrumentListingWatchState> implements QuoteListener {

  staticDataService = new StaticDataServiceClient(Login.grpcContext.serviceUrl, null, null)
  quoteService: QuoteService
  listingContext: ListingContext
  ticketController: TicketController

  watchMap: Map<number, ListingWatch> = new Map()

  constructor(props: InstrumentListingWatchProps) {
    super(props);

    this.quoteService = props.quoteService
    this.ticketController = props.ticketController

    let initialState: InstrumentListingWatchState = {
      watches: Array.from(this.watchMap.values())
    }

    this.state = initialState;

    this.addListing = this.addListing.bind(this);

    this.props.node.setEventListener("save", (p) => {
      let persistentConfig: PersistentConfig = { listingIds: Array.from(this.state.watches.map(l => l.listing.getId())) }
      this.props.model.doAction(Actions.updateNodeAttributes(props.node.getId(), { config: persistentConfig }))
    });

    if (this.props.node.getConfig() && this.props.node.getConfig()) {
      let persistentConfig: PersistentConfig = this.props.node.getConfig();
      this.addListingLineByIds(persistentConfig.listingIds)
    }

    this.listingContext = props.listingContext

    this.openBuyDialog = this.openBuyDialog.bind(this);
    this.openSellDialog = this.openSellDialog.bind(this);
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
        if (listings) {
          listings.getListingsList().forEach((listing) => {
            this.addListingLine(listing)
          })
        }
      })
  }

  private addListingLine(listing: Listing) {

    if (!this.watchMap.get(listing.getId())) {
      let line = new ListingWatch(listing)


      this.watchMap.set(listing.getId(), line);

      let lines = this.state.watches.slice(0)
      lines.push(line)

      this.setState({
        watches: lines
      });

      this.quoteService.SubscribeToQuote(listing.getId(), this)

    }

  }

  onQuote(quote: Quote): void {

    let line = this.watchMap.get(quote.getListingid())
    if (line) {

      line.quote = quote
      let lines = this.state.watches.slice(0)
      this.setState({
        watches: lines
      });

    } else {
      logDebug("received quote update for non-existent watch, quote id:" + quote.getListingid())
    }

  }

  public render() {

    return (

      <div className="bp3-dark">

        <InstrumentListingSearchBar add={this.addListing} />
        <Table enableRowResizing={false} numRows={this.state.watches.length} className="bp3-dark" selectionModes={SelectionModes.ROWS_AND_CELLS}
          onSelection={this.onSelection} bodyContextMenuRenderer={this.renderContextMenu}>
          <Column name="Id" cellRenderer={this.renderId} />
          <Column name="Symbol" cellRenderer={this.renderSymbol} />
          <Column name="Name" cellRenderer={this.renderName} />
          <Column name="Mic" cellRenderer={this.renderMic} />
          <Column name="Country" cellRenderer={this.renderCountry} />
          <Column name="Bid Size" cellRenderer={this.renderBidSize} />
          <Column name="Bid Px" cellRenderer={this.renderBidPrice} />
          <Column name="Ask Px" cellRenderer={this.renderAskPrice} />
          <Column name="Ask Size" cellRenderer={this.renderAskSize} />
        </Table>
      </div>
    );
  }

  private renderId = (row: number) => <Cell>{this.state.watches[row].Id()}</Cell>;
  private renderSymbol = (row: number) => <Cell>{this.state.watches[row].Symbol()}</Cell>;
  private renderName = (row: number) => <Cell>{this.state.watches[row].Name()}</Cell>;
  private renderMic = (row: number) => <Cell>{this.state.watches[row].Mic()}</Cell>;
  private renderCountry = (row: number) => <Cell>{this.state.watches[row].Country()}</Cell>;
  private renderBidSize = (row: number) => <Cell>{this.state.watches[row].BidSize()}</Cell>;
  private renderBidPrice = (row: number) => <Cell>{this.state.watches[row].BidPrice()}</Cell>;
  private renderAskPrice = (row: number) => <Cell>{this.state.watches[row].AskPrice()}</Cell>;
  private renderAskSize = (row: number) => <Cell>{this.state.watches[row].AskPrice()}</Cell>;

  renderContextMenu = () => {
    return (

      <Menu >
        <MenuItem  onClick={this.openBuyDialog} disabled={this.listingContext.selectedListing === undefined}>
          Buy
         </MenuItem>
        <MenuItem divider />
        <MenuItem onClick={this.openSellDialog} disabled={this.listingContext.selectedListing === undefined}>
          Sell
         </MenuItem>
      </Menu>

    );
  };

  private openBuyDialog(e: React.TouchEvent<HTMLDivElement> | React.MouseEvent<HTMLDivElement>) {

    if (this.listingContext.selectedListing) {
      this.ticketController.openTicket(Side.BUY, this.listingContext.selectedListing)
    }

  }

  private openSellDialog(e: any) {
    if (this.listingContext.selectedListing) {
      this.ticketController.openTicket(Side.SELL, this.listingContext.selectedListing)
    }
  }

  private onSelection = (selectedRegions: IRegion[]) => {
    let selectedWatches: Map<number, ListingWatch> = new Map<number, ListingWatch>()

    for (let region of selectedRegions) {

      let firstRowIdx: number;
      let lastRowIdx: number;

      if (region.rows) {
        firstRowIdx = region.rows[0]
        lastRowIdx = region.rows[1]
      } else {
        firstRowIdx = 0
        lastRowIdx = this.state.watches.length - 1
      }

      for (let i = firstRowIdx; i <= lastRowIdx; i++) {
        let watch = this.state.watches[i]
        selectedWatches.set(watch.Id(), watch)

        if (i === firstRowIdx) {
          this.listingContext.setSelectedListing(watch.listing)
        }

      }
    }

  }


}

class ListingWatch {

  listing: Listing;
  quote?: Quote;

  constructor(listing: Listing) {
    this.listing = listing
  }

  Id(): number {
    return this.listing.getId()
  }

  Symbol(): string {
    let i = this.listing.getInstrument()
    if (i) {
      return i.getDisplaysymbol()
    }

    return ""
  }

  Name(): string {
    let i = this.listing.getInstrument()
    if (i) {
      return i.getName()
    }

    return ""
  }

  Mic(): string {
    let m = this.listing.getMarket()
    if (m) {
      return m.getMic()
    }

    return ""
  }

  Country(): string {
    let m = this.listing.getMarket()
    if (m) {
      return m.getCountrycode()
    }

    return ""
  }

  BidSize(): string {
    if (this.quote) {
      if (this.quote.getDepthList().length >= 1) {
        let depth = this.quote.getDepthList()[0]
        let sz = toNumber(depth.getBidsize())
        if (sz) {
          return sz.toString()
        }
      }
    }

    return ""
  }

  BidPrice(): string {
    if (this.quote) {
      if (this.quote.getDepthList().length >= 1) {
        let depth = this.quote.getDepthList()[0]
        let sz = toNumber(depth.getBidprice())
        if (sz) {
          return sz.toString()
        }
      }
    }

    return ""
  }

  AskSize(): string {
    if (this.quote) {
      if (this.quote.getDepthList().length >= 1) {
        let depth = this.quote.getDepthList()[0]
        let sz = toNumber(depth.getAsksize())
        if (sz) {
          return sz.toString()
        }
      }
    }

    return ""
  }

  AskPrice(): string {
    if (this.quote) {
      if (this.quote.getDepthList().length >= 1) {
        let depth = this.quote.getDepthList()[0]
        let sz = toNumber(depth.getAskprice())
        if (sz) {
          return sz.toString()
        }
      }
    }

    return ""
  }

}
