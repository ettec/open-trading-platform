import { Button, InputGroup } from "@blueprintjs/core";
import React from 'react';
import ReactTable from 'react-table';
import "react-table/react-table.css";
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

const marketDataService = new ClientMarketDataServiceClient('http://192.168.1.100:32365', null, null)



export default class MarketDepth extends React.Component<{}, MarketDepthState > {

    stream?: grpcWeb.ClientReadableStream<Book>;    
   // marketDepthSource : EventSource;
    id : string;
    
    constructor() {
        super({});
        this.setState({symbol:''});

        this.id = v4();

        var subscription = new SubscribeRequest()
        subscription.setSubscriberid(this.id)

        this.stream = marketDataService.subscribe(subscription, Login.grpcContext.grpcMetaData)

        this.stream.on('data', (response: Book) => {
          console.log('Received book' + response)
          let blotterState : MarketDepthState =  {...this.state,... {
            book: response,
          }}

          this.setState( blotterState );
        });
        this.stream.on('status', (status: grpcWeb.Status) => {
          if (status.metadata) {
            console.log('Received metadata');
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
        marketDataService.addSubscription(subscription, Login.grpcContext.grpcMetaData, (err, response )=>{
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
                <ReactTable<LocalBookLine> 
                  
                  data={clonedDepth}
                  columns={[
                    {
                      columns: [
                        {
                          Header: "Bid Size",
                          accessor: "bidSize"
                        },
                        {
                          Header: "Bid Px",
                          accessor: "bidPrice"
                        },
                        {
                          Header: "Ask Px",
                          accessor: "askPrice"
                        },
                        {
                          Header: "Ask Size",
                          accessor: "askSize"
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
                  
                />)
                <br />



              </div>


        );
    }

}