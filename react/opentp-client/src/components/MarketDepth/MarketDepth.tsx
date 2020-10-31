import { Button, Position, Tooltip } from "@blueprintjs/core";
import { ColDef, ColumnApi, ColumnState, GridApi, GridReadyEvent, RowDataTransaction } from "ag-grid-community";
import { ApplyColumnStateParams } from "ag-grid-community/dist/lib/columnController/columnApi";
import { AgGridReact } from "ag-grid-react/lib/agGridReact";
import { Actions, Model, TabNode } from "flexlayout-react";
import * as grpcWeb from 'grpc-web';
import React, { Component } from 'react';
import { getListingShortName } from "../../common/modelutilities";
import { ClobQuote } from "../../serverapi/clobquote_pb";
import { Listing } from "../../serverapi/listing_pb";
import { Side } from "../../serverapi/order_pb";
import { StaticDataServiceClient } from "../../serverapi/StaticdataserviceServiceClientPb";
import { ListingService } from "../../services/ListingService";
import { QuoteListener, QuoteService } from "../../services/QuoteService";
import { GlobalColours } from "../Container/Colours";
import { ListingContext } from "../Container/Contexts";
import { AgGridColumnChooserController, TicketController } from "../Container/Controllers";
import Login from "../Login";
import { DepthLine } from "./MarketDepthView";
import { MarketDepthView } from "./MarketDepthView";




interface MarketDepthProps {
  node: TabNode,
  model: Model,
  quoteService: QuoteService,
  listingContext: ListingContext
  listingService: ListingService
  ticketController: TicketController,
  colController: AgGridColumnChooserController
}

interface MarketDepthState {
  listing?: Listing,
  locked: boolean
  lockedListingId: number
  selectedIdx?: number
}

interface MarketDepthConfig {
  colState: ColumnState[]
  lockedListingId: number
  locked: boolean
}


const fieldName = (name: keyof DepthLine) => name;

const columnDefs: ColDef[] = [

  {
    headerName: 'Bid Mic',
    field: fieldName('bidMic'),
    width: 80,
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
  },
  {
    headerName: 'Ask Px',
    field: fieldName('askPrice'),
    width: 80,
  },
  {
    headerName: 'Ask Qty',
    field: fieldName('askSize'),
    width: 80,
  },
  {
    headerName: 'Ask Mic',
    field: fieldName('askMic'),
    width: 80,
  },

]


export default class MarketDepth extends Component<MarketDepthProps, MarketDepthState> implements QuoteListener {

  stream?: grpcWeb.ClientReadableStream<ClobQuote>;

  quoteService: QuoteService;
  listingService: ListingService;
  ticketController: TicketController;
  colController: AgGridColumnChooserController;

  staticDataService = new StaticDataServiceClient(Login.grpcContext.serviceUrl, null, null)

  depth: MarketDepthView
  gridApi?: GridApi;
  gridColumnApi?: ColumnApi;


  constructor(props: MarketDepthProps) {
    super(props);

    this.toggleLock = this.toggleLock.bind(this)

    this.quoteService = props.quoteService
    this.listingService = props.listingService
    this.ticketController = props.ticketController
    this.colController = props.colController

    this.openBuyDialog = this.openBuyDialog.bind(this);
    this.openSellDialog = this.openSellDialog.bind(this);
    this.onGridReady = this.onGridReady.bind(this)
    this.onSelectionChanged = this.onSelectionChanged.bind(this)
    let displayedDepth = 10
    this.depth = new MarketDepthView(displayedDepth, this.listingService, () => {
      if (this.gridApi) {
        this.gridApi.refreshCells()
      }
    })


    let config = this.props.node.getConfig() as MarketDepthConfig

    this.props.node.setEventListener("save", (p) => {


      let lockedListingId = -1
      if (this.state.listing) {
        lockedListingId = this.state.listing.getId()
      }

      let colState = new Array<ColumnState>()
      if (this.gridColumnApi) {
        colState = this.gridColumnApi.getColumnState()
      }

      let persistentConfig: MarketDepthConfig = {
        colState: colState,
        locked: this.state.locked,
        lockedListingId: lockedListingId
      }


      this.props.model.doAction(Actions.updateNodeAttributes(props.node.getId(), { config: persistentConfig }))
    });


    this.state = {
      locked: config?.locked,
      lockedListingId: config?.lockedListingId,
    };
  }

