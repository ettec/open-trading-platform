import { AnchorButton, Classes, Colors, Dialog, FormGroup, Intent, Label, NumericInput, MenuItem, Button } from '@blueprintjs/core';
import { Error } from 'grpc-web';
import React, { CSSProperties } from 'react';
import { getListingLongName, getListingShortName } from '../../common/modelutilities';
import { logDebug, logError } from '../../logging/Logging';
import { ClobQuote } from '../../serverapi/clobquote_pb';
import { ExecutionVenueClient } from '../../serverapi/ExecutionvenueServiceClientPb';
import { CreateAndRouteOrderParams, OrderId, ModifyOrderParams } from '../../serverapi/executionvenue_pb';
import { Listing, TickSizeEntry } from '../../serverapi/listing_pb';
import { Side, Order } from '../../serverapi/order_pb';
import { QuoteListener, QuoteService } from '../../services/QuoteService';
import { toDecimal64, toNumber } from '../../util/decimal64Conversion';
import { TicketController } from "../Container";
import Login from '../Login';
import { GlobalColours } from '../Colours';
import { Select, ItemRenderer } from '@blueprintjs/select';
import VwapParams from './Strategies/VwapParams/VwapParams';

interface OrderTicketState {
  listing?: Listing,
  quote?: ClobQuote,
  quantity: number,
  price: number,
  side: Side,
  isOpen: boolean,
  usePortal: boolean
  orderToModify: Order | null,
  destination?: string,
  destinations: Array<string>,
}

interface OrderTicketProps {
  tickerController: TicketController,
  quoteService: QuoteService
}

