import { Checkbox, Colors, Icon, Menu } from '@blueprintjs/core';
import { Cell, Column, ColumnHeaderCell, IRegion } from "@blueprintjs/table";
import "@blueprintjs/table/lib/css/table.css";
import React from 'react';
import ReactCountryFlag from "react-country-flag";
import { roundToTick } from '../../common/modelutilities';
import { Order, OrderStatus, Side } from '../../serverapi/order_pb';
import { ListingService } from '../../services/ListingService';
import { GlobalColours } from '../Colours';
import '../TableView/TableCommon.css';
import TableView, { TableViewProperties } from '../TableView/TableView';
import { OrdersView, OrderView } from './OrderView';


export interface OrderBlotterProps extends TableViewProperties {
  listingService: ListingService
}

export interface OrderBlotterState  {
  orders: OrderView[];
  columns: Array<JSX.Element>
  columnWidths: Array<number>
  visibleStates: Set<OrderStatus>
}



export default abstract class OrderBlotter<P extends OrderBlotterProps , S extends OrderBlotterState> extends TableView<P , S>{

  private view : OrdersView

  constructor(props: P) {
    super(props)
    this.view = new OrdersView( props.listingService, ()=>{this.setState({ ...this.state, orders: this.view.getOrders() })})
    this.view.setSort((a:Order, b:Order) : number => {
       
      let aCreated = a.getCreated()
      let bCreated = b.getCreated() 
      if( aCreated  && bCreated) {
          return aCreated.getSeconds() - bCreated.getSeconds()
      } else {
          return 0
      }
  })

  }

 protected clearOrders() {
  this.view.clear()
  this.setState({ ...this.state, orders: this.view.getOrders() })
 }

  protected addOrUpdateOrder(order: Order) {
    this.view.addOrUpdateOrder(order)
    // Set state called twice, the table does not always update immediately unless this is done.
    this.setState({ ...this.state, orders: this.view.getOrders() })
    this.setState({ ...this.state, orders: this.view.getOrders() })
  }



  getColumns() {
    return [<Column key="id" id="id" name="Id" cellRenderer={this.renderId} />,
    <Column key="side" id="side" name="Side" cellRenderer={this.renderSide} />,
    <Column key="symbol" id="symbol" name="Symbol" cellRenderer={this.renderSymbol} />,
    <Column key="mic" id="mic" name="Mic" cellRenderer={this.renderMic} />,
    <Column key="country" id="country" name="Country" cellRenderer={this.renderCountry} />,
    <Column key="quantity" id="quantity" name="Quantity" cellRenderer={this.renderQuantity} />,
    <Column key="price" id="price" name="Price" cellRenderer={this.renderPrice} />,
    <Column key="status" id="status" name="Status" cellRenderer={this.renderStatus} columnHeaderCellRenderer={this.renderStatusHeader}/>,
    <Column key="targetStatus" id="targetStatus" name="Target Status" cellRenderer={this.renderTargetStatus} />,
    <Column key="remQty" id="remQty" name="Rem. Qty" cellRenderer={this.renderRemQty} />,
    <Column key="exposedQty" id="exposedQty" name="Exp. Qty" cellRenderer={this.renderExpQty} />,
    <Column key="tradedQty" id="tradedQty" name="Traded Qty" cellRenderer={this.renderTrdQty} />,
    <Column key="avgPrice" id="avgPrice" name="Avg Price" cellRenderer={this.renderAvgPrice} />,
    <Column key="listingId" id="listingId" name="Listing Id" cellRenderer={this.renderListingId} />,
    <Column key="created" id="created" name="Created" cellRenderer={this.renderCreated} />,
    <Column key="placedWith" id="placedWith" name="Placed With" cellRenderer={this.renderPlacedWith} />,
    <Column key="errorMsg" id="errorMsg" name="Error" cellRenderer={this.renderErrorMsg} />
    ];
  }

