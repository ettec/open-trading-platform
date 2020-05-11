import { IRegion } from "@blueprintjs/table";
import "@blueprintjs/table/lib/css/table.css";
import * as React from "react";
import { Order, OrderStatus } from '../../serverapi/order_pb';
import '../TableView/TableCommon.css';
import '../TableView/TableLayout.ts';
import { OrderView } from './OrderView';
import { reorderColumnData } from "../TableView/TableLayout";



export interface BlotterState {
  
  columns: Array<JSX.Element>
  columnWidths: Array<number>
}


export default class Blotter<P,S extends BlotterState >  extends React.Component<P, S>{



  columnResized = (index: number, size: number) => {
    let newColWidths = this.state.columnWidths.slice();
    newColWidths[index] = size
    let blotterState = {
        ...this.state, ...{
            columnWidths: newColWidths
        }
    }

    this.setState(blotterState)

}

onColumnsReordered = (oldIndex: number, newIndex: number, length: number) => {

    let newCols = reorderColumnData(oldIndex, newIndex, length, this.state.columns)
    let newColWidths = reorderColumnData(oldIndex, newIndex, length, this.state.columnWidths)

    let blotterState = {
        ...this.state, ...{
            columns: newCols,
            columnWidths: newColWidths
        }
    }

    this.setState(blotterState)
}


static getSelectedOrdersFromRegions(selectedRegions: IRegion[], orders: OrderView[]): Array<Order> {
    let newSelectedOrders: Array< Order> = new Array<Order>();

    let selectedOrderArray: Array<OrderView> = Blotter.getSelectedItems(selectedRegions, orders);

    for( let orderView of selectedOrderArray ) {
      newSelectedOrders.push(orderView.getOrder())
    }


    return newSelectedOrders;
  }


  private static getSelectedItems<T>(selectedRegions: IRegion[], items: T[]) {
    let selectedOrderArray: Array<T> = new Array<T>();
    for (let region of selectedRegions) {
      let firstRowIdx: number;
      let lastRowIdx: number;
      if (region.rows) {
        firstRowIdx = region.rows[0];
        lastRowIdx = region.rows[1];
      }
      else {
        firstRowIdx = 0;
        lastRowIdx = items.length - 1;
      }
      for (let i = firstRowIdx; i <= lastRowIdx; i++) {
        let item = items[i];
        if (item) {
          selectedOrderArray.push(item);
        }
      }
    }
    return selectedOrderArray;
  }

  static cancelleableOrders(orders: Array< Order>): Array<Order> {

    let result = new Array<Order>()
    for (let order of orders) {
      if (order.getStatus() === OrderStatus.LIVE) {
        result.push(order)
      }
    }

    return result
  }

}