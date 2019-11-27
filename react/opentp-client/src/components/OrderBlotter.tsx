import React from 'react';
import './OrderBlotter.css';
import ReactTable, { RowInfo } from 'react-table';
import "react-table/react-table.css";
import v4 from 'uuid';
import { Order, OrderStatus, Side } from '../model/Model';
import { ContextMenu, MenuItem, ContextMenuTrigger } from "react-contextmenu";
 

interface BlotterState {
  orders : Order[];
}

interface Props {
  onOrderSelected: (order: Order) => void
  selectedOrder?: Order;
}

export default class OrderBlotter extends React.Component<Props, BlotterState > {

    orderMap : Map<string, Order>;

    ordersSource : EventSource;
    id : string;

    constructor(props: Props) {
        super(props);

        this.id = v4();

        this.orderMap = new Map<string, Order>();
        /*
        this.orderMap.set("1a",  {
          id: "1a",
          instrumentId: "abc",
          qty: 30,
          price: 45,
          orderStatus: OrderStatus.Live,
          targetStatus: OrderStatus.None,
          side: Side.Buy

        })
        this.orderMap.set("2a",  {
          id: "2a",
          instrumentId: "abc2",
          qty: 30,
          price: 45,
          orderStatus: OrderStatus.Live,
          targetStatus: OrderStatus.None,
          side: Side.Buy

        })
        this.orderMap.set("3a",  {
          id: "3a",
          instrumentId: "abc3",
          qty: 30,
          price: 45,
          orderStatus: OrderStatus.Live,
          targetStatus: OrderStatus.None,
          side: Side.Buy

        })
        
        */
        let blotterState : BlotterState = {
          orders: Array.from(this.orderMap.values())
        }

        this.state =  blotterState;
        
        
        this.ordersSource = new EventSource("http://192.168.1.100:31887/proxy/subscribe-to-topic/orders?subscriberId=" + this.id);

        this.ordersSource.addEventListener( "orders", e  => {

          console.log("Message Event " + e);
          
          const messageEvent =  e as MessageEvent;

          
          let order : Order  = JSON.parse(messageEvent.data) as Order;

          this.orderMap.set(order.id, order);

          console.log("Values " + this.orderMap.values());


          let blotterState : BlotterState = {
            orders: Array.from(this.orderMap.values()),
            
            //selectionChanged: this.state.selectionChanged
          }

          this.setState( blotterState );
        })

        this.ordersSource.onerror = function(e) {
          console.log("EventSource failed." + e);
        };

        this.ordersSource.onopen = ( e: Event) => {
          console.log("Opened SSE connection")
        }; 

    }

    
    cancelOrder = (e: any, data: Order) => {
      if(data) {
        http://192.168.1.102:32413/order-management/cancel-order?orderId=00a2fdb5-7521-4f44-a985-2cddc7a19222

        fetch('http://192.168.1.102:32413/order-management/cancel-order?orderId=' + data.id, {
          method: 'POST',
          mode: 'no-cors'
        })
          .then(
            response =>{ console.log(response.statusText) } 
            )
          
          .catch(error => {
            throw new Error(error);
          });
      }
    }

    modifyOrder = (e: any, data: Order) => {
      if(data) {
        window.alert("modify order" + data.id);
      }
    }


    public render() {

        
        const myClonedArray  =  Object.assign([], this.state.orders);

        return ( 
              <div>
                
                <ContextMenuTrigger id="orderblottermenu">
                    
                  

                <ReactTable<Order> 
                
                  
                  data={myClonedArray}
                  columns={[
                    {
                      columns: [
                        {
                          Header: "Order Id",
                          accessor: "id"
                        },
                        {
                          Header: "Instrument Id",
                          accessor: "instrumentId"
                        },
                        {
                          Header: "Status",
                          accessor: "orderStatus"
                        },
                        {
                          Header: "Target Status",
                          accessor: "targetStatus"
                        },
                        {
                          Header: "Price",
                          accessor: "price"
                        },
                        {
                          Header: "Quantity",
                          accessor: "qty"
                        },
                        {
                          Header: "Side",
                          accessor: "side"
                        },
                        {
                          Header: "Traded Quantity",
                          accessor: "tradedQty"
                        },
                        {
                          Header: "Remaining Quantity",
                          accessor: "remainingQty"
                        },
                        {
                          Header: "Error",
                          accessor: "errorMsg"
                        }
                      ]
                    }
                  ]}
                  
                  showPaginationBottom={false}
                  defaultPageSize={200}
                  style={{
                    height: 20*41 + "px" // This will force the table body to overflow and scroll, since there is not enough room
                  }}
                  className="-striped -highlight"
                  
                  getTrProps={(state: any , rowInfo : RowInfo | undefined) => {
                
                    if (rowInfo && rowInfo.original ) {

                      let backgroundstyle :any;
                      if( this.props.selectedOrder && this.props.selectedOrder.id === rowInfo.original.id) {
                        backgroundstyle = {
                          background: '#00afec'
                        }
                      } else {
                        backgroundstyle = {};
                      }

                      return {

                       onClick: (e: any) => {
                            
                           let blotterState : BlotterState = {
                            orders: Array.from(this.orderMap.values()),
                          }
                          this.props.onOrderSelected(rowInfo.original)

                          this.setState(blotterState);
                          
                        }, 
                        style: backgroundstyle
                      }
                    }else{
                      return {}
                    }
                  }} 



                />)
                <br />

                </ContextMenuTrigger>
            
            <ContextMenu id="orderblottermenu" >
              <MenuItem data={this.props.selectedOrder} onClick={this.cancelOrder}  >
                Cancel Order
              </MenuItem>
              <MenuItem divider />
              <MenuItem data={this.props.selectedOrder}  onClick={this.modifyOrder}>
                Modify Order
              </MenuItem>
              
            </ContextMenu>

              </div>


        );
    }

}