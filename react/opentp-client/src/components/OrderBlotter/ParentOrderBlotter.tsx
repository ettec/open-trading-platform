import { Menu } from '@blueprintjs/core';
import { IMenuContext, IRegion, SelectionModes, Table } from "@blueprintjs/table";
import "@blueprintjs/table/lib/css/table.css";
import { Actions, Model, TabNode } from 'flexlayout-react';
import * as grpcWeb from 'grpc-web';
import React from 'react';
import v4 from 'uuid';
import log from 'loglevel';
import { Empty } from '../../serverapi/modelcommon_pb';
import { ExecutionVenueClient } from '../../serverapi/ExecutionvenueServiceClientPb';
import { CancelOrderParams } from '../../serverapi/executionvenue_pb';
import { Order, OrderStatus } from '../../serverapi/order_pb';
import { ListingService } from '../../services/ListingService';
import { OrderService } from '../../services/OrderService';
import { OrderContext } from "../Container/Contexts";
import { ChildOrderBlotterController, ExecutionsController, OrderHistoryBlotterController, TicketController } from "../Container/Controllers";
import Login from '../Login';
import '../TableView/TableCommon.css';
import { getColIdsInOrder, getConfiguredColumns, TableViewConfig, TableViewProperties } from '../TableView/TableView';
import OrderBlotter, { OrderBlotterState } from './OrderBlotter';
import { OrderView } from './OrderView';
import { Destinations } from '../../common/destinations';



interface ParentOrderBlotterState extends OrderBlotterState {
  selectedOrders: Array<Order>
}

interface ParentOrderBlotterProps extends TableViewProperties {
  node: TabNode,
  model: Model,
  orderContext: OrderContext
  orderService: OrderService
  childOrderBlotterController: ChildOrderBlotterController
  orderHistoryBlotterController: OrderHistoryBlotterController
  executionsController: ExecutionsController
  ticketController: TicketController
  listingService: ListingService
}



