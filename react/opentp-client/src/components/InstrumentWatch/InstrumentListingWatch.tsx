import { Icon } from '@blueprintjs/core';
import { ColDef, ColumnApi, ColumnState, GridApi, GridReadyEvent, RefreshCellsParams, RowDataTransaction } from 'ag-grid-community';
import { ApplyColumnStateParams } from 'ag-grid-community/dist/lib/columnController/columnApi';
import { AgGridReact } from 'ag-grid-react/lib/agGridReact';
import { Actions, Model, TabNode } from "flexlayout-react";
import React, { Component } from 'react';
import { Listing } from "../../serverapi/listing_pb";
import { Side } from "../../serverapi/order_pb";
import { ListingService } from "../../services/ListingService";
import { QuoteService } from "../../services/QuoteService";
import { GlobalColours } from '../Container/Colours';
import { ListingContext } from "../Container/Contexts";
import { AgGridColumnChooserController, TicketController } from "../Container/Controllers";
import { CountryFlagRenderer } from '../OrderBlotter/ParentOrderBlotterAgGrid';
import '../TableView/TableCommon.css';
import InstrumentListingSearchBar from "./InstrumentListingSearchBar";
import { InstrumentListingWatchesView, ListingWatchView, WatchEventType } from './InstrumentListingWatchView';


export class DirectionalPriceRenderer extends Component<any, any> {
  constructor(props: any) {
    super(props);

    this.state = {
      value: this.props.value,
    };
  }


  render() {

    let price = this.state?.value?.price
    let direction = this.state?.value?.direction

    if (price) {
      if (direction) {
        if (direction > 0) {
          return <span><Icon icon="arrow-up" style={{ color: GlobalColours.UPTICK }} />{price}</span>;
        }

        if (direction < 0) {
          return <span><Icon icon="arrow-down" style={{ color: GlobalColours.DOWNTICK }} />{price}</span>;
        }

        return <span>{price}</span>;
      } else {
        return <span>{price}</span>;
      }
    } else {
      return <span></span>
    }
  }
}





const fieldName = (name: keyof ListingWatchView) => name;

const columnDefs: ColDef[] = [

  {
    headerName: 'Symbol',
    field: fieldName('symbol'),
    width: 80,
    rowDrag: true
  },

  {
    headerName: 'Id',
    field: fieldName('listingId'),
    width: 170,

  },

  {
    headerName: 'Name',
    field: fieldName('name'),
    width: 80,

  },
  {
    headerName: 'Mic',
    field: fieldName('mic'),
    width: 80,

  },
  {
    headerName: 'Country',
    field: fieldName('countryCode'),
    width: 80,
    cellRenderer: 'countryFlagRenderer'
  },
  {
    headerName: 'Bid Qty',
    field: fieldName('bidSize'),
    width: 80,
  },

  {
    headerName: 'Bid Px',
    field: fieldName('bidPrice'),
    width: 80,
    cellRenderer: 'directionalPriceRenderer'
  },

  {
    headerName: 'Ask Qty',
    field: fieldName('askSize'),
    width: 80,
  },

  {
    headerName: 'Ask Px',
    field: fieldName('askPrice'),
    width: 80,
    cellRenderer: 'directionalPriceRenderer'
  },
  {
    headerName: 'Last Qty',
    field: fieldName('lastSize'),
    width: 80,
  },

  {
    headerName: 'Last Px',
    field: fieldName('lastPrice'),
    width: 80,
    cellRenderer: 'directionalPriceRenderer'
  },
  {
    headerName: 'Traded Vol',
    field: fieldName('tradedVolume'),
    width: 80,
  },

]

interface InstrumentListingWatchState {
  selectedWatches: Array<ListingWatchView>
}

interface InstrumentListingWatchProps {
  node: TabNode,
  model: Model,
  quoteService: QuoteService,
  listingContext: ListingContext,
  ticketController: TicketController,
  listingService: ListingService,
  colController: AgGridColumnChooserController
}


interface InstrumentWatchPersistentConfig {
  colState: ColumnState[]
  listingIds: number[]
}

export default class InstrumentListingWatch extends Component<InstrumentListingWatchProps, InstrumentListingWatchState>  {

  quoteService: QuoteService
  listingContext: ListingContext
  ticketController: TicketController
  listingService: ListingService

  gridApi?: GridApi;
  gridColumnApi?: ColumnApi;


  watchMap: Map<number, ListingWatchView> = new Map()

  private watchesView: InstrumentListingWatchesView




