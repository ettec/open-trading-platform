import { IItemRendererProps, ItemRenderer, Select, ItemListPredicate } from "@blueprintjs/select";
import { MenuItem } from "@blueprintjs/core";
import React, { Component, RefObject } from 'react';
import "react-table/react-table.css";
import v4 from 'uuid';
import { SearchDisplayInstrument } from '../model/Model';
import './OrderBlotter.css';
import { Button } from "@blueprintjs/core";


const InstrumentSelect = Select.ofType<SearchDisplayInstrument>();



const renderInstrument: ItemRenderer<SearchDisplayInstrument> = (inst, { handleClick, modifiers, query }) => {
    if (!modifiers.matchesPredicate) {
        return null;
    }
    const text = `${inst.symbol}`;
    return (
        <MenuItem
            active={modifiers.active}
            disabled={modifiers.disabled}
            label={inst.name}
            key={inst.id}
            onClick={handleClick}
            text={highlightText(text, query)}
        />
    );
};




function highlightText(text: string, query: string) {
    let lastIndex = 0;
    const words = query
        .split(/\s+/)
        .filter(word => word.length > 0)
        .map(escapeRegExpChars);
    if (words.length === 0) {
        return [text];
    }
    const regexp = new RegExp(words.join("|"), "gi");
    const tokens: React.ReactNode[] = [];
    while (true) {
        const match = regexp.exec(text);
        if (!match) {
            break;
        }
        const length = match[0].length;
        const before = text.slice(lastIndex, regexp.lastIndex - length);
        if (before.length > 0) {
            tokens.push(before);
        }
        lastIndex = regexp.lastIndex;
        tokens.push(<strong key={lastIndex}>{match[0]}</strong>);
    }
    const rest = text.slice(lastIndex);
    if (rest.length > 0) {
        tokens.push(rest);
    }
    return tokens; 
}

function escapeRegExpChars(text: string) {
    return text.replace(/([.*+?^=!:${}()|\[\]\/\\])/g, "\\$1");
}

interface InstrumentSearchBarState {
    items: SearchDisplayInstrument[]
    selected?: SearchDisplayInstrument;
}

interface InstrumentSearchBarProps {
    add : (instrument? : SearchDisplayInstrument) => void
}


export default class InstrumentSearchBar extends React.Component<InstrumentSearchBarProps, InstrumentSearchBarState> {

    id: string;
    lastSearchString: string;
   


    constructor(props : InstrumentSearchBarProps) {
        super(props);

        let initialState = {
            items: []
        };

        this.state = initialState;

        this.id = v4();
        this.lastSearchString = ''

        
        this.handleQueryChange = this.handleQueryChange.bind(this);
       


    }

    handleQueryChange(query : string) {

        if( query != this.lastSearchString) {
            this.lastSearchString = query;
        } else {
            return
        }


        if( !query || query.length < 2) {
            let newState = { 
                ...this.state, ... {
                }
            }
            newState.items = []
            this.setState(newState) 
            return;

        }

        
        

        var fetchRequestString : string =  'http://192.168.1.100:31352/instrument-lookup/instruments-matching?searchString=' + query
        fetch(fetchRequestString, {
            method: 'GET',
        }) 
            .then(

                response => {

                    if( response.ok ) {
                        response.json().then(data => {

                            if( data != null ) {
                                let instrumentsFromJson : SearchDisplayInstrument[]  = data as SearchDisplayInstrument[];

                                let newState = { 
                                    ...this.state, ... {}
                                }

                                newState.items = instrumentsFromJson

                                this.setState(newState) 

                            } else {
                                let newState = { 
                                    ...this.state, ... {
                                    }
                                }

                                newState.items = []

                                this.setState(newState) 
                            }


                            
                        })
                    }
                }
            ).catch(error => {
                console.log(error)
            });
    }


    public render() {

        const selected = this.state.selected;
        const items = this.state.items;

        return (

            <div>

                <InstrumentSelect 
                    
                    items={items}
                    resetOnClose={true}
                    onItemSelect={this.handleValueChange}
                    onQueryChange={this.handleQueryChange}
                    noResults={<MenuItem disabled={true} text="No results." />}
                    itemRenderer={renderInstrument}>
                    
                    <Button
                        rightIcon="caret-down"
                        text={selected ? `${selected.symbol + " - " + selected.name} ` : "(No selection)"}
                    />
                </InstrumentSelect>
                <Button onClick={()=>this.props.add(this.state.selected)}>Add</Button>
            </div>
        );
    }


    private handleValueChange = (instrument: SearchDisplayInstrument) => {
        console.log("Selected: " + instrument.symbol)
        this.setState({
             ...this.state, ... {
                selected: instrument,
            }
        })
    };

}

