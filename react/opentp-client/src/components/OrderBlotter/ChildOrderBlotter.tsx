import { AnchorButton, Classes, Dialog, Intent } from '@blueprintjs/core';
import { IRegion, SelectionModes, Table } from "@blueprintjs/table";
import * as React from "react";
import { Order, OrderStatus } from '../../serverapi/order_pb';
import { ListingService } from '../../services/ListingService';
import { OrderService } from "../../services/OrderService";
import { ChildOrderBlotterController } from '../Container/Container';
import { getConfiguredColumns, TableViewConfig, TableViewProperties } from '../TableView/TableView';
import OrderBlotter, { OrderBlotterState } from "./OrderBlotter";
import { OrderView } from "./OrderView";

export interface ChildOrderProps extends TableViewProperties   {
    orderService: OrderService
    listingService: ListingService
    childOrderBlotterController: ChildOrderBlotterController
}

  
interface ChildOrderBlotterState extends OrderBlotterState {
    isOpen: boolean,
    usePortal: boolean
    parentOrder?: Order
    selectedOrders: Array<Order>,
    width: number,
}


export default class ChildOrderBlotter  extends OrderBlotter<ChildOrderProps, ChildOrderBlotterState> {

    orderService: OrderService
    listingService: ListingService
    childOrderBlotterController: ChildOrderBlotterController

    constructor(props: ChildOrderProps) {
        super(props)

        this.orderService = props.orderService
        this.listingService = props.listingService
        this.childOrderBlotterController = props.childOrderBlotterController

        this.childOrderBlotterController.setBlotter(this)
        let visibleStates = new Set<OrderStatus>() 
        visibleStates.add(OrderStatus.LIVE)
        visibleStates.add(OrderStatus.CANCELLED)
        visibleStates.add(OrderStatus.FILLED)
        visibleStates.add(OrderStatus.NONE)

        this.state = {
            isOpen: false,
            usePortal: false,
            columns: new Array<JSX.Element>(),
            columnWidths: new Array<number>(),
            orders: new Array<OrderView>(10),
            selectedOrders: new Array<Order>(),
            width: 0,
            visibleStates: visibleStates
        }

    }

    protected  getTableName() : string {
        return "Child Orders"
    }


    getTitle(order? : Order ) : string {
        if( order ) {
            return "Children of order " + order.getId()
        }

        return ""
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
                    <Table enableRowResizing={false} numRows={this.state.orders.length} className="bp3-dark" selectionModes={SelectionModes.ROWS_AND_CELLS}
                        onSelection={this.onSelection} enableColumnReordering={true}
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



    private onSelection = (selectedRegions: IRegion[]) => {
        let newSelectedOrders: Array<Order> = this.getSelectedOrdersFromRegions(selectedRegions, this.state.orders);

        let blotterState: ChildOrderBlotterState = {
            ...this.state, ...{
                selectedOrders: newSelectedOrders,
            }
        }

        this.setState(blotterState)
    }

   

    open(parentOrder: Order, orders: Array<Order>,  config: TableViewConfig, width: number) {


        let [cols, widths] = getConfiguredColumns(this.getColumns(), config);

        this.clearOrders()

        for (let order of orders) {
            this.addOrUpdateOrder(order)
        }

        let state =
        {
            parentOrder: parentOrder,
            isOpen: true,
            usePortal: false,
            columns: cols,
            columnWidths: widths,
            selectedOrders: new Array<Order>(),
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
