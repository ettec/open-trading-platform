import { AnchorButton, Classes, Dialog, Intent } from '@blueprintjs/core';
import { ColDef, ColumnApi } from 'ag-grid-community';
import { AgGridReact } from 'ag-grid-react/lib/agGridReact';
import * as grpcWeb from 'grpc-web';
import log from 'loglevel';
import * as React from "react";
import { toNumber } from '../common/decimal64Conversion';
import { Listing } from '../serverapi/listing_pb';
import { Timestamp } from '../serverapi/modelcommon_pb';
import { OrderHistory } from '../serverapi/orderdataservice_pb';
import { Order } from '../serverapi/order_pb';
import { ListingService } from '../services/ListingService';
import { OrderService } from "../services/OrderService";
import { ExecutionsController } from "./Container/Controllers";
import { CountryFlagRenderer } from './OrderBlotter/ParentOrderBlotterAgGrid';



export interface ExecutionsProps {
    orderService: OrderService
    listingService: ListingService
    executionsController: ExecutionsController
}


interface ExecutionsState {
    isOpen: boolean,
    usePortal: boolean
    parentOrder?: Order

    executions: Array<ExecutionView>;
    width: number
}

class ExecutionView {
    time?: string;
    id: string;
    quantity?: number;
    price?: number;
    listing?: Listing;
    symbol: string | undefined;
    mic: string | undefined;
    countryCode: string | undefined;

    constructor(order: Order, updateTime: Timestamp, listing?: Listing) {

        this.id = order.getLastexecid()
        this.quantity = toNumber(order.getLastexecquantity())
        this.price = toNumber(order.getLastexecprice())
        this.listing = listing

        if (updateTime) {
            let date = new Date(updateTime.getSeconds() * 1000)
            this.time = date.toLocaleTimeString()
        }

        this.symbol = this.listing?.getInstrument()?.getDisplaysymbol()
        this.mic = this.listing?.getMarket()?.getMic()
        this.countryCode = this.listing?.getMarket()?.getCountrycode()
    }

}

const fieldName = (name: keyof ExecutionView) => name;

const columnDefs: ColDef[] = [
    {
        headerName: 'Time',
        field: fieldName('time'),
        width: 80,

    },
    {
        headerName: 'Id',
        field: fieldName('id'),
        width: 80
    },
    {
        headerName: 'Symbol',
        field: fieldName('symbol'),
        width: 80
    },
    {
        headerName: 'Mic',
        field: fieldName('mic'),
        width: 80
    },
    {
        headerName: 'Country',
        field: fieldName('countryCode'),
        width: 80,
        cellRenderer: 'countryFlagRenderer'
    },
    {
        headerName: 'Quantity',
        field: fieldName('quantity'),
        width: 80
    },
    {
        headerName: 'Price',
        field: fieldName('price'),
        width: 80
    }

]


export default class Executions extends React.Component<ExecutionsProps, ExecutionsState> {

    gridColumnApi?: ColumnApi;
    orderService: OrderService
    listingService: ListingService
    executionsController: ExecutionsController

    constructor(props: ExecutionsProps) {
        super(props)

        this.orderService = props.orderService
        this.listingService = props.listingService
        this.executionsController = props.executionsController

        this.executionsController.setView(this)

        this.state = {
            isOpen: false,
            usePortal: false,
            executions: new Array<ExecutionView>(10),
            width: 0
        }

    }

    protected getTableName(): string {
        return "Executions"
    }

    getTitle(order?: Order): string {
        if (order) {
            return "Executions for order " + order.getId()
        }

        return ""
    }

    render() {
        return (
            <Dialog
                icon="tick"
                onClose={this.handleClose}
                title={this.getTitle(this.state.parentOrder)}
                style={{ minWidth: this.state.width }}
                {...this.state}
                className="bp3-dark">
                <div className={Classes.DIALOG_BODY} >
                    <AgGridReact className="ag-theme-balham-dark" 
                        domLayout='autoHeight'
                        frameworkComponents={{
                            countryFlagRenderer: CountryFlagRenderer,
                        }}

                        defaultColDef={{
                            sortable: false,
                            filter: true,
                            resizable: true
                        }}
                        columnDefs={columnDefs}
                        rowData={this.state.executions}
                    />
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



    open(order: Order, width: number) {

        this.orderService.GetOrderHistory(order, (err: grpcWeb.Error, history: OrderHistory) => {

            if (err) {
                log.error("failed to get order history:" + err)
                return
            }


            let executionIds = new Set<String>()

            let newViews = new Array<ExecutionView>()
            for (let update of history.getUpdatesList()) {



                let order = update.getOrder()
                let time = update.getTime()
                if (order && time) {
                    if (order.getLastexecid() !== "" && !executionIds.has(order.getLastexecid())) {
                        executionIds.add(order.getLastexecid())
                        let listing = this.listingService.GetListingImmediate(order.getListingid())
                        let view = new ExecutionView(order, time, listing)
                        newViews.push(view)
                    }
                }
            }


            let blotterState = {
                ...this.state, ...{
                    parentOrder: order,
                    isOpen: true,
                    usePortal: false,
                    width: width,
                    executions: newViews,
                }
            }

            this.setState(blotterState)

        })


    }


    private handleClose = () => {
        this.setState({
            ...this.state, ...{
                isOpen: false,
            }
        })
    };

}
