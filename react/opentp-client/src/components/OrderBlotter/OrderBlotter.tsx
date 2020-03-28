import { Colors, Menu } from '@blueprintjs/core';
import { Cell, Column, IMenuContext, IRegion, SelectionModes, Table } from "@blueprintjs/table";
import "@blueprintjs/table/lib/css/table.css";
import * as grpcWeb from 'grpc-web';
import React from 'react';
import { MenuItem } from "react-contextmenu";
import v4 from 'uuid';
import { logDebug, logGrpcError } from '../../logging/Logging';
import { Empty } from '../../serverapi/common_pb';
import { OrderRouterClient } from '../../serverapi/Order-routerServiceClientPb';
import { OrderId, CancelOrderParams } from '../../serverapi/order-router_pb';
import { Order, OrderStatus } from '../../serverapi/order_pb';
import { Timestamp } from '../../serverapi/modelcommon_pb';
import { ViewServiceClient } from '../../serverapi/View-serviceServiceClientPb';
import { SubscribeToOrders } from '../../serverapi/view-service_pb';
import { ListingService } from '../../services/ListingService';

import { OrderContext } from '../Container';
import Login from '../Login';
import { Listing } from '../../serverapi/listing_pb';
import '../TableView/TableCommon.css';
import '../TableView/TableLayout.ts';
import { TabNode, Model, Actions } from 'flexlayout-react';
import { OrderView } from './OrderView';
import TableViewConfig, { getColIdsInOrder, getColumnState, reorderColumnData } from '../TableView/TableLayout';

interface OrderBlotterState {

  orders: OrderView[];
  selectedOrders: Map<string, Order>
  columns: Array<JSX.Element>
  columnWidths: Array<number>
}

interface OrderBlotterProps {
  node: TabNode,
  model: Model,
  orderContext: OrderContext
  listingService: ListingService
}



export default class OrderBlotter extends React.Component<OrderBlotterProps, OrderBlotterState> {

  viewService = new ViewServiceClient(Login.grpcContext.serviceUrl, null, null)
  executionVenueService = new OrderRouterClient(Login.grpcContext.serviceUrl, null, null)
  listingService: ListingService

  orderMap: Map<string, number>;

  stream?: grpcWeb.ClientReadableStream<Order>;

  id: string;


  constructor(props: OrderBlotterProps) {
    super(props);

    this.id = v4();

    let columns = [<Column key="id" id="id" name="Id" cellRenderer={this.renderId} />,
    <Column key="side" id="side" name="Side" cellRenderer={this.renderSide} />,
    <Column key="symbol" id="symbol" name="Symbol" cellRenderer={this.renderSymbol} />,
    <Column key="mic" id="mic" name="Mic" cellRenderer={this.renderMic} />,
    <Column key="country" id="country" name="Country" cellRenderer={this.renderCountry} />,
    <Column key="quantity" id="quantity" name="Quantity" cellRenderer={this.renderQuantity} />,
    <Column key="price" id="price" name="Price" cellRenderer={this.renderPrice} />,
    <Column key="status" id="status" name="Status" cellRenderer={this.renderStatus} />,
    <Column key="targetStatus" id="targetStatus" name="Target Status" cellRenderer={this.renderTargetStatus} />,
    <Column key="remQty" id="remQty" name="Rem Qty" cellRenderer={this.renderRemQty} />,
    <Column key="tradedQty" id="tradedQty" name="Traded Qty" cellRenderer={this.renderTrdQty} />,
    <Column key="avgPrice" id="avgPrice" name="Avg Price" cellRenderer={this.renderAvgPrice} />,
    <Column key="listingId" id="listingId" name="Listing Id" cellRenderer={this.renderListingId} />,
    <Column key="created" id="created" name="Created" cellRenderer={this.renderCreated} />,
    <Column key="placedWith" id="placedWith" name="Placed With" cellRenderer={this.renderPlacedWith} />
    ]







    let view = new Array<OrderView>(50)

    let config = this.props.node.getConfig()

    let { defaultCols, defaultColWidths } = getColumnState(columns, config);

    let blotterState: OrderBlotterState = {
      orders: view,
      selectedOrders: new Map<string, Order>(),
      columns: defaultCols,
      columnWidths: defaultColWidths
    }

    this.state = blotterState;

    this.props.node.setEventListener("save", (p) => {

      let cols = this.state.columns
      let colOrderIds = getColIdsInOrder(cols);

      let persistentConfig: TableViewConfig = {
        columnWidths: this.state.columnWidths,
        columnOrder: colOrderIds
      }

      this.props.model.doAction(Actions.updateNodeAttributes(props.node.getId(), { config: persistentConfig }))
    }

    );

    this.listingService = props.listingService

    this.orderMap = new Map<string, number>();




    let after = new Timestamp()

    let startOfLocalDay = new Date()
    startOfLocalDay.setHours(0, 0, 0, 0)
    after.setSeconds(Math.floor(startOfLocalDay.getTime() / 1000))
    let sto = new SubscribeToOrders()
    sto.setAfter(after)

    this.stream = this.viewService.subscribe(sto, Login.grpcContext.grpcMetaData)

    this.stream.on('data', (order: Order) => {

      let idx = this.orderMap.get(order.getId())
      let orderView: OrderView
      let newOrders = [...this.state.orders]
      if (idx === undefined) {
        idx = this.orderMap.size
        let orderLength = this.state.orders.length
        if (idx >= orderLength) {
          newOrders = new Array<OrderView>(orderLength * 2)
          for (let i = 0; i < orderLength; i++) {
            newOrders[i] = this.state.orders[i]
          }
        }

        this.orderMap.set(order.getId(), idx);
        orderView = new OrderView(order)

        newOrders[idx] = orderView
      } else {
        orderView = new OrderView(order)
        newOrders[idx] = orderView
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


      this.listingService.GetListing(order.getListingid(), (listing) => {
        let params = new CancelOrderParams()
        params.setOrderid(order.getId())
        params.setListing(listing)

        this.executionVenueService.cancelOrder(params, Login.grpcContext.grpcMetaData, (err: grpcWeb.Error, response: Empty) => {
          if (err) {
            logGrpcError("error cancelling order", err)
          }
        })

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
          bodyContextMenuRenderer={this.renderBodyContextMenu} onSelection={this.onSelection} enableColumnReordering={true}
          onColumnsReordered={this.onColumnsReordered} enableColumnResizing={true} onColumnWidthChanged={this.columnResized} columnWidths={this.state.columnWidths}>
          {this.state.columns}
        </Table>
      </div>
    );
  }

  columnResized = (index: number, size: number) => {
    let newColWidths = this.state.columnWidths.slice();
    newColWidths[index] = size
    let blotterState: OrderBlotterState = {
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



