import { Alignment, Button, Navbar } from "@blueprintjs/core";
import FlexLayout, { Model, TabNode } from "flexlayout-react";
import "flexlayout-react/style/dark.css";
import React from 'react';
import OrderBlotterContainer from '../containers/OrderBlotterContainer';
import InstrumentWatchView from "./InstrumentWatchView";
import MarketDepth from './MarketDepth';
import OrderTicket from './OrderTicket';


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


    state: Model;
    factory: (node: TabNode) => React.ReactNode;
    readonly configKey : string = "open-oms-config"; 


    constructor() {
        super({}, {});

        let layoutString : string | null= localStorage.getItem(this.configKey);

        let layoutJson : {}
        if( layoutString ) {
             layoutJson =  JSON.parse(layoutString) ;
        } else {
             layoutJson = this.defaultJson;
        }

        this.state = FlexLayout.Model.fromJson(layoutJson);

        this.factory = (node: TabNode) => {
            var component = node.getComponent();
            if (component === "order-ticket") {
                return <OrderTicket/>
            }
            if (component === "order-blotter") {
                return <OrderBlotterContainer />;
            }
            if (component === "market-depth") {
                return <MarketDepth />;
            }
            if (component === "instrument-watch") {
                return<InstrumentWatchView node={node} model={this.state} />;
            }
            if(component === "nav-bar") {
                return <Navbar/>;
            }

        }

        this.onSave = this.onSave.bind(this);
    }

    onSave() {
        var jsonStr = JSON.stringify(this.state!.toJson(), null, "\t");
        localStorage.setItem(this.configKey, jsonStr);
        console.log("JSON IS:" + jsonStr);
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
                <Navbar.Heading>Open OMS</Navbar.Heading>
                <Navbar.Divider />
                <Button className="bp3-minimal" icon="home" text="Home" />
                <Button className="bp3-minimal" icon="floppy-disk" text="Save" onClick={this.onSave}/>
            </Navbar.Group>
        </Navbar>
        </div>

        <div className="contents">
        {contents}
        </div>
    </div>);


    } 

}
