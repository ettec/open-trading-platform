import { AnchorButton, Classes, Dialog, Intent } from '@blueprintjs/core';
import { ColDef, ColumnApi, ColumnState, GridReadyEvent } from 'ag-grid-community';
import { AgGridReact } from 'ag-grid-react/lib/agGridReact';
import * as React from "react";
import { Order } from '../../serverapi/order_pb';
import { ListingService } from '../../services/ListingService';
import { OrderService } from "../../services/OrderService";
import { ChildOrderBlotterController } from "../Container/Controllers";
import { OrderView } from "./OrderView";
import { CountryFlagRenderer } from './ParentOrderBlotterAgGrid';

export interface ChildOrderProps {
    orderService: OrderService
    listingService: ListingService
    childOrderBlotterController: ChildOrderBlotterController
}


interface ChildOrderBlotterState {
    isOpen: boolean,
    usePortal: boolean
    orders: Array<OrderView>
    parentOrder?: Order
    width: number,
    colStates: ColumnState[],
    colDefs: ColDef[],
}


export default class ChildOrderBlotter extends React.Component<ChildOrderProps, ChildOrderBlotterState> {

    gridColumnApi?: ColumnApi;
    orderService: OrderService
    listingService: ListingService
    childOrderBlotterController: ChildOrderBlotterController

    constructor(props: ChildOrderProps) {
        super(props)

        this.orderService = props.orderService
        this.listingService = props.listingService
        this.childOrderBlotterController = props.childOrderBlotterController

        this.childOrderBlotterController.setBlotter(this)

        this.state = {
            isOpen: false,
            usePortal: false,
            orders: new Array<OrderView>(10),
            width: 0,
            colStates: new Array<ColumnState>(),
            colDefs: new Array<ColDef>(),
            parentOrder: undefined
        }

        this.onGridReady = this.onGridReady.bind(this);
    }

    protected getTableName(): string {
        return "Child Orders"
    }


    getTitle(order?: Order): string {
        if (order) {
            return "Children of order " + order.getId()
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
                icon="fork"
                onClose={this.handleClose}
                title={this.getTitle(this.state.parentOrder)}
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
                                sortable: true,
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





    open(parentOrder: Order, orders: Array<Order>, colStates: ColumnState[], colDefs: ColDef[], width: number) {

        let orderViews = Array<OrderView>()
        for (let order of orders) {
            let view = new OrderView(order, ()=>{}, ()=>{})
            let listing = this.listingService.GetListingImmediate(order.getListingid())
            if (listing) {
                view.setListing(listing)
            }

            orderViews.push(view)
        }


        let state =
        {
            parentOrder: parentOrder,
            isOpen: true,
            usePortal: false,
            colStates: colStates,
            colDefs: colDefs,
            width: width,
            orders: orderViews
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