  constructor(props: InstrumentListingWatchProps) {
    super(props);

    this.quoteService = props.quoteService
    this.ticketController = props.ticketController
    this.listingService = props.listingService

    this.watchesView = new InstrumentListingWatchesView(props.listingService, props.quoteService, (watch: ListingWatchView, eventType: WatchEventType) => {

      if (this.gridApi) {

        switch (eventType) {
          case WatchEventType.Add:
            let rt: RowDataTransaction = {
              add: [watch]
            }

            this.gridApi.applyTransaction(rt)
            let addedNode = this.gridApi.getRowNode(watch.listingId.toString())

            this.gridApi.ensureIndexVisible(addedNode.rowIndex, 'bottom')
            break;
          case WatchEventType.Update:
            let updateNode = this.gridApi.getRowNode(watch.listingId.toString())
            if (updateNode) {
              let rf: RefreshCellsParams = {
                rowNodes: [updateNode]
              }

              this.gridApi?.refreshCells(rf)
            }
            break;
          case WatchEventType.Remove:
            let removetx: RowDataTransaction = {
              remove: [watch]
            }

            this.gridApi.applyTransaction(removetx)
            break;
        }

      }
    })


    this.props.node.setEventListener("save", (p) => {

      let colState = new Array<ColumnState>()
      if (this.gridColumnApi) {
        colState = this.gridColumnApi.getColumnState()
      }

      let listingIds = new Array<number>()
      if (this.gridApi) {
        this.gridApi.forEachNode(n => { listingIds.push(n.data.listingId) })
      }

      let persistentConfig: InstrumentWatchPersistentConfig = {
        colState: colState,
        listingIds: listingIds
      }

      this.props.model.doAction(Actions.updateNodeAttributes(props.node.getId(), { config: persistentConfig }))
    });

    let initialState: InstrumentListingWatchState = {
      selectedWatches: new Array<ListingWatchView>()
    }
    this.state = initialState

    this.listingContext = props.listingContext
    this.addListing = this.addListing.bind(this);
    this.openBuyDialog = this.openBuyDialog.bind(this);
    this.openSellDialog = this.openSellDialog.bind(this);
    this.removeListings = this.removeListings.bind(this);
    this.onGridReady = this.onGridReady.bind(this);
    this.onSelectionChanged = this.onSelectionChanged.bind(this)


  }

  onGridReady(params: GridReadyEvent) {
    this.gridApi = params.api;
    this.gridColumnApi = params.columnApi;


    let config: InstrumentWatchPersistentConfig = this.props.node.getConfig()
    let initialColConfig = config.colState

    if (initialColConfig) {
      let colState: ApplyColumnStateParams = {
        state: initialColConfig,
        applyOrder: true
      }
      this.gridColumnApi.applyColumnState(colState)
    }

    let listingIds = config.listingIds
    if (listingIds) {
      listingIds.forEach(id => this.watchesView.addListing(id))
    }
  };

  addListing(listing?: Listing) {
    if (listing) {
      this.watchesView.addListing(listing.getId())
    }
  }

  protected getTableName(): string {
    return "Instrument Watch"
  }

  public render() {

    return (

      <div className="bp3-dark">
        <InstrumentListingSearchBar add={this.addListing} />
        <div className="ag-theme-balham-dark" style={{ width: "100%", height: "100%" }}>
          <AgGridReact
            rowSelection={'multiple'}
            getRowNodeId={(data: ListingWatchView) => { return data.listingId.toString() }}
            onSelectionChanged={this.onSelectionChanged}
            domLayout='autoHeight'
            rowDragManaged={true}

            frameworkComponents={{
              countryFlagRenderer: CountryFlagRenderer,
              directionalPriceRenderer: DirectionalPriceRenderer
            }}

            defaultColDef={{
              sortable: false,
              filter: true,
              resizable: true
            }}
            columnDefs={columnDefs}
            onGridReady={this.onGridReady}
          />
        </div>
      </div>

    );
  }

  /*
  renderContextMenu = () => {
    return (

      <Menu >
        
        <Menu.Item icon="arrow-left" text="Buy" onClick={this.openBuyDialog} disabled={this.listingContext.selectedListing === undefined}>
        </Menu.Item>
        <Menu.Divider />
        <Menu.Item icon="arrow-right" text="Sell" onClick={this.openSellDialog} disabled={this.listingContext.selectedListing === undefined}>
        </Menu.Item>
        <Menu.Divider />
        <Menu.Item text="Move Listings Up" onClick={this.moveListingsUp} disabled={this.listingContext.selectedListing === undefined}>
        </Menu.Item>
        <Menu.Item text="Move Listings Down" onClick={this.moveListingsDown} disabled={this.listingContext.selectedListing === undefined}>
        </Menu.Item>
        <Menu.Item text="Remove Listings" onClick={this.removeListings} disabled={this.listingContext.selectedListing === undefined}>
        </Menu.Item>
        <Menu.Divider />
        <Menu.Item text="Edit Visible Columns" onClick={() => this.editVisibleColumns()}  >
        </Menu.Item>
        
      </Menu>

    );
  };
  */

  private removeListings() {

    /*
    for (let watch of this.state.selectedWatches) {
      this.watchMap.delete(watch.listingId)
      this.quoteService.UnsubscribeFromQuote(watch.listingId, this)
    }

    let newWatches = Array<ListingWatchView>()

    for (let watch of this.state.watches) {
      if (this.watchMap.has(watch.listingId)) {
        newWatches.push(watch)
      }
    }

    let blotterState = {
      ...this.state, ...{
        watches: newWatches
      }
    }

    this.setState(blotterState)
    */
  }




  private openBuyDialog() {

    if (this.listingContext.selectedListing) {
      this.ticketController.openNewOrderTicket(Side.BUY, this.listingContext.selectedListing)
    }

  }

  private openSellDialog() {
    if (this.listingContext.selectedListing) {
      this.ticketController.openNewOrderTicket(Side.SELL, this.listingContext.selectedListing)
    }
  }


  onSelectionChanged() {

    let selectedWatches = new Array<ListingWatchView>()
    if (this.gridApi) {
      var selectedRows = this.gridApi.getSelectedRows();

      selectedRows.forEach(function (selectedRow, index) {

        let watch: ListingWatchView = selectedRow

        selectedWatches.push(watch)

      });

    }

    let newState: InstrumentListingWatchState = {
      ...this.state, ...{
        selectedWatches: selectedWatches,
      }
    }

    this.setState(newState)

    if (selectedWatches.length > 0) {
      let listing = selectedWatches[0]?.getListing()
      if (listing)
        this.listingContext.setSelectedListing(listing)
    }
  }



}



