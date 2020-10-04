import { AnchorButton, Classes, Dialog, Intent } from '@blueprintjs/core';
import { Cell, Column, SelectionModes, Table } from "@blueprintjs/table";
import * as grpcWeb from 'grpc-web';
import * as React from "react";
import log from 'loglevel';
import { Listing } from '../serverapi/listing_pb';
import { Timestamp } from '../serverapi/modelcommon_pb';
import { Order } from '../serverapi/order_pb';
import { OrderHistory } from '../serverapi/orderdataservice_pb';
import { ListingService } from '../services/ListingService';
import { OrderService } from "../services/OrderService";
import { toNumber } from '../common/decimal64Conversion';
import { ExecutionsController } from "./Container/Controllers";
import TableView, { getConfiguredColumns, TableViewProperties } from './TableView/TableView';



export interface ExecutionsProps extends TableViewProperties {
    orderService: OrderService
    listingService: ListingService
    executionsController: ExecutionsController
}


interface ExecutionsState {
    isOpen: boolean,
    usePortal: boolean
    parentOrder?: Order
    columns: Array<JSX.Element>
    columnWidths: Array<number>
    executions: Array<ExecutionView>;
    width: number
}

class ExecutionView {
    time?: Date;
    id: string;
    quantity?: number;
    price?: number;
    listing?: Listing;

    constructor(order: Order, updateTime: Timestamp, listing?: Listing) {

        this.id = order.getLastexecid()
        this.quantity = toNumber(order.getLastexecquantity())
        this.price = toNumber(order.getLastexecprice())
        this.listing = listing

        if (updateTime) {
            this.time = new Date(updateTime.getSeconds() * 1000)
        }
    }

    getSymbol(): string | undefined {
        return this.listing?.getInstrument()?.getDisplaysymbol()
    }

    getMic(): string | undefined {
        return this.listing?.getMarket()?.getMic()
    }

    getCountryCode(): string | undefined {
        return this.listing?.getMarket()?.getCountrycode()
    }


}


export default class Executions extends TableView<ExecutionsProps, ExecutionsState> {

    orderService: OrderService
    listingService: ListingService
    executionsController: ExecutionsController

    constructor(props: ExecutionsProps) {
        super(props)

        this.orderService = props.orderService
        this.listingService = props.listingService
        this.executionsController = props.executionsController

        this.executionsController.setView(this)

        let columns = this.getColumns()
        let [defaultCols, defaultColWidths] = getConfiguredColumns(columns);

        this.state = {
            isOpen: false,
            usePortal: false,
            columns: defaultCols,
            columnWidths: defaultColWidths,
            executions: new Array<ExecutionView>(10),
            width: 0
        }

    }

    protected  getTableName() : string {
        return "Executions"
    }

    getColumns() {
        return [<Column key="time" id="time" name="Time" cellRenderer={this.renderTime} />,
        <Column key="id" id="id" name="Id" cellRenderer={this.renderId} />,
        <Column key="symbol" id="symbol" name="Symbol" cellRenderer={this.renderSymbol} />,
        <Column key="mic" id="mic" name="Mic" cellRenderer={this.renderMic} />,
        <Column key="country" id="country" name="Country" cellRenderer={this.renderCountry} />,
        <Column key="quantity" id="quantity" name="Quantity" cellRenderer={this.renderQuantity} />,
        <Column key="price" id="price" name="Price" cellRenderer={this.renderPrice} />,


        ];
    }

    private renderId = (row: number) => <Cell>{Array.from(this.state.executions)[row]?.id}</Cell>;

    private renderQuantity = (row: number) => <Cell>{Array.from(this.state.executions)[row]?.quantity}</Cell>;
    private renderSymbol = (row: number) => <Cell>{Array.from(this.state.executions)[row]?.getSymbol()}</Cell>;
    private renderMic = (row: number) => <Cell>{Array.from(this.state.executions)[row]?.getMic()}</Cell>;
    private renderCountry = (row: number) => <Cell>{Array.from(this.state.executions)[row]?.getCountryCode()}</Cell>;
    private renderPrice = (row: number) => <Cell>{Array.from(this.state.executions)[row]?.price}</Cell>;
    private renderTime = (row: number) => {
        let created = Array.from(this.state.executions)[row]?.time

        if (created) {
            return <Cell>{created.toLocaleTimeString()}</Cell>
        } else {
            return <Cell></Cell>
        }
    }


    getTitle(order? : Order ) : string {
        if( order ) {
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
                    <Table enableRowResizing={false} numRows={this.state.executions.length} className="bp3-dark" selectionModes={SelectionModes.ROWS_AND_CELLS}
                        enableColumnReordering={true}
                        onColumnsReordered={this.onColumnsReordered} enableColumnResizing={true} onColumnWidthChanged={this.columnResized} columnWidths={this.state.columnWidths}>
                        {this.state.columns}
                    </Table>



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
