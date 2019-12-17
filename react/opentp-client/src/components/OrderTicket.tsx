import { AnchorButton, Classes, Colors, Dialog, FormGroup, Intent, Label, NumericInput, Spinner } from '@blueprintjs/core';
import { Error } from 'grpc-web';
import React, { CSSProperties } from 'react';
import { getListingLongName, getListingShortName } from '../common/modelutilities';
import { logGrpcError, logDebug } from '../logging/Logging';
import { Decimal64 } from '../serverapi/common_pb';
import { ExecutionVenueClient } from '../serverapi/Execution-venueServiceClientPb';
import { CreateAndRouteOrderParams, OrderId } from '../serverapi/execution-venue_pb';
import { Listing } from '../serverapi/listing_pb';
import { Side } from '../serverapi/order_pb';
import { toNumber } from '../util/decimal64Conversion';
import { ListingContext, TicketController } from './Container';
import Login from './Login';
import { QuoteService, QuoteListener } from '../services/QuoteService';
import { Quote } from '../serverapi/market-data-service_pb';

interface OrderTicketState {
  listing?: Listing,
  quote?: Quote,
  quantity: number,
  price: number,
  side: Side,
  isOpen: boolean,
  usePortal: boolean
}

interface OrderTicketProps {
  tickerController: TicketController,
  quoteService: QuoteService
}

export default class OrderTicket extends React.Component<OrderTicketProps, OrderTicketState> implements QuoteListener {

  executionVenueService = new ExecutionVenueClient(Login.grpcContext.serviceUrl, null, null)
  quoteService: QuoteService

  constructor(props: OrderTicketProps) {
    super(props);

    this.quoteService = props.quoteService
    props.tickerController.setOrderTicket(this)

    this.state = {
      quantity: 0,
      price: 0,
      side: Side.BUY,
      isOpen: false,
      usePortal: true,
    };

    this.sendOrder = this.sendOrder.bind(this);
  }

  onQuote(recQuote: Quote): void {
      let state: OrderTicketState = {
        ...this.state, ... {
          quote: recQuote
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

  private getListingShortName(): string {
    let side = this.state.side
    if (this.state && this.state.listing && side !== undefined) {

      return this.getSideAsString(side) + " " + getListingShortName(this.state.listing)
    }

    return " "
  }

  private getListingFullName(): string {
    let side = this.state.side
    if (this.state && this.state.listing && side !== undefined) {

      return getListingLongName(this.state.listing)
    }

    return " "

  }

  // TOdo - bind in tick size to numeric HTMLFormControlsCollection, use quote to populate in, add hotkeys for ticket 

  public render() {

    let listing = this.state.listing
    //let quote = this.state.quote
    if (listing ) {

      let sizeIncrement = toNumber(listing.getSizeincrement())

      return (

        <Dialog
          icon="exchange"
          onClose={this.handleClose}
          title={this.getListingShortName()}
          {...this.state}
          className="bp3-dark">
          <div className={Classes.DIALOG_BODY}>

            <Label>{this.getListingFullName()}</Label>
            <FormGroup
              label="Quantity"
              labelFor="quantity-input">
              <NumericInput
                id="quantity-input"
                stepSize={sizeIncrement}
                value={this.state.quantity}
                onChange={
                  (e: any) => {
                    this.setState({ quantity: e.target.value })
                  }

                }
              />
            </FormGroup>
            <FormGroup
              label="Price"
              labelFor="price-input">
              <NumericInput
                id="price-input"
                value={this.state.price}
                onChange={
                  (e: any) => {
                    this.setState({ price: e.target.value })
                  }

                }
              />
            </FormGroup>

          </div>
          <div className={Classes.DIALOG_FOOTER}>
            <div className={Classes.DIALOG_FOOTER_ACTIONS}>
              <AnchorButton onClick={this.sendOrder}
                intent={Intent.PRIMARY} style={this.getBuySellButtonStyle(this.state.side)}><h2>
                  {this.getSideAsString(this.state.side)}</h2>
              </AnchorButton>
            </div>
          </div>


        </Dialog>
      );
    } else {
      return <div></div>
    }

  }


  private getBuySellButtonStyle(side: Side): CSSProperties {

    let color = Colors.DARK_GRAY1
    switch (side) {
      case Side.BUY:
        color = Colors.BLUE5
        break
      case Side.SELL:
        color = Colors.ROSE4
        break

    }

    return { background: color }
  }

  openTicket(newSide: Side, newListing: Listing) {

    logDebug("opening ticket")

    let existingQuote = this.quoteService.SubscribeToQuote(newListing.getId(), this)

    let defaultPrice;
    let defaultQuantity;
    if (existingQuote) {
      if (existingQuote.getDepthList().length > 0) {
        let depthLine = existingQuote.getDepthList()[0]
        if (this.state.side === Side.SELL) {
          defaultPrice = toNumber(depthLine.getBidprice())
          defaultQuantity = toNumber(depthLine.getBidsize())
        } else {
          defaultPrice = toNumber(depthLine.getAskprice())
          defaultQuantity = toNumber(depthLine.getAsksize())
        }
      }
    }

    if (!defaultPrice) {
      defaultPrice = 0
    }

    if (!defaultQuantity) {
      defaultQuantity = 0
    }

    let state: OrderTicketState = {
      ...this.state, ... {
        side: newSide,
        isOpen: true,
        listing: newListing,
        price: defaultPrice,
        quantity: defaultQuantity,
        quote: existingQuote
      }
    }

    this.setState(state)
    logDebug("opened ticket")
  }

  private handleClose = () => {

    logDebug("closing ticket")

    if (this.state.listing) {
      this.quoteService.UnsubscribeFromQuote(this.state.listing.getId(), this)
    }

    this.setState({
      ...this.state, ... {
        isOpen: false,
        defaultPrice: undefined,
        defaultQuantity: undefined,
        price: 0,
        quantity: 0,
        listing: undefined,
        quote: undefined
      }
    })

    logDebug("closed ticket")
  };


  private sendOrder(event: React.MouseEvent<HTMLElement>) {

    let listing = this.state.listing
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