import { Layout } from "flexlayout-react";
import { Listing } from "../../serverapi/listing_pb";
import { Order, Side } from "../../serverapi/order_pb";
import Executions from "../Executions";
import ChildOrderBlotter from "../OrderBlotter/ChildOrderBlotter";
import OrderHistoryBlotter from "../OrderBlotter/OrderHistoryBlotter";
import OrderTicket from '../OrderTicket/OrderTicket';
import QuestionDialog from "./QuestionDialog";
import ColumnChooser from "../TableView/ColumnChooser";
import { TableViewConfig } from "../TableView/TableView";
import ViewNameDialog from "./ViewNameDialog";
import ColumnChooserAgGrid from "../TableView/ColumnChooseAgGrid";
import { Column, ColumnState } from "ag-grid-community";


export class AgGridColumnChooserController {

    private dialog?: ColumnChooserAgGrid;

    setDialog(dialog: ColumnChooserAgGrid) {
        this.dialog = dialog;
    }

    open(tableName: string, colStates: ColumnState[], cols: Column[], callback: (columns: ColumnState[] | undefined) => void) {
        if (this.dialog) {
            this.dialog.open(tableName, colStates, cols, callback);
        }
    }

}


export class ColumnChooserController {

    private dialog?: ColumnChooser;

    setDialog(dialog: ColumnChooser) {
        this.dialog = dialog;
    }

    open(tableName: string, visibleColumns: JSX.Element[], widths: number[], allColumns: JSX.Element[], callback: (newVisibleCols: JSX.Element[] | undefined,
        widths: number[] | undefined) => void) {
        if (this.dialog) {
            this.dialog.open(tableName, visibleColumns, widths, allColumns, callback);
        }
    }

}

export class QuestionDialogController {

    private dialog?: QuestionDialog;

    setDialog(dialog: QuestionDialog) {
        this.dialog = dialog;
    }

    open(question: string, title: string, callback: (response: boolean) => void) {
        if (this.dialog) {
            this.dialog.open(question, title, callback);
        }
    }

}

export class ViewNameDialogController {

    private dialog?: ViewNameDialog;

    setDialog(dialog: ViewNameDialog) {
        this.dialog = dialog;
    }

    open(component: string, componentDislayName: string, layout: Layout) {
        if (this.dialog) {
            this.dialog.open(component, componentDislayName, layout);
        }
    }

}



export class ExecutionsController {

    private executions?: Executions;

    setView(executions: Executions) {
        this.executions = executions;
    }

    open(order: Order, width: number) {
        if (this.executions) {
            this.executions.open(order, width);
        }
    }

}

export class OrderHistoryBlotterController {

    private orderHistoryBlotter?: OrderHistoryBlotter;

    setBlotter(orderHistoryBlotter: OrderHistoryBlotter) {
        this.orderHistoryBlotter = orderHistoryBlotter;
    }

    openBlotter(order: Order, config: TableViewConfig, width: number) {
        if (this.orderHistoryBlotter) {
            this.orderHistoryBlotter.open(order, config, width);
        }
    }

}


export class ChildOrderBlotterController {

    private childOrderBlotter?: ChildOrderBlotter;

    setBlotter(childOrderBlotter: ChildOrderBlotter) {
        this.childOrderBlotter = childOrderBlotter;
    }

    openBlotter(parentOrder: Order, orders: Array<Order>,
        config: TableViewConfig, width: number) {
        if (this.childOrderBlotter) {
            this.childOrderBlotter.open(parentOrder, orders, config, width);
        }
    }

}



export class TicketController {

    private orderTicket?: OrderTicket;

    setOrderTicket(orderTicket: OrderTicket) {
        this.orderTicket = orderTicket;
    }

    openNewOrderTicket(side: Side, listing: Listing) {
        if (this.orderTicket) {
            this.orderTicket.openNewOrderTicket(side, listing);
        }
    }

    openOrderTicketWithDefaultPriceAndQty(newSide: Side, newListing: Listing, defaultPrice?: number, defaultQuantity?: number) {
        if (this.orderTicket) {
            this.orderTicket.openOrderTicketWithDefaultPriceAndQty(newSide, newListing, defaultPrice, defaultQuantity);
        }
    }

    openModifyOrderTicket(order: Order, listing: Listing) {
        if (this.orderTicket) {
            this.orderTicket.openModifyOrderTicket(order, listing);
        }
    }

}
