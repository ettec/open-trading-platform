import { Button, Position, Tooltip } from '@blueprintjs/core';
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
import { CountryFlagRenderer, DirectionalPriceRenderer } from '../TableView/Renderers';
import '../TableView/TableCommon.css';
import InstrumentListingSearchBar from "./InstrumentListingSearchBar";
import { InstrumentListingWatchesView, ListingWatchView, WatchEventType } from './InstrumentListingWatchView';


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
  colController: AgGridColumnChooserController

  gridApi?: GridApi;
  gridColumnApi?: ColumnApi;

  private watchesView: InstrumentListingWatchesView

  constructor(props: InstrumentListingWatchProps) {
    super(props);

    this.quoteService = props.quoteService
    this.ticketController = props.ticketController
    this.listingService = props.listingService
    this.colController = props.colController

    this.watchesView = new InstrumentListingWatchesView(props.listingService, props.quoteService, (watches: ListingWatchView[], eventType: WatchEventType) => {

      if (this.gridApi) {

        switch (eventType) {
          case WatchEventType.Add:
            let rt: RowDataTransaction = {
              add: watches
            }

            this.gridApi.applyTransaction(rt)
            let addedNode = this.gridApi.getRowNode(watches[watches.length - 1].listingId.toString())

            this.gridApi.ensureIndexVisible(addedNode.rowIndex, 'bottom')
            break;
          case WatchEventType.Update:
            for (let watch of watches) {
              let updateNode = this.gridApi.getRowNode(watch.listingId.toString())
              if (updateNode) {
                let rf: RefreshCellsParams = {
                  rowNodes: [updateNode]
                }

                this.gridApi?.refreshCells(rf)
              }
            }
            break;
          case WatchEventType.Remove:
            let removetx: RowDataTransaction = {
              remove: watches
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

  editVisibleColumns = () => {

    if (this.gridColumnApi) {


      this.colController.open(this.getTableName(), this.gridColumnApi.getColumnState(),
        this.gridColumnApi.getAllColumns(), (newColumnsState: ColumnState[] | undefined) => {
          if (newColumnsState) {
            let colState: ApplyColumnStateParams = {
              state: newColumnsState,
              applyOrder: true
            }
            this.gridColumnApi?.applyColumnState(colState)
          }
        })
    }

  }

  public render() {

    return (


      <div style={{ width: "100%", height: "100%", display: 'flex', flexDirection: 'column', alignItems: "centre" }}>
        <div className="bp3-dark" style={{ display: 'flex', flexDirection: 'row', paddingTop: 0, alignItems: "left" }}>
          <div style={{ flexGrow: 1, flexDirection: 'row', display: 'flex' }}>
            <InstrumentListingSearchBar add={this.addListing} />

            <Button minimal={true} text="Remove" onClick={this.removeListings} disabled={this.state.selectedWatches.length === 0} />
            <span style={{minWidth:40}}></span>
            <Button  text="Buy" onClick={this.openBuyDialog} disabled={this.state.selectedWatches.length === 0} 
            style={{ minWidth: 80,  backgroundColor:GlobalColours.BUYBKG  }} />
            <span style={{minWidth:5}}></span>
            <Button  text="Sell" onClick={this.openSellDialog} disabled={this.state.selectedWatches.length === 0}
            style={{ minWidth: 80, backgroundColor:GlobalColours.SELLBKG  }} />
            

          </div>
          <div >
            <Tooltip
              content={<span>Edit Visible Columns</span>}
              position={Position.LEFT_BOTTOM}
              usePortal={false}

            >
              <Button minimal={true} icon="manually-entered-data" onClick={() => this.editVisibleColumns()} />
            </Tooltip>
          </div>
        </div>

        <div className="ag-theme-balham-dark" style={{ width: "100%", height: "100%" }}>
          <AgGridReact
            rowSelection={'multiple'}
            getRowNodeId={(data: ListingWatchView) => { return data?.listingId?.toString() }}
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


  private removeListings() {
    this.watchesView.removeListings(this.state.selectedWatches.map(w => w.listingId))
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



