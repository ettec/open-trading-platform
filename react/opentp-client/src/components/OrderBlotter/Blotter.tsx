import { Colors } from '@blueprintjs/core';
import { Cell, Column, IRegion } from "@blueprintjs/table";
import "@blueprintjs/table/lib/css/table.css";
import React from 'react';
import { Order, OrderStatus } from '../../serverapi/order_pb';
import '../TableView/TableCommon.css';
import { reorderColumnData } from '../TableView/TableLayout';
import '../TableView/TableLayout.ts';
import { OrderView } from './OrderView';


export interface BlotterState {

  orders: OrderView[];
  columns: Array<JSX.Element>
  columnWidths: Array<number>
}


export default class Blotter<P,S extends BlotterState >  extends React.Component<P, S>{






  getColumns() {
    return [<Column key="id" id="id" name="Id" cellRenderer={this.renderId} />,
    <Column key="side" id="side" name="Side" cellRenderer={this.renderSide} />,
    <Column key="symbol" id="symbol" name="Symbol" cellRenderer={this.renderSymbol} />,
    <Column key="mic" id="mic" name="Mic" cellRenderer={this.renderMic} />,
    <Column key="country" id="country" name="Country" cellRenderer={this.renderCountry} />,
    <Column key="quantity" id="quantity" name="Quantity" cellRenderer={this.renderQuantity} />,
    <Column key="price" id="price" name="Price" cellRenderer={this.renderPrice} />,
    <Column key="status" id="status" name="Status" cellRenderer={this.renderStatus} />,
    <Column key="targetStatus" id="targetStatus" name="Target Status" cellRenderer={this.renderTargetStatus} />,
    <Column key="remQty" id="remQty" name="Rem. Qty" cellRenderer={this.renderRemQty} />,
    <Column key="exposedQty" id="exposedQty" name="Exp. Qty" cellRenderer={this.renderExpQty} />,
    <Column key="tradedQty" id="tradedQty" name="Traded Qty" cellRenderer={this.renderTrdQty} />,
    <Column key="avgPrice" id="avgPrice" name="Avg Price" cellRenderer={this.renderAvgPrice} />,
    <Column key="listingId" id="listingId" name="Listing Id" cellRenderer={this.renderListingId} />,
    <Column key="created" id="created" name="Created" cellRenderer={this.renderCreated} />,
    <Column key="placedWith" id="placedWith" name="Placed With" cellRenderer={this.renderPlacedWith} />
    ];
  }

  private renderId = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.id}</Cell>;
  private renderSide = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.side}</Cell>;
  private renderQuantity = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.quantity}</Cell>;
  private renderSymbol = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.getSymbol()}</Cell>;
  private renderMic = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.getMic()}</Cell>;
  private renderCountry = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.getCountryCode()}</Cell>;
  private renderPrice = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.price}</Cell>;
  private renderListingId = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.listingId}</Cell>;
  private renderRemQty = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.remainingQuantity}</Cell>;
  private renderExpQty = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.exposedQuantity}</Cell>;
  private renderTrdQty = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.tradedQuantity}</Cell>;
  private renderAvgPrice = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.avgTradePrice}</Cell>;
  private renderPlacedWith = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.placedWith}</Cell>;

  private renderCreated = (row: number) => {
    let created = Array.from(this.state.orders)[row]?.created

    if (created) {
      return <Cell>{created.toLocaleTimeString()}</Cell>
    } else {
      return <Cell></Cell>
    }
  }

  private renderStatus = (row: number) => {
    let orderView = Array.from(this.state.orders)[row]
    let statusStyle = {}
    if (orderView) {
      switch (orderView.getOrder().getStatus()) {
        case OrderStatus.LIVE:
          statusStyle = { background: Colors.GREEN3 }
          break
        case OrderStatus.CANCELLED:
          statusStyle = { background: Colors.RED3 }
          break
        case OrderStatus.FILLED:
          statusStyle = { background: Colors.BLUE3 }
          break
      }
    }



    return <Cell style={statusStyle}>{orderView?.status}</Cell>
  }


  private renderTargetStatus = (row: number) => {
    let orderView = Array.from(this.state.orders)[row]
    let statusStyle = {}
    if (orderView) {
      if (orderView.getOrder().getTargetstatus() !== OrderStatus.NONE) {
        statusStyle = { background: Colors.ORANGE3 }
      }
    }


    return <Cell style={statusStyle}>{orderView?.targetStatus}</Cell>
  }



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