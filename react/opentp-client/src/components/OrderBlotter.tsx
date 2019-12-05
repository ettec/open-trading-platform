import React from 'react';
import './OrderBlotter.css';

import v4 from 'uuid';
import { ContextMenu, MenuItem, ContextMenuTrigger } from "react-contextmenu";
import { ViewServiceClient } from '../serverapi/View-serviceServiceClientPb';
import * as grpcWeb from 'grpc-web'
import Login from './Login';
import { SubscribeToOrders } from '../serverapi/view-service_pb';
import { Order, Side, OrderStatus } from '../serverapi/order_pb';
import { Decimal64 } from '../serverapi/common_pb';
import { number, string } from 'prop-types';
import { toNumber } from '../util/decimal64Conversion';
import { Table, Column, Cell, SelectionModes, IMenuContext, IRegion } from "@blueprintjs/table";
import "@blueprintjs/table/lib/css/table.css"
import { Menu } from '@blueprintjs/core';
import { start } from 'repl';

  

interface BlotterState {
  orders: OrderView[];
  selectedOrders : Map<string, Order>
}

interface Props {
  onOrderSelected: (order: Order) => void
  selectedOrder?: Order;
}





class OrderView {

  version: number;
  id: string;
  side: string;
  quantity?: number;
  price?: number;
  listingId: string;
  remainingQuantity?: number;
  tradedQuantity?: number;
  avgTradePrice?: number;
  status: string;
  targetStatus: string;
  order: Order;

  constructor(order: Order) {
    this.order = order;
    this.version = order.getVersion()
    this.id = order.getId()

    switch (order.getSide()) {
      case Side.BUY:
        this.side = "Buy"
        break;
      case Side.SELL:
        this.side = "Sell"
        break;
      default:
        this.side = "Unknown Side: " + order.getSide()
    }

    this.quantity = toNumber(order.getQuantity())
    this.price = toNumber(order.getPrice())
    this.listingId = order.getListingid()
    this.remainingQuantity = toNumber(order.getRemainingquantity())
    this.tradedQuantity = toNumber(order.getTradedquantity())
    this.avgTradePrice = toNumber(order.getAvgtradeprice())
    this.status = this.getStatusString(order.getStatus())
    this.targetStatus = this.getStatusString(order.getTargetstatus())


  }

  private getStatusString(status: OrderStatus): string {
    switch (status) {
      case OrderStatus.CANCELLED:
        return "Cancelled"
      case OrderStatus.FILLED:
        return "Filled"
      case OrderStatus.LIVE:
        return "Live"
      case OrderStatus.NONE:
        return "None"
      default:
        return "Unknown status: " + status
    }


  }

}

export default class OrderBlotter extends React.Component<Props, BlotterState> {

  viewService = new ViewServiceClient(Login.grpcContext.serviceUrl, null, null)

  orderMap: Map<string, OrderView>;

  stream?: grpcWeb.ClientReadableStream<Order>;


  //   ordersSource : EventSource;
  id: string;

  constructor(props: Props) {
    super(props);

    this.id = v4();

    this.orderMap = new Map<string, OrderView>();
 
    let blotterState: BlotterState = {
      orders: Array.from(this.orderMap.values()),
      selectedOrders: new Map<string, Order>()
    }

    this.state = blotterState;

    this.stream = this.viewService.subscribe(new SubscribeToOrders(), Login.grpcContext.grpcMetaData)

    this.stream.on('data', (order: Order) => {
      console.log('Received an order' + order)


      let orderView = new OrderView(order)

      this.orderMap.set(order.getId(), orderView);

      console.log("Values " + this.orderMap.values());

      let blotterState: BlotterState = {
        ...this.state, ... {
          orders: Array.from(this.orderMap.values()),
        }
      }


      this.setState(blotterState);
    });
    this.stream.on('status', (status: grpcWeb.Status) => {
      if (status.metadata) {
        console.log('Received metadata');
        console.log(status.metadata);
      }
    });
    this.stream.on('error', (err: grpcWeb.Error) => {
      console.log('Received error:' + err)
    });
    this.stream.on('end', () => {
      console.log('stream end signal received');
    });


  }


  cancelOrder = (e: any, data: Order) => {
    if (data) {
      http://192.168.1.102:32413/order-management/cancel-order?orderId=00a2fdb5-7521-4f44-a985-2cddc7a19222

      fetch('http://192.168.1.102:32413/order-management/cancel-order?orderId=' + data.getId(), {
        method: 'POST',
        mode: 'no-cors'
      })
        .then(
          response => { console.log(response.statusText) }
        )

        .catch(error => {
          throw new Error(error);
        });
    }
  }

  modifyOrder = (e: any, data: Order) => {
    if (data) {
      window.alert("modify order" + data.getId());
    }
  }


  public render() {


    const myClonedArray = Object.assign([], this.state.orders);

    return (
      <div>
        
        <Table enableRowResizing={false} numRows={this.orderMap.size} className="bp3-dark" selectionModes={SelectionModes.ROWS_AND_CELLS}
        bodyContextMenuRenderer={this.renderBodyContextMenu} onSelection={this.onSelection}>
                    <Column name="Id" cellRenderer={this.renderId} />            
                    <Column name="Side" cellRenderer={this.renderSide} />
                    <Column name="Quantity" cellRenderer={this.renderQuantity} />
                    <Column name="Price" cellRenderer={this.renderPrice} />
        </Table>

      </div>


    );
  }

  private renderId = (row: number) => <Cell>{Array.from(this.orderMap.values())[row].id}</Cell>;
  private renderSide = (row: number) => <Cell>{Array.from(this.orderMap.values())[row].side}</Cell>;
  private renderQuantity = (row: number) => <Cell>{Array.from(this.orderMap.values())[row].quantity}</Cell>;
  private renderPrice = (row: number) => <Cell>{Array.from(this.orderMap.values())[row].price}</Cell>;


  private onSelection = (selectedRegions: IRegion[]) => {
    let selectedOrders : Map<string, Order> = new Map<string,Order>()


    for( let region of selectedRegions)  {

      let firstRowIdx : number;
      let lastRowIdx : number;


      if( region.rows ) {
        firstRowIdx = region.rows[0]
        lastRowIdx = region.rows[1]
      }  else {
        firstRowIdx = 0
        lastRowIdx = this.state.orders.length -1
      }

      for (let i = firstRowIdx; i < lastRowIdx; i++) {
        let order = this.state.orders[i].order
        selectedOrders.set(order.getId(), order)
      }
    }


    
  }

  private renderBodyContextMenu = (context: IMenuContext) => {
    return (
        <Menu>
             <MenuItem data={this.props.selectedOrder} onClick={this.cancelOrder} disabled={this.state.selectedOrders.size==0} >
            Cancel Order
              </MenuItem>
          <MenuItem divider />
          <MenuItem data={this.props.selectedOrder} onClick={this.modifyOrder}>
            Modify Order
              </MenuItem>
        </Menu>
    );
};

}
