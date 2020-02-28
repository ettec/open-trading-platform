import { Menu } from '@blueprintjs/core';
import { Cell, Column, IRegion, SelectionModes, Table } from "@blueprintjs/table";
import { Actions, Model, TabNode } from "flexlayout-react";
import React from 'react';
import { MenuItem } from "react-contextmenu";
import { logDebug } from "../logging/Logging";
import { Listing } from "../serverapi/listing_pb";
import { Side } from "../serverapi/order_pb";
import { ListingId } from "../serverapi/static-data-service_pb";
import { ListingService } from "../services/ListingService";
import { QuoteListener, QuoteService } from "../services/QuoteService";
import { toNumber } from "../util/decimal64Conversion";
import { ListingContext, TicketController } from "./Container";
import InstrumentListingSearchBar from "./InstrumentListingSearchBar";
import './TableCommon.css';
import { ClobQuote } from '../serverapi/clobquote_pb';


interface InstrumentListingWatchState {
  watches: Array<ListingWatch>
}

interface InstrumentListingWatchProps {
  node: TabNode,
  model: Model,
  quoteService: QuoteService,
  listingContext: ListingContext,
  ticketController: TicketController,
  listingService: ListingService
}

interface PersistentConfig {
  listingIds: number[]
}

export default class InstrumentListingWatch extends React.Component<InstrumentListingWatchProps, InstrumentListingWatchState> implements QuoteListener {

  quoteService: QuoteService
  listingContext: ListingContext
  ticketController: TicketController
  listingService: ListingService

  watchMap: Map<number, ListingWatch> = new Map()

  constructor(props: InstrumentListingWatchProps) {
    super(props);

    this.quoteService = props.quoteService
    this.ticketController = props.ticketController
    this.listingService = props.listingService

    this.addListing = this.addListing.bind(this);

    this.props.node.setEventListener("save", (p) => {
      let persistentConfig: PersistentConfig = { listingIds: Array.from(this.state.watches.map(l => l.listingId)) }
      this.props.model.doAction(Actions.updateNodeAttributes(props.node.getId(), { config: persistentConfig }))
    });

    let initialWatches = new Array<ListingWatch>()
    if (this.props.node.getConfig() && this.props.node.getConfig()) {
      let persistentConfig: PersistentConfig = this.props.node.getConfig();
      persistentConfig.listingIds.forEach(id => {
        let watch = this.getListingWatch(id)
        initialWatches.push(watch)
      })
    }

    let initialState: InstrumentListingWatchState = {
      watches: initialWatches
    }
    this.state = initialState

    this.listingContext = props.listingContext
    this.openBuyDialog = this.openBuyDialog.bind(this);
    this.openSellDialog = this.openSellDialog.bind(this);
  }

  addListing(listing?: Listing) {

    if (listing) {
      if (!this.watchMap.get(listing.getId())) {
        let watch = this.getListingWatch(listing.getId())
        let newWatches = this.state.watches.slice(0)
        newWatches.push(watch)
        this.setState(
          {
            watches: newWatches
          }
        )
      }
    }
  }



  private getListingWatch(listingId: number): ListingWatch {


    let line = new ListingWatch(listingId)

    let id = new ListingId()
    id.setListingid(listingId)

    this.listingService.GetListing(listingId, (listing: Listing)=> {
      line.listing = listing

      this.setState({
        watches: Object.assign([], this.state.watches)
      })
    })
  
    this.watchMap.set(listingId, line);
    this.quoteService.SubscribeToQuote(listingId, this)

    return line

  }

  onQuote(quote: ClobQuote): void {

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
  private renderAskSize = (row: number) => <Cell>{this.state.watches[row].AskSize()}</Cell>;

  renderContextMenu = () => {
    return (

      <Menu >
        <MenuItem onClick={this.openBuyDialog} disabled={this.listingContext.selectedListing === undefined}>
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
          if (watch.listing) {
            this.listingContext.setSelectedListing(watch.listing)
          }

        }

      }
    }

  }


}

class ListingWatch {

  listingId: number;
  listing?: Listing;
  quote?: ClobQuote;

  constructor(listingId: number) {
    this.listingId = listingId
  }

  Id(): number {
    return this.listingId
  }

  Symbol(): string {
    if (this.listing) {
      let i = this.listing.getInstrument()
      if (i) {
        return i.getDisplaysymbol()
      }
    }

    return ""
  }

  Name(): string {
    if (this.listing) {
      let i = this.listing.getInstrument()
      if (i) {
        return i.getName()
      }
    }


    return ""
  }

  Mic(): string {
    if (this.listing) {
      let m = this.listing.getMarket()
      if (m) {
        return m.getMic()
      }
    }
    return ""
  }

  Country(): string {
    if (this.listing) {
      let m = this.listing.getMarket()
      if (m) {
        return m.getCountrycode()
      }
    }
    return ""
  }

  BidSize(): string {
    if (this.quote) {
      if (this.quote.getBidsList().length >= 1) {
        let depth = this.quote.getBidsList()[0]
        let sz = toNumber(depth.getSize())
        if (sz) {
          return sz.toString()
        }
      }
    }

    return ""
  }

  BidPrice(): string {
    if (this.quote) {
      if (this.quote.getBidsList().length >= 1) {
        let depth = this.quote.getBidsList()[0]
        let sz = toNumber(depth.getPrice())
        if (sz) {
          return sz.toString()
        }
      }
    }

    return ""
  }

  AskSize(): string {
    if (this.quote) {
      if (this.quote.getOffersList().length >= 1) {
        let depth = this.quote.getOffersList()[0]
        let sz = toNumber(depth.getSize())
        if (sz) {
          return sz.toString()
        }
      }
    }

    return ""
  }

  AskPrice(): string {
    if (this.quote) {
      if (this.quote.getOffersList().length >= 1) {
        let depth = this.quote.getOffersList()[0]
        let sz = toNumber(depth.getPrice())
        if (sz) {
          return sz.toString()
        }
      }
    }

    return ""
  }

}
