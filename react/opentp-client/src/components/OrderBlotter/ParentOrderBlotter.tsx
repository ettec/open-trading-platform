import { Menu } from '@blueprintjs/core';
import { IMenuContext, IRegion, SelectionModes, Table } from "@blueprintjs/table";
import "@blueprintjs/table/lib/css/table.css";
import { Actions, Model, TabNode } from 'flexlayout-react';
import * as grpcWeb from 'grpc-web';
import React from 'react';
import v4 from 'uuid';
import { logGrpcError } from '../../logging/Logging';
import { Empty } from '../../serverapi/common_pb';
import { ExecutionVenueClient } from '../../serverapi/ExecutionvenueServiceClientPb';
import { CancelOrderParams } from '../../serverapi/executionvenue_pb';
import { Listing } from '../../serverapi/listing_pb';
import { Order } from '../../serverapi/order_pb';
import { ListingService } from '../../services/ListingService';
import { OrderService } from '../../services/OrderService';
import { ChildOrderBlotterController, OrderContext, OrderHistoryBlotterController, ExecutionsController, ColumnChooserController, TicketController } from '../Container';
import Login from '../Login';
import '../TableView/TableCommon.css';
import { getColIdsInOrder, getConfiguredColumns, TableViewConfig, TableViewState, TableViewProperties } from '../TableView/TableView';
import OrderBlotter from './OrderBlotter';
import { OrderView } from './OrderView';



interface ParentOrderBlotterState extends TableViewState {

  orders: OrderView[];
  selectedOrders: Array< Order>
  

}

interface ParentOrderBlotterProps extends TableViewProperties{
  node: TabNode,
  model: Model,
  orderContext: OrderContext
  listingService: ListingService
  orderService: OrderService
  childOrderBlotterController: ChildOrderBlotterController
  orderHistoryBlotterController: OrderHistoryBlotterController
  executionsController: ExecutionsController
  ticketController: TicketController
}



export default class ParentOrderBlotter extends OrderBlotter<ParentOrderBlotterProps, ParentOrderBlotterState> {


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

    let view = new Array<OrderView>(50)

    let config = this.props.node.getConfig()

    let [defaultCols, defaultColWidths] = getConfiguredColumns(columns, config);

  
    let blotterState: ParentOrderBlotterState = {
      orders: view,
      selectedOrders: new Array<Order>(),
      columns: defaultCols,
      columnWidths: defaultColWidths,
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
    this.childOrderBlotterController = props.childOrderBlotterController
    this.orderHistoryBlotterController = props.orderHistoryBlotterController 
    this.executionsController = props.executionsController
    this.ticketController = props.ticketController
    this.orderService = props.orderService

    this.orderMap = new Map<string, number>();


    

  }

  protected  getTableName() : string {
    return "Order Blotter"
}

public componentDidMount(): void {
  this.orderService.SubscribeToAllParentOrders((order: Order) => {

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

    newOrders.sort((o1,o2)=>{
      if( o1.created && o2.created) {
        return o1.created.getTime() - o2.created.getTime()
      }
      return 0
    })

    let blotterState: TableViewState = {
      ...this.state, ...{
        orders: newOrders
      }
    }

    orderView.listing = this.listingService.GetListingImmediate(order.getListingid())

    this.setState(blotterState);

    if (!orderView.listing) {
      this.listingService.GetListing(order.getListingid(), (listing: Listing) => {
        let newOrders = [...this.state.orders]
        let idx = this.orderMap.get(order.getId())
        orderView.listing = listing
        if (idx) {
          newOrders[idx] = orderView
          let blotterState: ParentOrderBlotterState = {
            ...this.state, ...{
              orders: newOrders
            }
          }

          this.setState(blotterState);
        }

      })
    }


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

    this.childOrderBlotterController.openBlotter(order.value, childOrders,  config,
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
            logGrpcError("error cancelling order", err)
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
      <div>

        <Table enableRowResizing={false} numRows={this.state.orders.length} className="bp3-dark" selectionModes={SelectionModes.ROWS_AND_CELLS}
          bodyContextMenuRenderer={this.renderBodyContextMenu} onSelection={this.onSelection} enableColumnReordering={true}
          onColumnsReordered={this.onColumnsReordered} enableColumnResizing={true} onColumnWidthChanged={this.columnResized} columnWidths={this.state.columnWidths}>
          {this.state.columns}
        </Table>
      </div>
    );
  }



  private onSelection = (selectedRegions: IRegion[]) => {
    let newSelectedOrders: Array< Order> = this.getSelectedOrdersFromRegions(selectedRegions, this.state.orders);

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
        <Menu.Item icon="delete"  text="Cancel Order" onClick={() => this.cancelOrder(cancelleableOrders)} disabled={cancelleableOrders.length === 0} >
        </Menu.Item>
        <Menu.Divider />
        <Menu.Item icon="edit" text="Modify Order" onClick={() => this.modifyOrder(cancelleableOrders[0])} disabled={cancelleableOrders.length === 0}>
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