  onGridReady(params: GridReadyEvent) {
    this.gridApi = params.api;
    this.gridColumnApi = params.columnApi;


    let config: MarketDepthConfig = this.props.node.getConfig()
    let initialColConfig = config?.colState

    if (initialColConfig) {
      let colState: ApplyColumnStateParams = {
        state: initialColConfig,
        applyOrder: true
      }
      this.gridColumnApi.applyColumnState(colState)
    }

    let rt: RowDataTransaction = {
      add: this.depth.lines
    }
    this.gridApi.applyTransaction(rt)

    this.props.listingContext.addListener((listing: Listing) => {

      if (!this.state.locked) {

        this.setListing(listing)
      }
    })

    if (this.state.locked) {
      this.listingService.GetListing(this.state.lockedListingId, (response: Listing) => {
        this.setListing(response)
      })
    }

  };


  private setListing(listing: Listing): void {
    if (this.state && this.state.listing) {
      if (this.state.listing === listing) {
        return
      }

      this.quoteService.UnsubscribeFromQuote(this.state.listing.getId(), this)
    }


    let quote = this.quoteService.SubscribeToQuote(listing.getId(), this)

    let state: MarketDepthState = {
      ...this.state, ...{
        listing: listing,
      }
    }

    if (quote) {
      this.depth.setQuote(quote)
    }

    this.setState(state)
  }



  protected getTableName(): string {
    return "Market Depth"
  }

  onQuote(quote: ClobQuote): void {
    this.depth.setQuote(quote)
  }

  toggleLock(): void {
    let state: MarketDepthState = {
      ...this.state, ...{
        locked: !this.state.locked,
      }
    }
    this.setState(state);
  }

  public render() {
    return (
      <div style={{ width: "100%", height: "100%", display: 'flex', flexDirection: 'column', alignItems: "centre" }}>
        <div className="bp3-dark" style={{ display: 'flex', flexDirection: 'row', paddingTop: 0, alignItems: "left" }}>
          <div style={{ flexGrow: 1, flexDirection: 'row', display: 'flex' }}>
            <Button  icon={this.state.locked ? "lock" : "unlock"} onClick={this.toggleLock}>{this.getListingLabel()}</Button>
            <span style={{ minWidth: 40 }}></span>
            <Button text="Buy" onClick={this.openBuyDialog} disabled={!this.state.selectedIdx}
              style={{ minWidth: 80, backgroundColor: GlobalColours.BUYBKG }} />
            <span style={{ minWidth: 5 }}></span>
            <Button text="Sell" onClick={this.openSellDialog} disabled={!this.state.selectedIdx}
              style={{ minWidth: 80, backgroundColor: GlobalColours.SELLBKG }} />
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
            rowSelection={'single'}
            getRowNodeId={(data: DepthLine) => { return data?.idx }}
            onSelectionChanged={this.onSelectionChanged}
            suppressLoadingOverlay={true}
            rowDragManaged={true}
            defaultColDef={{
              sortable: false,
              filter: false,
              resizable: true
            }}
            columnDefs={columnDefs}
            onGridReady={this.onGridReady}
          />
        </div>
      </div>
    );
  }

  onSelectionChanged() {

    if (this.gridApi) {
      
      var selectedRows = this.gridApi.getSelectedRows();
      if (selectedRows.length > 0) {
        let depth: DepthLine = selectedRows[0]
        
        let newState: MarketDepthState = {
          ...this.state, ...{
            selectedIdx: depth.idx,
          }
        }

        this.setState(newState)
      } else {
        let newState: MarketDepthState = {
          ...this.state, ...{
            selectedIdx: undefined,
          }
        }
        this.setState(newState)
      }
    }
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

  private openBuyDialog() {

    let price: number | undefined;
    let quantity: number | undefined;

    if (this.state?.selectedIdx) {
      let result = this.depth.getDepthAtIdx(this.state.selectedIdx, Side.SELL)
      if (result) {
        price = result.price
        quantity = result.quantity
      }

    }

    if (this.state.listing) {
      this.ticketController.openOrderTicketWithDefaultPriceAndQty(Side.BUY, this.state.listing, price, quantity)
    }

  }

  private openSellDialog() {
    let price: number | undefined;
    let quantity: number | undefined;

    if (this.state?.selectedIdx) {
      let result = this.depth.getDepthAtIdx(this.state.selectedIdx, Side.BUY)
      price = result.price
      quantity = result.quantity
    }

    if (this.state.listing) {
      this.ticketController.openOrderTicketWithDefaultPriceAndQty(Side.SELL, this.state.listing, price, quantity)
    }
  }

  private getListingLabel(): string {
    if (this.state && this.state.listing) {
      return getListingShortName(this.state.listing)
    }

    return "(No Selection) "
  }

}



