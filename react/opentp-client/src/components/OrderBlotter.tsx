import { Menu, Colors } from '@blueprintjs/core';
import { Cell, Column, IMenuContext, IRegion, SelectionModes, Table } from "@blueprintjs/table";
import "@blueprintjs/table/lib/css/table.css";
import * as grpcWeb from 'grpc-web';
import React from 'react';
import { MenuItem } from "react-contextmenu";
import v4 from 'uuid';
import { logDebug, logGrpcError } from '../logging/Logging';
import { Order, OrderStatus, Side } from '../serverapi/order_pb';
import { ViewServiceClient } from '../serverapi/View-serviceServiceClientPb';
import { SubscribeToOrders } from '../serverapi/view-service_pb';
import { toNumber } from '../util/decimal64Conversion';
import { OrderContext } from './Container';
import Login from './Login';
import './OrderBlotter.css';
import { ExecutionVenueClient } from '../serverapi/Execution-venueServiceClientPb';
import { OrderId } from '../serverapi/execution-venue_pb';
import { Empty } from '../serverapi/common_pb';

interface OrderBlotterState {
  orders: OrderView[];
  selectedOrders: Map<string, Order>
}

interface OrderBlotterProps {
  orderContext: OrderContext
}


class OrderView {

  version: number;
  id: string;
  side: string;
  quantity?: number;
  price?: number;
  listingId: number;
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

export default class OrderBlotter extends React.Component<OrderBlotterProps, OrderBlotterState> {

  viewService = new ViewServiceClient(Login.grpcContext.serviceUrl, null, null)
  executionVenueService = new ExecutionVenueClient(Login.grpcContext.serviceUrl, null, null)

  orderMap: Map<string, OrderView>;

  stream?: grpcWeb.ClientReadableStream<Order>;

  id: string;

  constructor(props: OrderBlotterProps) {
    super(props);

    this.id = v4();

    this.orderMap = new Map<string, OrderView>();

    let blotterState: OrderBlotterState = {
      orders: Array.from(this.orderMap.values()),
      selectedOrders: new Map<string, Order>()
    }

    this.state = blotterState;

    this.stream = this.viewService.subscribe(new SubscribeToOrders(), Login.grpcContext.grpcMetaData)

    this.stream.on('data', (order: Order) => {

      let orderView = new OrderView(order)

      this.orderMap.set(order.getId(), orderView);

      let blotterState: OrderBlotterState = {
        ...this.state, ...{
          orders: Array.from(this.orderMap.values()),
        }
      }


      this.setState(blotterState);
    });
    this.stream.on('status', (status: grpcWeb.Status) => {
      if (status.metadata) {
        logDebug('Received metadata:' + status.metadata);
      }
    });
    this.stream.on('error', (err: grpcWeb.Error) => {
      logGrpcError('Error subscribing to orders:', err)
    });
    this.stream.on('end', () => {
      logDebug('stream end signal received');
    });


  }


  cancelOrder = (e: any, orders: Array<Order>) => {

    orders.forEach(order => {
      let orderId = new OrderId()
      orderId.setOrderid(order.getId())
      this.executionVenueService.cancelOrder(orderId, Login.grpcContext.grpcMetaData, (err: grpcWeb.Error, response: Empty) => {
        if (err) {
          logGrpcError("error cancelling order", err)
        }
      })
    });

  }

  modifyOrder = (e: any, data: Order) => {
    if (data) {
      window.alert("modify order" + data.getId());
    }
  }


