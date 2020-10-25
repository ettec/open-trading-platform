import { AnchorButton, Button, Classes, Dialog, Intent } from '@blueprintjs/core';
import { Cell, SelectionModes, Table } from "@blueprintjs/table";
import { ColDef, ColumnApi, ColumnState, GridReadyEvent } from 'ag-grid-community';
import { AgGridReact } from 'ag-grid-react/lib/agGridReact';
import * as grpcWeb from 'grpc-web';
import * as React from "react";
import { Timestamp } from '../../serverapi/modelcommon_pb';
import { OrderHistory } from '../../serverapi/orderdataservice_pb';
import { Order } from '../../serverapi/order_pb';
import { ListingService } from '../../services/ListingService';
import { OrderService } from "../../services/OrderService";
import { OrderHistoryBlotterController } from "../Container/Controllers";
import { TableViewProperties } from "../TableView/TableView";
import { OrderView } from "./OrderView";
import { CountryFlagRenderer } from "./ParentOrderBlotterAgGrid"

export interface OrderHistoryBlotterProps extends TableViewProperties {
    orderService: OrderService
    listingService: ListingService
    orderHistoryBlotterController: OrderHistoryBlotterController
}

interface OrderHistoryBlotterState {
    orders: OrderView[],
    isOpen: boolean,
    usePortal: boolean
    order?: Order
    updates: OrderUpdateView[];
    width: number,
    colStates: ColumnState[],
    colDefs: ColDef[],
}

const fieldName = (name: keyof OrderUpdateView) => name;

export default class OrderHistoryBlotter extends React.Component<OrderHistoryBlotterProps, OrderHistoryBlotterState> {

    gridColumnApi?: ColumnApi;
    orderService: OrderService
    listingService: ListingService
    orderHistoryBlotterController: OrderHistoryBlotterController


    constructor(props: OrderHistoryBlotterProps) {
        super(props)

        this.orderService = props.orderService
        this.listingService = props.listingService
        this.orderHistoryBlotterController = props.orderHistoryBlotterController

        this.orderHistoryBlotterController.setBlotter(this)

        this.state = {
            isOpen: false,
            usePortal: false,
            updates: new Array<OrderUpdateView>(0),
            orders: new Array<OrderView>(),
            width: 0,
            colStates: new Array<ColumnState>(),
            colDefs: new Array<ColDef>(),
        }

        this.onGridReady = this.onGridReady.bind(this);
    }

    protected getTableName(): string {
        return "Child Orders"
    }

    getTitle(order?: Order): string {
        if (order) {
            return "History of order " + order.getId()
        }

        return ""
    }

    onGridReady(params: GridReadyEvent) {

        this.gridColumnApi = params.columnApi;
        this.gridColumnApi.setColumnState(this.state.colStates)

    }



    render() {
        return (
            <Dialog
                icon="bring-data"
                onClose={this.handleClose}
                title={this.getTitle(this.state.order)}
                style={{ minWidth: this.state.width }}
                {...this.state}
                className="bp3-dark">
                <div className={Classes.DIALOG_BODY} >
                    <div className="ag-theme-balham-dark"  >
                        <AgGridReact
                        domLayout='autoHeight' 
                            frameworkComponents={{
                                countryFlagRenderer: CountryFlagRenderer,
                            }}

                            defaultColDef={{
                                sortable: false,
                                filter: true,
                                resizable: true
                            }}
                            columnDefs={this.state.colDefs}
                            rowData={this.state.orders}
                            onGridReady={this.onGridReady}
                        />
                    </div>

                </div>
                <div className={Classes.DIALOG_FOOTER}>
                    <div className={Classes.DIALOG_FOOTER_ACTIONS}>
                        <AnchorButton onClick={this.handleClose}
                            intent={Intent.PRIMARY}>Close
                        </AnchorButton>
                    </div>
                </div>
            </Dialog>
        )
    }







    open(order: Order, colStates: ColumnState[], colDefs: ColDef[], width: number) {



        this.orderService.GetOrderHistory(order, (err: grpcWeb.Error, history: OrderHistory) => {

            let newViews = new Array<OrderUpdateView>()
            for (let update of history.getUpdatesList()) {

                let order = update.getOrder()
                let time = update.getTime()
                if (order && time) {
                    let view = new OrderUpdateView(order, time)
                    let listing = this.listingService.GetListingImmediate(order.getListingid())
                    if (listing) {
                        view.setListing(listing)
                    }

                    newViews.push(view)
                }
            }


            let blotterState = {
                ...this.state, ...{
                    updates: newViews,
                    orders: newViews
                }
            }

            this.setState(blotterState)

        })

        let colDefsWithTime = new Array<ColDef>()

        

        let timeCol: ColDef = {
            headerName: 'Time',
            field: fieldName('time'),
            width: 80,

        };

        colDefsWithTime.push(timeCol)
        colDefs.forEach(d=>colDefsWithTime.push(d)) 
        
        let state: OrderHistoryBlotterState =
        {
            order: order,
            isOpen: true,
            usePortal: false,

            updates: new Array<OrderUpdateView>(),
            orders: new Array<OrderView>(),
            width: width,
            colStates: colStates,
            colDefs: colDefsWithTime,
        }

        this.setState(state)
    }


    private handleClose = () => {
        this.setState({
            ...this.state, ...{
                isOpen: false,
            }
        })
    };


}

class OrderUpdateView extends OrderView {

    time?: string;

    constructor(order: Order, updateTime: Timestamp) {
        super(order)

        if (updateTime) {
            let date = new Date(updateTime.getSeconds() * 1000)
            this.time = date.toLocaleTimeString()
        }
    }

}
