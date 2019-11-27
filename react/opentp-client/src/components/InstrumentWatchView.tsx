import { TabNode, Model, Actions } from "flexlayout-react";
import React from 'react';
import ReactTable from 'react-table';
import "react-table/react-table.css";
import v4 from 'uuid';
import { InstrumentWatchLine, SearchDisplayInstrument, Instrument } from '../model/Model';
import InstrumentSearchBar from "./InstrumentSearchBar";
import './OrderBlotter.css';





interface InstrumentWatchState {
  watches: InstrumentWatchLine[]
}

interface InstrumentWatchProps {
  node : TabNode,
  model: Model,
}

interface PersistentConfig {
  instrumentIds: number[]
}



export default class InstrumentWatchView extends React.Component<InstrumentWatchProps, InstrumentWatchState> {


  watchMap: Map<number, InstrumentWatchLine> = new Map()

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
        this.addInstrumentId(id)
      })

    }


  }

  addInstrument(instrument?: SearchDisplayInstrument) {

    if (instrument) {

      var instId = instrument.id
      if (this.watchMap.has(instId)) {
        return;
      }

      this.addInstrumentId(instId);
    }

  }

  private addInstrumentId(instId: number) {
    var fetchRequestString: string = 'http://192.168.1.100:31352/instrument-lookup/instrument/' + instId;
    fetch(fetchRequestString, {
      method: 'GET',
    })
      .then(response => {
        if (response.ok) {
          response.json().then(data => {
            if (data != null) {
              let instrument: Instrument = data as Instrument;
              let instrumentWatchLine: InstrumentWatchLine = {
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
    var watches: InstrumentWatchLine[];
    if (this.state) {
      watches = Object.assign([], this.state.watches);
    } else {
      watches = []
    }

    const clonedWatches = watches;


    return (

      <div className="bp3-dark">

        <InstrumentSearchBar add={this.addInstrument} />

        <ReactTable<InstrumentWatchLine>

          data={clonedWatches}
          columns={[
            {
              columns: [
                {
                  Header: "Name",
                  accessor: "name"
                },
                {
                  Header: "Symbol",
                  accessor: "symbol"
                },
                {
                  Header: "Type",
                  accessor: "type"
                },
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
            height: 20 * 41 + "px" // This will force the table body to overflow and scroll, since there is not enough room
          }}
          className="-striped -highlight"

        />)
                <br />



      </div>


    );
  }

}