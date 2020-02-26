import { AnchorButton, Classes, Colors, Dialog, FormGroup, Intent, Label, NumericInput } from '@blueprintjs/core';
import { Error } from 'grpc-web';
import React, { CSSProperties } from 'react';
import { getListingLongName, getListingShortName } from '../common/modelutilities';
import { logGrpcError, logDebug } from '../logging/Logging';
import { Decimal64 } from '../serverapi/common_pb';
import { ExecutionVenueClient } from '../serverapi/Execution-venueServiceClientPb';
import { CreateAndRouteOrderParams, OrderId } from '../serverapi/execution-venue_pb';
import { Listing } from '../serverapi/listing_pb';
import { TickSizeEntry } from '../serverapi/listing_pb';
import { Side } from '../serverapi/order_pb';
import { toNumber, toDecimal64 } from '../util/decimal64Conversion';
import Login from './Login';
import { QuoteService, QuoteListener } from '../services/QuoteService';
import {  TicketController } from "./Container";
import { ClobQuote } from '../serverapi/clobquote_pb';

interface OrderTicketState {
  listing?: Listing,
  quote?: ClobQuote,
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

  onQuote(recQuote: ClobQuote): void {
      let state: OrderTicketState = {
        ...this.state,...{
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


  private getTickSize(price:number, listing:Listing):number {
    let tt = listing.getTicksize()
    if( tt ) {
      let el = tt.getEntriesList()
      for(var entry of el) {
        let tickSize = this.tickSizeFromEntry(entry, price)
        if( tickSize ) {
          return tickSize
        }
      }
    }

    return 1
  }

  private tickSizeFromEntry(entry :TickSizeEntry, price: number  ):number | undefined {
    let lowerBound = toNumber(entry.getLowerpricebound())
    let upperBound = toNumber(entry.getUpperpricebound())

    if( lowerBound !== undefined && upperBound !== undefined ) {
      if( price >= lowerBound && price <=upperBound) {
        return toNumber(entry.getTicksize())
      }
    }
    
    return undefined
  }


  private getAskText(quote?: ClobQuote):string {
    if( quote ) {
      let best = this.getBestBidAndAsk(quote)
      return "Ask: " + best.bestAskQuantity +"@" + best.bestAskPrice
    } else {
      return "Ask: <>"
    }
  }

  private getBidText(quote?: ClobQuote):string {
    if( quote ) {
      let best = this.getBestBidAndAsk(quote)
      return "Bid: " + best.bestBidQuantity +"@" + best.bestBidPrice
    } else {
      return "Bid: <>"
    }
  }

  public render() {

    let listing = this.state.listing
    //let quote = this.state.quote
    if (listing ) {

      let sizeIncrement = toNumber(listing.getSizeincrement())
      let tickSize = this.getTickSize(this.state.price, listing)

      return (

        <Dialog
          icon="exchange"
          onClose={this.handleClose}
          title={this.getListingShortName()}
          {...this.state}
          className="bp3-dark">
          <div className={Classes.DIALOG_BODY}>

            <Label>{this.getListingFullName()}</Label>
            <Label style={{color:Colors.LIME3}}>{this.getBidText(this.state.quote)}</Label>
            <Label style={{color:Colors.ORANGE3}}>{this.getAskText(this.state.quote)}</Label>

            <FormGroup
              label="Quantity"
              labelFor="quantity-input">
              <NumericInput
                id="quantity-input"
                stepSize={sizeIncrement}
                minorStepSize={sizeIncrement}
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
                stepSize={tickSize}
                minorStepSize={tickSize}
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


  getBestBidAndAsk(quote: ClobQuote): BestBidAndAsk {
    let result = new BestBidAndAsk()
    if (quote.getBidsList().length > 0) {
      let depthLine = quote.getBidsList()[0]

      result.bestBidPrice = toNumber(depthLine.getPrice())
      result.bestBidQuantity = toNumber(depthLine.getSize())
    }

    if (quote.getOffersList().length > 0) {
      let depthLine = quote.getOffersList()[0]

      result.bestAskPrice = toNumber(depthLine.getPrice())
      result.bestAskQuantity = toNumber(depthLine.getSize())
    }


    return result
  }



  openTicket(newSide: Side, newListing: Listing) {

    let existingQuote = this.quoteService.SubscribeToQuote(newListing.getId(), this)

    let defaultPrice;
    let defaultQuantity;
    if (existingQuote) {
      let best = this.getBestBidAndAsk(existingQuote)
      if (this.state.side === Side.SELL) {
        defaultPrice = best.bestBidPrice
        defaultQuantity = best.bestBidQuantity
      } else {
        defaultPrice = best.bestAskPrice
        defaultQuantity = best.bestAskQuantity
      }
    }

    if (!defaultPrice) {
      defaultPrice = 0
    }

    if (!defaultQuantity) {
      defaultQuantity = 0
    }

    let state: OrderTicketState = {
      ...this.state,...{
        side: newSide,
        isOpen: true,
        listing: newListing,
        price: defaultPrice,
        quantity: defaultQuantity,
        quote: existingQuote
      }
    }

    this.setState(state)
  }

  private handleClose = () => {

    if (this.state.listing) {
      this.quoteService.UnsubscribeFromQuote(this.state.listing.getId(), this)
    }

    this.setState({
      ...this.state,...{
        isOpen: false,
        defaultPrice: undefined,
        defaultQuantity: undefined,
        price: 0,
        quantity: 0,
        listing: undefined,
        quote: undefined
      }
    })

  };


  private sendOrder(event: React.MouseEvent<HTMLElement>) {

    let listing = this.state.listing
    let side = this.state.side
    if (listing) {

      let croParams = new CreateAndRouteOrderParams()
      croParams.setListing(listing)

      croParams.setSide(side)
      croParams.setQuantity(toDecimal64(this.state.quantity))
      croParams.setPrice(toDecimal64(this.state.price))

      this.executionVenueService.createAndRouteOrder(croParams, Login.grpcContext.grpcMetaData, (err: Error,
        response: OrderId) => {
        if (err) {
          logGrpcError("error whilst sending order:", err)
        }
        logDebug("create and route order result:" + response.getOrderid())
      })

    }

    this.handleClose()

  }

}

class BestBidAndAsk {
  bestBidPrice? :number
  bestBidQuantity? :number
  bestAskPrice? :number
  bestAskQuantity? :number
}