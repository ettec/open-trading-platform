import { Button, InputGroup } from "@blueprintjs/core";
import React from 'react';


import v4 from 'uuid';
import './OrderBlotter.css';
import * as grpcWeb from 'grpc-web'
import {ClientMarketDataServiceClient} from '../serverapi/CmdsServiceClientPb'
import {SubscribeRequest, Book, BookLine, Subscription} from '../serverapi/cmds_pb'
import {LocalBookLine} from '../model/Model'
import { GrpcConsumer } from "./GrpcContextProvider/GrpcContextProvider";
import Login from "./Login";



interface MarketDepthState {
  symbol?: string,
  book: Book
}

export default class MarketDepth extends React.Component<{}, MarketDepthState > {

    stream?: grpcWeb.ClientReadableStream<Book>;    
   // marketDepthSource : EventSource;
    id : string;

    marketDataService = new ClientMarketDataServiceClient(Login.grpcContext.serviceUrl, null, null)
    
    constructor() {
        super({});
        this.setState({symbol:''});

        this.id = v4();

        var subscription = new SubscribeRequest()
        subscription.setSubscriberid(this.id)

        this.stream = this.marketDataService.subscribe(subscription, Login.grpcContext.grpcMetaData)

        this.stream.on('data', (response: Book) => {
          let blotterState : MarketDepthState =  {...this.state,... {
            book: response,
          }}

          this.setState( blotterState );
        });
        this.stream.on('status', (status: grpcWeb.Status) => {
          if (status.metadata) {
            console.log(status.metadata);
          }
        });
        this.stream.on('error', (err: grpcWeb.Error) => {
          console.log('Received error')
        });
        this.stream.on('end', () => {
          console.log('stream end signal received');
        });

        this.handleSymbolChange = this.handleSymbolChange.bind(this);
        this.onSubscribe = this.onSubscribe.bind(this);

    }


    handleSymbolChange(e:any) {


      if( e.target && e.target.value) {

        let newSymbol : string = e.target.value;
        this.setState(state=>({symbol:newSymbol}))  
      }


    }


    
    onSubscribe() {

      console.log("Subscribe to:" + this.state.symbol)

      if( this.state.symbol != null) {
        var subscription = new Subscription()
        subscription.setSubscriberid(this.id)
        subscription.setListingid(this.state.symbol)
        console.log("adding subscription:" + subscription)
        this.marketDataService.addSubscription(subscription, Login.grpcContext.grpcMetaData, (err, response )=>{
            if( err ) {
              console.log("failed to add subscription:" + err)
              return
            }

            if( response ) {
              console.log("Add subscription response:" + response)
            }
            

        } )
      }

    }
    



    public render() {

        
      
        var depth:LocalBookLine[] = [];
        if( this.state && this.state.book ) {
            for( let line of this.state.book.getDepthList()) {
                depth.push({bidSize:line.getBidsize(),bidPrice:line.getBidprice(),
                  askSize:line.getAsksize(),askPrice:line.getAskprice()})

            }
        } 

        const clonedDepth = depth;

        

        return ( 
                
              <div className="bp3-dark">
                <InputGroup onChange={this.handleSymbolChange} />
                <Button onClick={this.onSubscribe } >Subscribe</Button>
            



              </div>


        );
    }

}