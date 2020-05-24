import { Alignment, Button, Navbar } from "@blueprintjs/core";
import FlexLayout, { Model, TabNode } from "flexlayout-react";
import "flexlayout-react/style/dark.css";
import React from 'react';
import { Listing } from "../serverapi/listing_pb";
import { Order, Side } from "../serverapi/order_pb";
import ListingServiceImpl, { ListingService } from "../services/ListingService";
import OrderServiceImpl, { OrderService } from "../services/OrderService";
import QuoteServiceImpl, { QuoteService } from "../services/QuoteService";
import Executions from "./Executions";
import InstrumentListingWatch from "./InstrumentListingWatch";
import MarketDepth from './MarketDepth';
import ChildOrderBlotter from "./OrderBlotter/ChildOrderBlotter";
import OrderHistoryBlotter from "./OrderBlotter/OrderHistoryBlotter";
import OrderBlotter from "./OrderBlotter/ParentOrderBlotter";
import OrderTicket from './OrderTicket';
import { TableViewConfig } from "./TableView/TableView";
import QuestionDialog from "./QuestionDialog";
import { OrderMonitorClient } from "../serverapi/OrdermonitorServiceClientPb";
import Login from "./Login";
import { CancelAllOrdersForOriginatorIdParams } from "../serverapi/ordermonitor_pb";
import { Error } from "grpc-web";
import { Empty } from "../serverapi/modelcommon_pb";
import { logError, logDebug } from "../logging/Logging";



export default class Container extends React.Component {

    defaultJson = {
        global: {},
        borders: [],
        layout: {
            "type": "row",
            "weight": 100,
            "children": [

                {
                    "type": "row",
                    "weight": 50,
                    "children": [
                        {


                            "type": "tabset",
                            "weight": 50,
                            "children": [

                                {
                                    "type": "tab",
                                    "weight": 50,
                                    "name": "Instrument Watch",
                                    "component": "instrument-watch",

                                }

                            ]
                        },
                        {
                            "type": "tabset",
                            "weight": 50,
                            "children": [

                                {
                                    "type": "tab",
                                    "weight": 50,
                                    "name": "Order Blotter",
                                    "component": "order-blotter",
                                }

                            ]
                        }
                    ]
                },
                {
                    "type": "row",
                    "weight": 50,
                    "children": [
                        {


                            "type": "tabset",
                            "weight": 50,
                            "children": [

                                {
                                    "type": "tab",
                                    "weight": 50,
                                    "name": "Market Depth",
                                    "component": "market-depth",
                                }

                            ]
                        },
                        {


                            "type": "tabset",
                            "weight": 50,
                            "children": [

                                {
                                    "type": "tab",
                                    "weight": 50,
                                    "name": "Order Ticket",
                                    "component": "order-ticket",
                                }

                            ]
                        }


                    ]
                }
            ]
        }
    };


    orderMonitorClient = new OrderMonitorClient(Login.grpcContext.serviceUrl, null, null)

    state: Model;
    factory: (node: TabNode) => React.ReactNode;
    readonly configKey: string = "open-oms-config";

    quoteService: QuoteService
    orderService: OrderService
    listingService: ListingService
    listingContext: ListingContext
    orderContext: OrderContext
    ticketController: TicketController
    childOrderBlotterController : ChildOrderBlotterController
    orderHistoryBlotterController : OrderHistoryBlotterController
    executionsController : ExecutionsController
    questionDialogController: QuestionDialogController
    



    constructor() {
        super({}, {});

        this.listingService = new ListingServiceImpl()
        this.quoteService = new QuoteServiceImpl(this.listingService)
        this.orderService = new OrderServiceImpl()
        this.listingContext = new ListingContext()
        this.orderContext = new OrderContext()
        this.ticketController = new TicketController()
        this.childOrderBlotterController = new ChildOrderBlotterController()
        this.orderHistoryBlotterController = new OrderHistoryBlotterController()
        this.executionsController = new ExecutionsController()
        this.questionDialogController = new QuestionDialogController()

        let layoutString: string | null = localStorage.getItem(this.configKey);

        let layoutJson: {}
        if (layoutString) {
            layoutJson = JSON.parse(layoutString);
        } else {
            layoutJson = this.defaultJson;
        }

        this.state = FlexLayout.Model.fromJson(layoutJson);

        this.factory = (node: TabNode) => {
            var component = node.getComponent();

            if (component === "order-blotter") {
                return <OrderBlotter executionsController={this.executionsController} orderHistoryBlotterController={this.orderHistoryBlotterController} childOrderBlotterController={this.childOrderBlotterController} listingService={this.listingService} orderService={this.orderService} orderContext={this.orderContext} node={node} model={this.state} />;
            }
            if (component === "market-depth") {
                return <MarketDepth listingContext={this.listingContext} quoteService={this.quoteService} listingService={this.listingService} node={node} model={this.state} />;
            }
            if (component === "instrument-watch") {
                return <InstrumentListingWatch listingService={this.listingService} ticketController={this.ticketController} listingContext={this.listingContext} quoteService={this.quoteService} node={node} model={this.state} />;
            }
            if (component === "nav-bar") {
                return <Navbar />;
            }
        }

        this.onSave = this.onSave.bind(this);
        this.onCancelAllOrders = this.onCancelAllOrders.bind(this);
    }