const DestinationSelect = Select.ofType<string>();
const renderDestination: ItemRenderer<string> = (destination, { handleClick, modifiers, query }) => {

  const text = `${destination}`;
  return (
    <MenuItem
      active={modifiers.active}
      disabled={modifiers.disabled}

      key={destination}
      onClick={handleClick}
      text={text}
    />
  );
};

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
      orderToModify: null,
      destinations: new Array<string>()
    };

    this.onSubmit = this.onSubmit.bind(this);
    this.handleQueryChange = this.handleQueryChange.bind(this);
  }

  onQuote(recQuote: ClobQuote): void {
    let state: OrderTicketState = {
      ...this.state, ...{
        quote: recQuote
      }
    }
    this.setState(state)
  }



  private getSubmitButtonAsString(): string {

    if (this.state.orderToModify) {
      return "Modify"
    } else {
      switch (this.state.side) {
        case Side.BUY:
          return "BUY"
        case Side.SELL:
          return "SELL"
        default:
          return "Side not recognised:" + this.state.side

      }

    }


  }

  private getDialogTitle(): string {

    if (this.state.orderToModify) {
      return "Modify Order " + this.state.orderToModify.getId()
    }


    let side = this.state.side
    if (this.state && this.state.listing && side !== undefined) {

      return this.getSubmitButtonAsString() + " " + getListingShortName(this.state.listing)
    }

    return " "
  }

  private getListingFullName(): string {
    if (this.state && this.state.listing) {
      return getListingLongName(this.state.listing)
    }

    return " "

  }


  private getTickSize(price: number, listing: Listing): number {
    let tt = listing.getTicksize()
    if (tt) {
      let el = tt.getEntriesList()
      for (var entry of el) {
        let tickSize = this.tickSizeFromEntry(entry, price)
        if (tickSize) {
          return tickSize
        }
      }
    }

    return 1
  }

  private tickSizeFromEntry(entry: TickSizeEntry, price: number): number | undefined {
    let lowerBound = toNumber(entry.getLowerpricebound())
    let upperBound = toNumber(entry.getUpperpricebound())

    if (lowerBound !== undefined && upperBound !== undefined) {
      if (price >= lowerBound && price <= upperBound) {
        return toNumber(entry.getTicksize())
      }
    }

    return undefined
  }


  private getAskText(quote?: ClobQuote): string {
    if (quote) {
      let best = this.getBestBidAndAsk(quote)
      if (best.bestAskPrice && best.bestAskQuantity) {
        return "Ask: " + best.bestAskQuantity + "@" + best.bestAskPrice
      }
    }
    return "Ask: <>"

  }

  private getBidText(quote?: ClobQuote): string {
    if (quote) {
      let best = this.getBestBidAndAsk(quote)
      if (best.bestBidPrice && best.bestBidQuantity) {
        return "Bid: " + best.bestBidQuantity + "@" + best.bestBidPrice
      }
    }

    return "Bid: <>"

  }



  public render() {

    let listing = this.state.listing
    if (listing) {

      let sizeIncrement = toNumber(listing.getSizeincrement())
      let tickSize = this.getTickSize(this.state.price, listing)


      const destination = this.state.destination;

      return (

        <Dialog
          icon="exchange"
          onClose={this.handleClose}
          title={this.getDialogTitle()}
          {...this.state}
          className="bp3-dark">
          <div className={Classes.DIALOG_BODY}>

            <Label>{this.getListingFullName()}</Label>
            <Label style={{ color: Colors.LIME3 }}>{this.getBidText(this.state.quote)}</Label>
            <Label style={{ color: Colors.ORANGE3 }}>{this.getAskText(this.state.quote)}</Label>

            <FormGroup
              label="Quantity"
              labelFor="quantity-input">
              <NumericInput
                id="quantity-input"
                stepSize={sizeIncrement}
                minorStepSize={sizeIncrement}
                value={this.state.quantity}
                onValueChange={
                  (num: number) => {
                    this.setState({ quantity: num })
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
                onValueChange={
                  (num: number) => {
                    this.setState({ price: num })
                  }

                }
              />
            </FormGroup>

            <Label htmlFor="input-b">Destination</Label>
            <DestinationSelect items={this.state.destinations}
              resetOnClose={true}
              onItemSelect={this.handleValueChange}
              onQueryChange={this.handleQueryChange}
              itemRenderer={renderDestination}
              noResults={<MenuItem disabled={true} text="No results." />}>
              <Button
                rightIcon="caret-down"
                text={destination ? destination : "Select Destination..."} />
            </DestinationSelect>

            <VwapParams></VwapParams>
          </div>
          <div className={Classes.DIALOG_FOOTER}>
            <div className={Classes.DIALOG_FOOTER_ACTIONS}>
              <AnchorButton onClick={this.onSubmit}
                intent={Intent.PRIMARY} style={this.getBuySellButtonStyle(this.state.side)}><h2>
                  {this.getSubmitButtonAsString()}</h2>
              </AnchorButton>
            </div>
          </div>


        </Dialog>
      );
    } else {
      return <div></div>
    }

  }

  private handleValueChange = (destination: string) => {
    this.setState({
      ...this.state, ...{
        destination: destination,
      }
    })
  };


  handleQueryChange(query: string) {

    let allDestinations = ["DMA", "VWAP"]
    let filteredDestinations = new Array<string>()

    if (query.length > 0) {
      query = query.toUpperCase()

      for (let dest of allDestinations) {
        if (dest.startsWith(query)) {
          filteredDestinations.push(dest)
        }
      }

    } else {
      filteredDestinations = allDestinations
    }

    let newState = {
      ...this.state, ...{
      }
    }

    newState.destinations = filteredDestinations

    this.setState(newState)

  }






  private getBuySellButtonStyle(side: Side): CSSProperties {

    let color = Colors.DARK_GRAY1
    switch (side) {
      case Side.BUY:
        color = GlobalColours.BUYBKG
        break
      case Side.SELL:
        color = GlobalColours.SELLBKG
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


  openModifyOrderTicket(order: Order, newListing: Listing) {

    let existingQuote = this.quoteService.SubscribeToQuote(newListing, this)


    let price = toNumber(order.getPrice())
    if (!price) {
      price = 0
    }

    let quantity = toNumber(order.getQuantity())
    if (!quantity) {
      quantity = 0
    }


    let destinations = new Array<string>()
    let destination = "DMA"
    if (order.getOwnerid() !== newListing.getMarket()?.getMic()) {
      destination = order.getOwnerid()
      destinations.push(order.getOwnerid())
    } else {
      destinations.push("DMA")
    }


    let state: OrderTicketState =
    {
      side: order.getSide(),
      isOpen: true,
      listing: newListing,
      price: price,
      quantity: quantity,
      quote: existingQuote,
      orderToModify: order,
      usePortal: true,
      destinations: destinations,
      destination: destination
    }


    this.setState(state)


  }


  openNewOrderTicket(newSide: Side, newListing: Listing) {

    let existingQuote = this.quoteService.SubscribeToQuote(newListing, this)

    let defaultPrice;
    let defaultQuantity;
    if (existingQuote) {
      let best = this.getBestBidAndAsk(existingQuote)
      if (newSide === Side.SELL) {
        defaultPrice = best.bestBidPrice
        defaultQuantity = best.bestBidQuantity
      } else {
        defaultPrice = best.bestAskPrice
        defaultQuantity = best.bestAskQuantity
      }
    }

    this.openOrderTicketWithDefaultPriceAndQty(newSide, newListing, defaultPrice, defaultQuantity,);
  }

  private handleClose = () => {

    if (this.state.listing) {
      this.quoteService.UnsubscribeFromQuote(this.state.listing.getId(), this)
    }

    this.setState({
      ...this.state, ...{
        isOpen: false,
      }
    })

  };



  public openOrderTicketWithDefaultPriceAndQty(newSide: Side, newListing: Listing, defaultPrice?: number, defaultQuantity?: number,) {

    if (!defaultPrice) {
      defaultPrice = 0;
    }

    if (!defaultQuantity) {
      defaultQuantity = 0;
    }

    let existingQuote = this.quoteService.SubscribeToQuote(newListing, this)

    let state: OrderTicketState = {
      side: newSide,
      isOpen: true,
      listing: newListing,
      price: defaultPrice,
      quantity: defaultQuantity,
      quote: existingQuote,
      usePortal: true,
      orderToModify: null,
      destinations: ["DMA", "VWAP"],
      destination: "DMA"

    };

    this.setState(state);
    return { defaultPrice, defaultQuantity };
  }

  private onSubmit(event: React.MouseEvent<HTMLElement>) {


    if (this.state.orderToModify) {

      let modifyParams = new ModifyOrderParams()
      modifyParams.setListingid(this.state.orderToModify.getListingid())
      modifyParams.setOwnerid(this.state.orderToModify.getOwnerid())
      modifyParams.setQuantity(toDecimal64(this.state.quantity))
      modifyParams.setPrice(toDecimal64(this.state.price))
      modifyParams.setOrderid(this.state.orderToModify.getId())

      logDebug("modify order " + this.state.orderToModify.getId() + " to " + toNumber(modifyParams.getQuantity()) +
        "@" + toNumber(modifyParams.getPrice()))

      this.executionVenueService.modifyOrder(modifyParams, Login.grpcContext.grpcMetaData, (err: Error) => {
        if (err) {
          let msg = "error whilst modifying order:" + err.message
          logError(msg)
          alert(msg)
        }
      })
    } else {
      let listing = this.state.listing
      let side = this.state.side
      let market = listing?.getMarket()
      if (listing && market) {

        let croParams = new CreateAndRouteOrderParams()
        croParams.setListingid(listing.getId())
        croParams.setDestination(market.getMic())
        croParams.setOrderside(side)
        croParams.setQuantity(toDecimal64(this.state.quantity))
        croParams.setPrice(toDecimal64(this.state.price))

        logDebug("sending order for " + toNumber(croParams.getQuantity()) + "@" + toNumber(croParams.getPrice()) + " of " +
          croParams.getListingid())
        croParams.setOriginatorid(Login.desk)
        croParams.setOriginatorref(Login.username)
        croParams.setRootoriginatorid(Login.desk)
        croParams.setRootoriginatorref(Login.username)

        this.executionVenueService.createAndRouteOrder(croParams, Login.grpcContext.grpcMetaData, (err: Error,
          response: OrderId) => {
          if (err) {
            let msg = "error whilst sending order:" + err.message
            logError(msg)
            alert(msg)
          } else {
            logDebug("create and route order created order with id:" + response.getOrderid())
          }

        })
      }
    }

    this.handleClose()

  }

}

class BestBidAndAsk {
  bestBidPrice?: number
  bestBidQuantity?: number
  bestAskPrice?: number
  bestAskQuantity?: number
}