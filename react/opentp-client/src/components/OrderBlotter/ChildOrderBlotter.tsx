import { AnchorButton, Classes, Dialog, Intent, Menu } from '@blueprintjs/core';
import { IMenuContext, IRegion, SelectionModes, Table } from "@blueprintjs/table";
import * as React from "react";
import { Order } from '../../serverapi/order_pb';
import { OrderService } from "../../services/OrderService";
import { ChildOrderBlotterController } from '../Container';
import TableViewConfig, { getConfiguredColumns } from "../TableView/TableLayout";
import Blotter from "./Blotter";
import { OrderView } from "./OrderView";

export interface ChildOrderProps {
    orderService: OrderService
    childOrderBlotterController: ChildOrderBlotterController
}


interface ChildOrderBlotterState {
    isOpen: boolean,
    usePortal: boolean
    parentOrder?: Order
    columns: Array<JSX.Element>
    columnWidths: Array<number>
    orders: OrderView[];
    selectedOrders: Map<string, Order>,
    width: number
}


export default class ChildOrderBlotter extends Blotter<ChildOrderProps, ChildOrderBlotterState> {

    orderService: OrderService
    childOrderBlotterController: ChildOrderBlotterController

    constructor(props: ChildOrderProps) {
        super(props)

        this.orderService = props.orderService
        this.childOrderBlotterController = props.childOrderBlotterController

        this.childOrderBlotterController.setBlotter(this)

        this.state = {
            isOpen: false,
            usePortal: false,
            columns: new Array<JSX.Element>(),
            columnWidths: new Array<number>(),
            orders: new Array<OrderView>(10),
            selectedOrders: new Map<string, Order>(),
            width: 0
        }

    }

    render() {
        return (
            <Dialog
                icon="bring-data"
                onClose={this.handleClose}
                title={this.state.parentOrder?.getId()}
                style={{ minWidth: this.state.width }}
                {...this.state}
                className="bp3-dark">
                <div className={Classes.DIALOG_BODY} >
                    <Table enableRowResizing={false} numRows={this.state.orders.length} className="bp3-dark" selectionModes={SelectionModes.ROWS_AND_CELLS}
                        bodyContextMenuRenderer={this.renderBodyContextMenu} onSelection={this.onSelection} enableColumnReordering={true}
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

    private renderBodyContextMenu = (context: IMenuContext) => {

        // let selectedOrders = Blotter.getSelectedOrdersFromRegions(context.getRegions(), this.state.orders)

        return (



            <Menu  >
                <Menu.Item text="Show History">
                </Menu.Item>

            </Menu>
        );
    };


    private onSelection = (selectedRegions: IRegion[]) => {
        let newSelectedOrders: Map<string, Order> = Blotter.getSelectedOrdersFromRegions(selectedRegions, this.state.orders);

        let blotterState: ChildOrderBlotterState = {
            ...this.state, ...{
                selectedOrders: newSelectedOrders,
            }
        }

        this.setState(blotterState)
    }

   


    open(parentOrder: Order, orders: Array<Order>, columns: Array<JSX.Element>, config: TableViewConfig, width: number) {


        let [cols, widths] = getConfiguredColumns(columns, config);

        let ordersView = new Array<OrderView>()

        for (let order of orders) {
            ordersView.push(new OrderView(order))
        }

        let state =
        {
            parentOrder: parentOrder,
            isOpen: true,
            usePortal: false,
            columns: cols,
            columnWidths: widths,
            orders: ordersView,
            selectedOrders: new Map<string, Order>(),
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

}
