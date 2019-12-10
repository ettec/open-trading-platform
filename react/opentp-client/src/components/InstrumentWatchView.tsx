import { TabNode, Model, Actions } from "flexlayout-react";
import React from 'react';

import v4 from 'uuid';
import { InstrumentWatchLine as ListingWatchLine} from '../model/Model';
import InstrumentSearchBar from "./InstrumentSearchBar";
import './OrderBlotter.css';
import { Listing } from "../serverapi/listing_pb";
import { ClientMarketDataServiceClient } from "../serverapi/CmdsServiceClientPb";
import Login from "./Login";
import { BookLine } from "../serverapi/cmds_pb";





interface InstrumentWatchState {
  watches: ListingWatchLine[]
}

interface InstrumentWatchProps {
  node : TabNode,
  model: Model,
}

interface PersistentConfig {
  instrumentIds: number[]
}



export default class InstrumentWatchView extends React.Component<InstrumentWatchProps, InstrumentWatchState> {

  // TODO: write a client side quote service to ensure single subscription per quote
  marketDataService = new ClientMarketDataServiceClient(Login.grpcContext.serviceUrl, null, null)

  watchMap: Map<number, ListingWatchLine> = new Map()

  constructor(props: InstrumentWatchProps) {
    super(props); 

    let initialState: InstrumentWatchState = {
      watches: Array.from(this.watchMap.values())
    }

    this.state = initialState;

    this.addInstrument = this.addInstrument.bind(this);

    this.props.node.setEventListener("save", (p)=> {
      let persistentConfig : PersistentConfig = {instrumentIds: Array.from(this.watchMap.keys())}
      this.props.model.doAction( Actions.updateNodeAttributes(props.node.getId(), {config:persistentConfig}))
    });

    if( this.props.node.getConfig() && this.props.node.getConfig()) {
      let persistentConfig : PersistentConfig = this.props.node.getConfig();
      persistentConfig.instrumentIds.forEach(id => {
        this.addListingLine(id)
      })

    }


  }

  addInstrument(listing?: Listing) {

    if (listing) {

      if (this.watchMap.has(listing.getId())) {
        return;
      }

      this.addListingLine(listing);
    }

  }

  private addListingLine(listing: Listing) {










    var fetchRequestString: string = 'http://192.168.1.100:31352/instrument-lookup/instrument/' + instId;
    fetch(fetchRequestString, {
      method: 'GET',
    })
      .then(response => {
        if (response.ok) {
          response.json().then(data => {
            if (data != null) {
              let instrument: Instrument = data as Instrument;
              let instrumentWatchLine: ListingWatchLine = {
                ...{
                name: instrument.canon.name,
                  symbol: instrument.symbols.IEX,
                  type: instrument.canon.type,
                  id: instrument.id
                }, ...{
                  bidPrice: "0", bidSize: "0",
                  askPrice: "0", askSize: "0"
                }
              };
              ;
              this.watchMap.set(instrumentWatchLine.id, instrumentWatchLine);
              this.setState({
                watches: Array.from(this.watchMap.values())
              });

            }
          });
        }
      }).catch(error => {
        console.log(error);
      });
  }

  public render() {
    var watches: ListingWatchLine[];
    if (this.state) {
      watches = Object.assign([], this.state.watches);
    } else {
      watches = []
    }

    const clonedWatches = watches;


    return (

      <div className="bp3-dark">

        <InstrumentSearchBar add={this.addInstrument} />



      </div>


    );
  }

}

class ListingWatchLine {
  listing?  : Listing;
  quote? : BookLine;
}