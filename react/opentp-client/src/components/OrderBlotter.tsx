import { Colors, Menu } from '@blueprintjs/core';
import { Cell, Column, IMenuContext, IRegion, SelectionModes, Table } from "@blueprintjs/table";
import "@blueprintjs/table/lib/css/table.css";
import * as grpcWeb from 'grpc-web';
import React from 'react';
import { MenuItem } from "react-contextmenu";
import v4 from 'uuid';
import { logDebug, logGrpcError } from '../logging/Logging';
import { Empty } from '../serverapi/common_pb';
import { ExecutionVenueClient } from '../serverapi/Execution-venueServiceClientPb';
import { OrderId } from '../serverapi/execution-venue_pb';
import { Order, OrderStatus, Side } from '../serverapi/order_pb';
import { ViewServiceClient } from '../serverapi/View-serviceServiceClientPb';
import { SubscribeToOrders } from '../serverapi/view-service_pb';
import { ListingService } from '../services/ListingService';
import { toNumber } from '../util/decimal64Conversion';
import { OrderContext } from './Container';
import Login from './Login';
import './OrderBlotter.css';
import { Listing } from '../serverapi/listing_pb';

interface OrderBlotterState {
  orders: OrderView[];
  selectedOrders: Map<string, Order>
}

interface OrderBlotterProps {
  orderContext: OrderContext
  listingService: ListingService
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
  private order: Order;
  listing?: Listing;

  constructor(order: Order) {
    this.id = ""
    this.version = 0;
    this.side = "";
    this.listingId = 0;
    this.status = "";
    this.targetStatus = "";
    this.order = order

    this.setOrder(order)
  }

  setOrder(order: Order) {
    this.order = order


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

  getOrder(): Order {
    return this.order
  }

  getStatusString(status: OrderStatus) {

    switch (status) {
      case OrderStatus.CANCELLED:
        return "Cancelled"
      case OrderStatus.FILLED:
        return "Filled"
      case OrderStatus.LIVE:
        return "Live"
      case OrderStatus.NONE:
        return "None"
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

export default class OrderBlotter extends React.Component<OrderBlotterProps, OrderBlotterState> {

  viewService = new ViewServiceClient(Login.grpcContext.serviceUrl, null, null)
  executionVenueService = new ExecutionVenueClient(Login.grpcContext.serviceUrl, null, null)
  listingService: ListingService

  orderMap: Map<string, number>;

  stream?: grpcWeb.ClientReadableStream<Order>;

  id: string;


  constructor(props: OrderBlotterProps) {
    super(props);

    this.id = v4();

    this.listingService = props.listingService

    this.orderMap = new Map<string, number>();

    let view = new Array<OrderView>(5)

    let blotterState: OrderBlotterState = {
      orders: view,
      selectedOrders: new Map<string, Order>()
    }

    this.state = blotterState;

    this.stream = this.viewService.subscribe(new SubscribeToOrders(), Login.grpcContext.grpcMetaData)

    this.stream.on('data', (order: Order) => {

      let idx = this.orderMap.get(order.getId())
      let orderView: OrderView
      let newOrders = this.state.orders
      if (!idx) {
        idx = this.orderMap.size
        let orderLength = this.state.orders.length
        if (idx >= orderLength) {
          newOrders = new Array<OrderView>(orderLength * 2)
          for( let i=0; i < orderLength; i++) {
              newOrders[i] = this.state.orders[i]
          }
        }

        this.orderMap.set(order.getId(), idx);
        orderView = new OrderView(order)

        newOrders[idx] = orderView
      } else {
        orderView = this.state.orders[idx]
        orderView.setOrder(order)
      }

      let blotterState: OrderBlotterState = {
        ...this.state, ...{
          orders: newOrders
        }
      }

      this.setState(blotterState);

      if (!orderView.listing) {
        this.listingService.GetListing(order.getListingid(), (listing: Listing) => {
          orderView.listing = listing
          let blotterState: OrderBlotterState = {
            ...this.state, ...{
              orders: this.state.orders
            }
          }

          this.setState(blotterState);
        })
      }



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

        <Table enableRowResizing={false} numRows={this.state.orders.length} className="bp3-dark" selectionModes={SelectionModes.ROWS_AND_CELLS}
          bodyContextMenuRenderer={this.renderBodyContextMenu} onSelection={this.onSelection} >
          <Column name="Id" cellRenderer={this.renderId} />
          <Column name="Side" cellRenderer={this.renderSide} />
          <Column name="Symbol" cellRenderer={this.renderSymbol} />
          <Column name="Mic" cellRenderer={this.renderMic} />
          <Column name="Country" cellRenderer={this.renderCountry} />
          <Column name="Quantity" cellRenderer={this.renderQuantity} />
          <Column name="Price" cellRenderer={this.renderPrice} />
          <Column name="Status" cellRenderer={this.renderStatus} />
          <Column name="Target Status" cellRenderer={this.renderTargetStatus} />
          <Column name="Rem Qty" cellRenderer={this.renderRemQty} />
          <Column name="Traded Qty" cellRenderer={this.renderTrdQty} />
          <Column name="Avg Price" cellRenderer={this.renderAvgPrice} />
          <Column name="Listing Id" cellRenderer={this.renderListingId} />

        </Table>

      </div>


    );
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
  private renderTrdQty = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.tradedQuantity}</Cell>;
  private renderAvgPrice = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.avgTradePrice}</Cell>;

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
        let orderView = this.state.orders[i];
        if (orderView) {
          newSelectedOrders.set(orderView.getOrder().getId(), orderView.getOrder());
        }
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