  private renderId = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.id}</Cell>;
  private renderQuantity = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.quantity}</Cell>;
  private renderSymbol = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.getSymbol()}</Cell>;
  private renderMic = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.getMic()}</Cell>;
  private renderPrice = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.price}</Cell>;
  private renderListingId = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.listingId}</Cell>;
  private renderRemQty = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.remainingQuantity}</Cell>;
  private renderExpQty = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.exposedQuantity}</Cell>;
  private renderTrdQty = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.tradedQuantity}</Cell>;
  private renderPlacedWith = (row: number) => <Cell>{Array.from(this.state.orders)[row]?.placedWith}</Cell>;
 
  private renderStatusHeader = (row: number) => {

      let icon = <Icon icon="filter" iconSize={12}></Icon>
      if( this.state.visibleStates.size !== 4) {
        icon = <Icon icon="filter-keep" iconSize={12}></Icon>
      }

    


      return <ColumnHeaderCell style={{display: 'flex',  flexDirection: 'row',  alignItems: "left", minWidth:40}} 
      name="Status" 
      menuRenderer={(index?: number): JSX.Element => {
        return <Menu><Checkbox checked={this.state.visibleStates.has(OrderStatus.LIVE)} onChange={()=>{this.onStatusFilterChange(OrderStatus.LIVE)}}>Live</Checkbox>
      <Checkbox checked={this.state.visibleStates.has(OrderStatus.FILLED) } onChange={()=>{this.onStatusFilterChange(OrderStatus.FILLED)}}>Filled</Checkbox>
      <Checkbox checked={this.state.visibleStates.has(OrderStatus.CANCELLED)} onChange={()=>{this.onStatusFilterChange(OrderStatus.CANCELLED)}}>Cancelled</Checkbox>
      <Checkbox checked={this.state.visibleStates.has(OrderStatus.NONE)} onChange={()=>{this.onStatusFilterChange(OrderStatus.NONE)}}>None</Checkbox>
      </Menu>}} >
        {icon}
      </ColumnHeaderCell>    
  }

  handleLiveChange() {
    let status = OrderStatus.LIVE
    this.onStatusFilterChange(status);
  }

 private onStatusFilterChange(status: OrderStatus) {
  let visibleStates = new Set<OrderStatus>();
  this.state.visibleStates.forEach((s) => { visibleStates.add(s); });


  if (visibleStates.has(status)) {
    visibleStates.delete(status);
  }
  else {
    visibleStates.add(status);
  }

  this.view.setFilter((order: Order)=>{return visibleStates.has(order.getStatus())})

  this.setState({ ...this.state, visibleStates: visibleStates, orders: this.view.getOrders() });
}
  
  private renderAvgPrice = (row: number) => {
    let orderView = Array.from(this.state.orders)[row]
    if( orderView && orderView.avgTradePrice) {
        if( orderView.listing ) {
          return <Cell>{roundToTick(orderView.avgTradePrice, orderView.listing )}</Cell>
        } else {
          return <Cell>{orderView.avgTradePrice}</Cell>
        }
    } else {
      return <Cell></Cell>
    }
  }

  private renderErrorMsg = (row: number) => {
    let orderView = Array.from(this.state.orders)[row]
    let statusStyle = {}
    if (orderView) {
      if (orderView.errorMsg.length >0 ) {
          statusStyle = { background: Colors.RED3 }  
      }
    }

    return <Cell style={statusStyle}>{orderView?.errorMsg}</Cell>
  }


  private renderCountry = (row: number) => {
    let country = Array.from(this.state.orders)[row]?.getCountryCode()

    if( country ) {
      return <Cell><ReactCountryFlag countryCode={country} /></Cell>
    } else {
      return <Cell></Cell>
    }

}

  private renderCreated = (row: number) => {
    let created = Array.from(this.state.orders)[row]?.created

    if (created) {
      return <Cell>{created.toLocaleTimeString()}</Cell>
    } else {
      return <Cell></Cell>
    }
  }


  private renderSide = (row: number) => {
    let orderView = Array.from(this.state.orders)[row]
    let statusStyle = {}
    if (orderView) {
      switch (orderView.getOrder().getSide()) {
        case Side.BUY:
          statusStyle = { background: GlobalColours.BUYBKG }
          break
        case Side.SELL:
          statusStyle = { background: GlobalColours.SELLBKG }
          break
      }
    }

    return <Cell style={statusStyle}>{orderView?.side}</Cell>
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


  getSelectedOrdersFromRegions(selectedRegions: IRegion[], orders: OrderView[]): Array<Order> {
    let newSelectedOrders: Array<Order> = new Array<Order>();

    let selectedOrderArray: Array<OrderView> = this.getSelectedItems(selectedRegions, orders);

    for (let orderView of selectedOrderArray) {
      newSelectedOrders.push(orderView.getOrder())
    }


    return newSelectedOrders;
  }



  static cancelleableOrders(orders: Array<Order>): Array<Order> {

    let result = new Array<Order>()
    for (let order of orders) {
      if (order.getStatus() === OrderStatus.LIVE) {
        result.push(order)
      }
    }

    return result
  }

}