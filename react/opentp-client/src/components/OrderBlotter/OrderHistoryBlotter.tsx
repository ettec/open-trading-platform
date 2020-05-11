import { AnchorButton, Classes, Dialog, Intent } from '@blueprintjs/core';
import { Cell, Column, SelectionModes, Table } from "@blueprintjs/table";
import * as grpcWeb from 'grpc-web';
import * as React from "react";
import { Timestamp } from '../../serverapi/modelcommon_pb';
import { Order } from '../../serverapi/order_pb';
import { OrderHistory } from '../../serverapi/viewservice_pb';
import { ListingService } from '../../services/ListingService';
import { OrderService } from "../../services/OrderService";
import { OrderHistoryBlotterController } from '../Container';
import  { getConfiguredColumns, TableViewConfig } from "../TableView/TableView";
import OrderBlotter from "./OrderBlotter";
import { OrderView } from "./OrderView";

export interface OrderHistoryBlotterProps {
    orderService: OrderService
    listingService: ListingService
    orderHistoryBlotterController: OrderHistoryBlotterController
}

interface OrderHistoryBlotterState {
    orders: OrderView[];
    isOpen: boolean,
    usePortal: boolean
    order?: Order
    columns: Array<JSX.Element>
    columnWidths: Array<number>
    updates: OrderUpdateView[];
    width: number
}


export default class OrderHistoryBlotter extends OrderBlotter<OrderHistoryBlotterProps, OrderHistoryBlotterState> {

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
            columns: new Array<JSX.Element>(),
            columnWidths: new Array<number>(),
            updates: new Array<OrderUpdateView>(0),
            orders: new Array<OrderView>(),
            width: 0
        }

    }

    render() {
        return (
            <Dialog
                icon="bring-data"
                onClose={this.handleClose}
                title={this.state.order?.getId()}
                style={{ minWidth: this.state.width }}
                {...this.state}
                className="bp3-dark">
                <div className={Classes.DIALOG_BODY} >
                    <Table enableRowResizing={false} numRows={this.state.updates.length} className="bp3-dark" selectionModes={SelectionModes.ROWS_AND_CELLS}
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







    open(order: Order, config: TableViewConfig, width: number) {



        this.orderService.GetOrderHistory(order, (err: grpcWeb.Error, history: OrderHistory) => {

            let newViews = new Array<OrderUpdateView>()
            for (let update of history.getUpdatesList()) {

                let order = update.getOrder()
                let time = update.getTime()
                if (order && time) {
                    let view = new OrderUpdateView(order, time)
                    let listing = this.listingService.GetListingImmediate(order.getListingid())
                    if( listing ) {
                        view.listing = listing
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


        let [cols, widths] = getConfiguredColumns(this.getColumns(), config);


        let newWidths = new Array<number>()
        newWidths.push(75)
        newWidths = newWidths.concat(widths)

        let timeCol = <Column key="time" id="time" name="Time" cellRenderer={this.renderUpdateTime} />

        let newCols = new Array<JSX.Element>()
        newCols.push(timeCol)
        newCols = newCols.concat(cols)

        let state =
        {
            order: order,
            isOpen: true,
            usePortal: false,
            columns: newCols,
            columnWidths: newWidths,
            updates: new Array<OrderUpdateView>(),
            orders: new Array<OrderUpdateView>(),
            width: width

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

    private renderUpdateTime = (row: number) => {
        let updateTime = Array.from(this.state.updates)[row]?.time

        if (updateTime) {
            return <Cell>{updateTime.toLocaleTimeString()}</Cell>
        } else {
            return <Cell></Cell>
        }
    }

}

class OrderUpdateView extends OrderView {

    time?: Date;

    constructor(order: Order, updateTime: Timestamp) {
        super(order)

        if (updateTime) {
            this.time = new Date(updateTime.getSeconds() * 1000)
        }
    }

}