    onSave() {
        var jsonStr = JSON.stringify(this.state!.toJson(), null, "\t");
        localStorage.setItem(this.configKey, jsonStr);
        console.log("JSON IS:" + jsonStr);
    }

    onCancelAllOrders() {
        this.questionDialogController.open("Cancel all orders?", "Cancel All Orders", (response: boolean)=> {
            var params = new CancelAllOrdersForOriginatorIdParams()
            params.setOriginatorid(Login.desk)

            this.orderMonitorClient.cancelAllOrdersForOriginatorId(params, Login.grpcContext.grpcMetaData, (err: Error,
                response: Empty) => {

                    if (err) {
                        let msg = "error whilst cancelling all orders:" + err.message
                        logError(msg)
                        alert(msg)
                      } else {
                        logDebug("cancelled all orders")
                      }

                } )
        })

        
    }


    public render() {



        let contents: React.ReactNode = "loading ...";
        if (this.state !== null) {
            contents = <FlexLayout.Layout
                ref="layout"
                model={this.state}
                factory={this.factory}
            />;
        }


        return (<div className="app" >

            <div className="toolbar" >
                <Navbar className="bp3-dark">
                    <Navbar.Group align={Alignment.LEFT}>
                        <Navbar.Heading>Open Trading Platform</Navbar.Heading>
                        <Navbar.Divider />
                        <Button className="bp3-minimal" icon="floppy-disk" text="Save Layout" onClick={this.onSave} />
                    </Navbar.Group>
                </Navbar>
            </div>
            <div>
                <OrderTicket quoteService={this.quoteService} tickerController={this.ticketController} ></OrderTicket>
                <ChildOrderBlotter childOrderBlotterController={this.childOrderBlotterController} orderService={this.orderService} listingService={this.listingService}></ChildOrderBlotter>
                <OrderHistoryBlotter orderHistoryBlotterController={this.orderHistoryBlotterController} orderService={this.orderService} listingService={this.listingService}></OrderHistoryBlotter>
                <Executions executionsController={this.executionsController} orderService={this.orderService} listingService={this.listingService}></Executions>
                <QuestionDialog controller={this.questionDialogController}></QuestionDialog>
            </div>

            <div className="contents">
                {contents}
            </div>
            <div className="toolbar" >
                <Navbar className="bp3-dark">
                    <Navbar.Group align={Alignment.LEFT}>
                        <Navbar.Heading>Status</Navbar.Heading>
                        <Navbar.Divider />
                        <Button className="bp3-minimal" icon="delete" text="Cancel All Orders" onClick={this.onCancelAllOrders} />
                    </Navbar.Group>
                </Navbar>
            </div>

        </div>);


    }

}

export class QuestionDialogController {

    private dialog?: QuestionDialog

    setDialog(dialog: QuestionDialog) {
        this.dialog = dialog
    }

    open(question : string, title: string, callback: (response: boolean)=>void) {
        if( this.dialog ) {
            this.dialog.open(question, title, callback)
        }
    }

}


export class ExecutionsController {

    private executions?: Executions;

    setView(executions: Executions) {
        this.executions = executions
    }

    open(order : Order, width:number) {
        if (this.executions) {
            this.executions.open(order,  width)
        }
    }

}

export class OrderHistoryBlotterController {

    private orderHistoryBlotter?: OrderHistoryBlotter;

    setBlotter(orderHistoryBlotter: OrderHistoryBlotter) {
        this.orderHistoryBlotter = orderHistoryBlotter
    }

    openBlotter(order : Order, config: TableViewConfig, width:number) {
        if (this.orderHistoryBlotter) {
            this.orderHistoryBlotter.open(order,  config, width)
        }
    }

}


export class ChildOrderBlotterController {

    private childOrderBlotter?: ChildOrderBlotter;

    setBlotter(childOrderBlotter: ChildOrderBlotter) {
        this.childOrderBlotter = childOrderBlotter
    }

    openBlotter(parentOrder : Order, orders: Array<Order>, 
         config: TableViewConfig, width:number) {
        if (this.childOrderBlotter) {
            this.childOrderBlotter.open(parentOrder,orders,  config, width)
        }
    }

}



export class TicketController {

    private orderTicket?: OrderTicket;

    setOrderTicket(orderTicket: OrderTicket) {
        this.orderTicket = orderTicket
    }

    openTicket(side: Side, listing: Listing) {
        if (this.orderTicket) {
            this.orderTicket.openTicket(side, listing)
        }
    }

}

export class ListingContext {

    selectedListing?: Listing

    private listeners: Array<(listing: Listing) => void>

    constructor() {
        this.listeners = new Array<(listing: Listing) => void>()

    }

    setSelectedListing(listing: Listing) {
        this.selectedListing = listing
        this.listeners.forEach(l => l(listing))
    }

    addListener(listener: (listing: Listing) => void) {
        if (this.selectedListing) {
            listener(this.selectedListing)
        }

        this.listeners.push(listener)
    }

}

export class OrderContext {

    selectedOrder?: Order
    private listeners: Array<(order: Order) => void>

    constructor() {
        this.listeners = new Array<(order: Order) => void>()
    }

    setSelectedOrder(order: Order) {
        this.selectedOrder = order
        this.listeners.forEach(l => l(order))
    }

    addListener(listener: (order: Order) => void) {
        if (this.selectedOrder) {
            listener(this.selectedOrder)
        }
        this.listeners.push(listener)
    }

}