  public render() {

    return (
      <div>

        <Table enableRowResizing={false} numRows={this.orderMap.size} className="bp3-dark" selectionModes={SelectionModes.ROWS_AND_CELLS}
          bodyContextMenuRenderer={this.renderBodyContextMenu} onSelection={this.onSelection} >
          <Column name="Id" cellRenderer={this.renderId} />
          <Column name="Side" cellRenderer={this.renderSide} />
          <Column name="Listing Id" cellRenderer={this.renderListingId} />
          <Column name="Quantity" cellRenderer={this.renderQuantity} />
          <Column name="Price" cellRenderer={this.renderPrice} />
          <Column name="Status" cellRenderer={this.renderStatus} />
          <Column name="Target Status" cellRenderer={this.renderTargetStatus} />
          <Column name="Rem Qty" cellRenderer={this.renderRemQty} />
          <Column name="Traded Qty" cellRenderer={this.renderTrdQty} />
          <Column name="Avg Price" cellRenderer={this.renderAvgPrice} />

        </Table>

      </div>


    );
  }

  private renderId = (row: number) => <Cell>{Array.from(this.orderMap.values())[row].id}</Cell>;
  private renderSide = (row: number) => <Cell>{Array.from(this.orderMap.values())[row].side}</Cell>;
  private renderQuantity = (row: number) => <Cell>{Array.from(this.orderMap.values())[row].quantity}</Cell>;
  private renderPrice = (row: number) => <Cell>{Array.from(this.orderMap.values())[row].price}</Cell>;
  private renderListingId = (row: number) => <Cell>{Array.from(this.orderMap.values())[row].listingId}</Cell>;
  private renderRemQty = (row: number) => <Cell>{Array.from(this.orderMap.values())[row].remainingQuantity}</Cell>;
  private renderTrdQty = (row: number) => <Cell>{Array.from(this.orderMap.values())[row].tradedQuantity}</Cell>;
  private renderAvgPrice = (row: number) => <Cell>{Array.from(this.orderMap.values())[row].avgTradePrice}</Cell>;

  private renderStatus = (row: number) => {
    let orderView = Array.from(this.orderMap.values())[row]
    let statusStyle = {}
    switch (orderView.order.getStatus()) {
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

    return <Cell style={statusStyle}>{orderView.status}</Cell>
  }


  private renderTargetStatus = (row: number) => {
    let orderView = Array.from(this.orderMap.values())[row]
    let statusStyle = {}
    if (orderView.order.getTargetstatus() !== OrderStatus.NONE) {
      statusStyle = { background: Colors.ORANGE3 }
    }


    return <Cell style={statusStyle}>{orderView.targetStatus}</Cell>
  }


  private onSelection = (selectedRegions: IRegion[]) => {
    let newSelectedOrders: Map<string, Order> = this.getSelectedOrdersFromRegions(selectedRegions);

    let blotterState: OrderBlotterState = {
      ...this.state, ...{
        selectedOrders: newSelectedOrders,
      }
    }

    this.setState(blotterState)
  }


  private getSelectedOrdersFromRegions(selectedRegions: IRegion[]): Map<string, Order> {
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
        lastRowIdx = this.state.orders.length - 1;
      }
      for (let i = firstRowIdx; i <= lastRowIdx; i++) {
        let order = this.state.orders[i].order;
        newSelectedOrders.set(order.getId(), order);
      }
    }
    return newSelectedOrders;
  }


  cancelleableOrders(orders: Map<string, Order>): Array<Order> {

    let result = new Array<Order>()
    for (let order of orders.values()) {
      if (order.getStatus() === OrderStatus.LIVE) {
        result.push(order)
      }
    }

    return result
  }


  private renderBodyContextMenu = (context: IMenuContext) => {

    let selectedOrders = this.getSelectedOrdersFromRegions(context.getRegions())
    let cancelleableOrders = this.cancelleableOrders(selectedOrders)

    return (
      <Menu>
        <MenuItem data={this.props.orderContext.selectedOrder} onClick={e => this.cancelOrder(e, cancelleableOrders)} disabled={cancelleableOrders.length === 0} >
          Cancel Order
              </MenuItem>
        <MenuItem divider />
        <MenuItem data={this.props.orderContext.selectedOrder} onClick={this.modifyOrder}>
          Modify Order
              </MenuItem>
      </Menu>
    );
  };

}