export default class ParentOrderBlotter
  extends OrderBlotter<ParentOrderBlotterProps, ParentOrderBlotterState> {


  executionVenueService = new ExecutionVenueClient(Login.grpcContext.serviceUrl, null, null)
  listingService: ListingService
  childOrderBlotterController: ChildOrderBlotterController
  orderHistoryBlotterController: OrderHistoryBlotterController
  executionsController: ExecutionsController
  ticketController: TicketController

  orderMap: Map<string, number>;

  orderService: OrderService

  id: string;


  constructor(props: ParentOrderBlotterProps) {
    super(props);

    this.id = v4();

    let columns = this.getColumns()

    let view = new Array<OrderView>()

    let config = props.node.getConfig()

    let [defaultCols, defaultColWidths] = getConfiguredColumns(columns, config);

    let visibleStates = new Set<OrderStatus>()
    visibleStates.add(OrderStatus.LIVE)
    visibleStates.add(OrderStatus.CANCELLED)
    visibleStates.add(OrderStatus.FILLED)
    visibleStates.add(OrderStatus.NONE)

    let blotterState: ParentOrderBlotterState = {
      orders: view,
      selectedOrders: new Array<Order>(),
      columns: defaultCols,
      columnWidths: defaultColWidths,
      visibleStates: visibleStates
    }

    this.state = blotterState;

    props.node.setEventListener("save", (p) => {

      let cols = this.state.columns
      let colOrderIds = getColIdsInOrder(cols);

      let persistentConfig: TableViewConfig = {
        columnWidths: this.state.columnWidths,
        columnOrder: colOrderIds
      }

      props.model.doAction(Actions.updateNodeAttributes(props.node.getId(), { config: persistentConfig }))
    }

    );

    this.listingService = props.listingService
    this.childOrderBlotterController = props.childOrderBlotterController
    this.orderHistoryBlotterController = props.orderHistoryBlotterController
    this.executionsController = props.executionsController
    this.ticketController = props.ticketController
    this.orderService = props.orderService

    this.orderMap = new Map<string, number>();




  }

  protected getTableName(): string {
    return "Order Blotter"
  }

  public componentDidMount(): void {
    this.orderService.SubscribeToAllParentOrders((order: Order) => {
      this.addOrUpdateOrder(order)
    })
  }


  showOrderHistory = (orders: IterableIterator<Order>) => {
    let order = orders.next()

    let cols = this.state.columns
    let colOrderIds = getColIdsInOrder(cols);

    let config: TableViewConfig = {
      columnWidths: this.state.columnWidths,
      columnOrder: colOrderIds
    }

    this.orderHistoryBlotterController.openBlotter(order.value, config,
      window.innerWidth)
  }



  showExecutions = (orders: IterableIterator<Order>) => {
    let order = orders.next()
    this.executionsController.open(order.value,
      window.innerWidth)
  }


  showChildOrders = (orders: IterableIterator<Order>) => {

    let order = orders.next()

    let childOrders = this.orderService.GetChildOrders(order.value)

    let cols = this.state.columns
    let colOrderIds = getColIdsInOrder(cols);

    let config: TableViewConfig = {
      columnWidths: this.state.columnWidths,
      columnOrder: colOrderIds
    }

    this.childOrderBlotterController.openBlotter(order.value, childOrders, config,
      window.innerWidth)

  }

  cancelOrder = (orders: Array<Order>) => {

    orders.forEach(order => {


      this.listingService.GetListing(order.getListingid(), (listing) => {
        let params = new CancelOrderParams()
        params.setOrderid(order.getId())
        params.setListingid(listing.getId())
        params.setOwnerid(order.getOwnerid())

        this.executionVenueService.cancelOrder(params, Login.grpcContext.grpcMetaData, (err: grpcWeb.Error, response: Empty) => {
          if (err) {
            log.error("error cancelling order", err)
          }
        })

      })


    });

  }

  modifyOrder = (order: Order) => {

    let listing = this.listingService.GetListingImmediate(order.getListingid())

    if (listing) {
      this.ticketController.openModifyOrderTicket(order, listing)
    }
  }

  public render() {

    return (
      <div style={{ width: "100%", height: "100%" }} >
        <Table enableRowResizing={false} numRows={this.state.orders.length} className="bp3-dark" selectionModes={SelectionModes.ROWS_AND_CELLS}
          bodyContextMenuRenderer={this.renderBodyContextMenu} onSelection={this.onSelection} enableColumnReordering={true}
          onColumnsReordered={this.onColumnsReordered} enableColumnResizing={true} onColumnWidthChanged={this.columnResized} columnWidths={this.state.columnWidths}>
          {this.state.columns}
        </Table>
      </div>
    );
  }



  private onSelection = (selectedRegions: IRegion[]) => {
    let newSelectedOrders: Array<Order> = this.getSelectedOrdersFromRegions(selectedRegions, this.state.orders);

    let blotterState: ParentOrderBlotterState = {
      ...this.state, ...{
        selectedOrders: newSelectedOrders,
      }
    }

    this.setState(blotterState)
  }


  private renderBodyContextMenu = (context: IMenuContext) => {

    let selectedOrders = this.getSelectedOrdersFromRegions(context.getRegions(), this.state.orders)
    let cancelleableOrders = OrderBlotter.cancelleableOrders(selectedOrders)

    return (

      <Menu  >
        <Menu.Item icon="delete" text="Cancel Order" onClick={() => this.cancelOrder(cancelleableOrders)} disabled={cancelleableOrders.length === 0} >
        </Menu.Item>
        <Menu.Divider />
        <Menu.Item icon="edit" text="Modify Order" onClick={() => this.modifyOrder(cancelleableOrders[0])} disabled={cancelleableOrders.length !== 1 ||
        cancelleableOrders[0].getOwnerid() === Destinations.VWAP || cancelleableOrders[0].getOwnerid() === Destinations.SMARTROUTER}>
        </Menu.Item>
        <Menu.Divider />
        <Menu.Item icon="fork" text="Child Orders" onClick={() => this.showChildOrders(selectedOrders.values())} disabled={selectedOrders.length === 0} >
        </Menu.Item>
        <Menu.Item icon="bring-data" text="History" onClick={() => this.showOrderHistory(selectedOrders.values())} disabled={selectedOrders.length === 0} >
        </Menu.Item>
        <Menu.Item icon="tick" text="Executions" onClick={() => this.showExecutions(selectedOrders.values())} disabled={selectedOrders.length === 0} >
        </Menu.Item>
        <Menu.Divider />
        <Menu.Item text="Edit Visible Columns" onClick={() => this.editVisibleColumns()}  >
        </Menu.Item>
      </Menu>
    );
  };

}



