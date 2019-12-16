import { Dialog, Classes, Tooltip, Button, AnchorButton, Intent } from '@blueprintjs/core';
import { Error } from 'grpc-web';
import React from 'react';
import { getListingLabel } from '../common/modelutilities';
import { logDebug, logGrpcError } from '../logging/Logging';
import { Decimal64 } from '../serverapi/common_pb';
import { ExecutionVenueClient } from '../serverapi/Execution-venueServiceClientPb';
import { CreateAndRouteOrderParams, OrderId } from '../serverapi/execution-venue_pb';
import { Listing } from '../serverapi/listing_pb';
import { Side } from '../serverapi/order_pb';
import { ListingContext, TicketController } from './Container';
import Login from './Login';


interface OrderTicketState {
  listing?: Listing,
  quantity: number,
  price: number,
  side: Side,
  isOpen: boolean,
  usePortal: boolean
}

interface OrderTicketProps {
  listingContext: ListingContext,
  tickerController: TicketController
}


export default class OrderTicket extends React.Component<OrderTicketProps, OrderTicketState> {

  executionVenueService = new ExecutionVenueClient(Login.grpcContext.serviceUrl, null, null)
  listingContext: ListingContext

  constructor(props: OrderTicketProps) {
    super(props);

    this.listingContext = props.listingContext
    props.tickerController.setOrderTicket(this)

    this.state = {
      quantity: 0,
      price: 0,
      side: Side.BUY,
      isOpen: false,
      usePortal: true,
    };

    props.listingContext.addListener((listing: Listing) => {
      let state: OrderTicketState = {
        ...this.state, ... {
          listing: listing,
        }
      }

      this.setState(state)

    })


    this.sendOrder = this.sendOrder.bind(this);
  }

  openTicket(newSide: Side) {

    let state: OrderTicketState = {
      ...this.state, ... {
        side: newSide,
        isOpen: true
      }
    }


    this.setState(state)


  }

  private getSideAsString(side: Side): string {
    switch (side) {
      case Side.BUY:
        return "BUY"
      case Side.SELL:
        return "SELL"
      default:
        return "Side not recognised:" + side

    }
  }

  private getListingLabel(): string {
    let side = this.state.side
    if (this.state && this.state.listing && side != undefined) {

      return this.getSideAsString(side) + " " + getListingLabel(this.state.listing)
    }

    return " "
  }

  public render() {

    return (

      <Dialog
        icon="exchange"
        onClose={this.handleClose}
        title={this.getListingLabel()}
        {...this.state}
        className="bp3-dark"
      >
        <div className={Classes.DIALOG_BODY}>

          <div>
            <label>Quantity </label>
            <input
              type="text"
              name="quantity"
              value={this.state.quantity}
              onChange={
                (e: any) => {
                  this.setState({ quantity: e.target.value })
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
                (e: any) => {
                  this.setState({ price: e.target.value })
                }

              }
            />
          </div>

        </div>
        <div className={Classes.DIALOG_FOOTER}>
          <div className={Classes.DIALOG_FOOTER_ACTIONS}>
            <AnchorButton onClick={this.sendOrder}
              intent={Intent.PRIMARY}><h2>
              {this.getSideAsString(this.state.side)}</h2>
            </AnchorButton>
          </div>
        </div>


      </Dialog>
    );

  }


  private handleClose = () => this.setState({
    ...this.state, ... {
      isOpen: false
    }
  });


  private sendOrder(event: React.MouseEvent<HTMLElement>) {

    let listing = this.listingContext.selectedListing
    let side = this.state.side
    if (listing && side) {

      let croParams = new CreateAndRouteOrderParams()
      croParams.setListing(listing)



      croParams.setSide(side)
      croParams.setQuantity(new Decimal64())

      this.executionVenueService.createAndRouteOrder(new CreateAndRouteOrderParams(), Login.grpcContext.grpcMetaData, (err: Error,
        response: OrderId) => {
        if (err) {
          logGrpcError("failed to send order:", err)
        }
      })

    }



  }







}