import React from 'react';
import query from 'query-string';
import { ExecutionVenueClient } from '../serverapi/Execution-venueServiceClientPb';
import Login from './Login';
import { CreateAndRouteOrderParams, OrderId } from '../serverapi/execution-venue_pb';
import { logError, logGrpcError } from '../logging/Logging';
import { Error } from 'grpc-web';
import { Side } from '../serverapi/order_pb';
import { Decimal64 } from '../serverapi/common_pb';
import { ListingContext } from './Container';
import { Listing } from '../serverapi/listing_pb';


interface OrderTicketState {
    listing?: Listing,
    quantity: number,
    price: number,
    side: string,
    instrumentId: string
}

interface OrderTicketProps {

  listingContext: ListingContext
}



export default class OrderTicket extends React.Component<OrderTicketProps , OrderTicketState> {

  executionVenueService = new ExecutionVenueClient(Login.grpcContext.serviceUrl, null, null)

    constructor(props: OrderTicketProps) {
        super(props);

        this.state = {
          quantity : 0,
          price : 0,
          side: '',
          instrumentId: ''
        };


        props.listingContext.addListener((listing:Listing)=>{
          let state: OrderTicketState = {
            ...this.state, ... {
              listing: listing,
            }
          }

          this.setState(state)
        })

    }      

   

    public render() {


        return ( 
            
          <div>
            <div>
            <label>Side </label>
            <input
              type="text"
              name="side"
              value={this.state.side}
              onChange={
                (e: any) =>{
                    this.setState({side: e.target.value})
                }
                    
              }
            />
            </div>
            <div>
            <label>Quantity </label>
            <input
              type="text"
              name="quantity"
              value={this.state.quantity}
              onChange={
                (e: any) =>{
                    this.setState({quantity: e.target.value})
                }
                    
              }
            />
            </div>
            <div>
            <label>Instrument </label>
            <input
              type="text"
              name="instrument"
              value={this.state.instrumentId}
              onChange={
                (e: any) =>{
                    this.setState({instrumentId: e.target.value})
                }

              }
            />
            </div>
            <div>
            <label>Price </label>
            <input
              type="text"
              name="price"
              value={this.state.price}
              onChange={
                (e: any) =>{
                    this.setState({price: e.target.value})
                }

              }
            />
            </div>
            <div>
              <button onClick={e=>this.sendOrder(this.state) }>Send</button>
            </div>
          </div>
        );

        }



        sendOrder(params: OrderTicketState ) {

          let croParams = new CreateAndRouteOrderParams()
          croParams.setListingid(params.instrumentId)
          
          let side: Side;
          if (params.side.toUpperCase() === "BUY") {
              side = Side.BUY
          } else if( params.side.toUpperCase() === "SELL") {
            side = Side.SELL
          } else {
            logError("Unrecognised side:" + params.side)
            return
          }

          croParams.setSide(side)
          croParams.setQuantity(new Decimal64())

          this.executionVenueService.createAndRouteOrder(new CreateAndRouteOrderParams(), Login.grpcContext.grpcMetaData, (err: Error,
            response: OrderId) => {
            if( err ) {
              logGrpcError("failed to send order:", err)
            }
          })
          
        }

        
       




}