import { IRegion } from "@blueprintjs/table";
import "@blueprintjs/table/lib/css/table.css";
import { Order, OrderStatus } from '../../serverapi/order_pb';
import '../TableView/TableCommon.css';
import '../TableView/TableLayout.ts';
import { OrderView } from './OrderView';




export default class Blotter {

static getSelectedOrdersFromRegions(selectedRegions: IRegion[], orders: OrderView[]): Map<string, Order> {
    let newSelectedOrders: Map<string, Order> = new Map<string, Order>();
    for (let region of selectedRegions) {
      let firstRowIdx: number;
      let lastRowIdx: number;
      if (region.rows) {
        firstRowIdx = region.rows[0];
        lastRowIdx = region.rows[1];
      }
      else {
        firstRowIdx = 0;
        lastRowIdx = orders.length - 1;
      }
      for (let i = firstRowIdx; i <= lastRowIdx; i++) {
        let orderView = orders[i];
        if (orderView) {
          newSelectedOrders.set(orderView.getOrder().getId(), orderView.getOrder());
        }
      }
    }
    return newSelectedOrders;
  }


  static cancelleableOrders(orders: Map<string, Order>): Array<Order> {

    let result = new Array<Order>()
    for (let order of orders.values()) {
      if (order.getStatus() === OrderStatus.LIVE) {
        result.push(order)
      }
    }

    return result
  }

}