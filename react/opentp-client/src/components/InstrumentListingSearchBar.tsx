import {  ItemRenderer, Select } from "@blueprintjs/select";
import { MenuItem } from "@blueprintjs/core";
import React from 'react';
import v4 from 'uuid';
import './TableView/TableCommon.css';
import { Button } from "@blueprintjs/core";
import Login from "./Login";
import { Listing } from "../serverapi/listing_pb";
import { logError, logGrpcError } from "../logging/Logging";
import { StaticDataServiceClient } from "../serverapi/StaticdataserviceServiceClientPb";
import { MatchParameters, Listings } from "../serverapi/staticdataservice_pb";

const ListingSelect = Select.ofType<Listing>();

const renderListing: ItemRenderer<Listing> = (listing, { handleClick, modifiers, query }) => {
    if (!modifiers.matchesPredicate) {
        return null;
    }
    var instrument = listing.getInstrument()
    if( !instrument ) {
        logError("instrument of listing " + listing + " is not set")
        return null;
    }

    var market = listing.getMarket()
    if( !market ) {
        logError("market of listing " + listing + " is not set")
        return null;
    }

    const text = `${instrument.getDisplaysymbol()}`;
    return (
        <MenuItem
            active={modifiers.active}
            disabled={modifiers.disabled}
            label={instrument.getName() + " - " + market.getName()}
            key={listing.getId()}
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
    return text.replace(/([.*+?^=!:${}()|[\]/\\])/g, "\\$1");
}

interface ListingSearchBarState {
    items: Listing[]
    selected?: Listing;
}

interface ListingSearchBarProps {
    add : (listing? : Listing) => void
}


export default class InstrumentListingSearchBar extends React.Component<ListingSearchBarProps, ListingSearchBarState> {

    id: string;
    lastSearchString: string;

    staticDataService = new StaticDataServiceClient(Login.grpcContext.serviceUrl, null, null)
   
    constructor(props : ListingSearchBarProps) {
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

        if( query !== this.lastSearchString) {
            this.lastSearchString = query;
        } else {
            return
        }


        if( !query || query.length < 2) {
            let newState = { 
                ...this.state,...{
                }
            }
            newState.items = []
            this.setState(newState) 
            return;

        }


        let p = new MatchParameters()
        p.setSymbolmatch(query)
        this.staticDataService.getListingsMatching(p, Login.grpcContext.grpcMetaData, (err, listings : Listings) => {

            if (err) {
                logGrpcError("failed to get listings matching:", err)
              return
            }
    
            let newState = { 
                ...this.state,...{}
            }

            newState.items = listings.getListingsList()

            this.setState(newState) 

          })
        
    }


    public render() {

        const selected = this.state.selected;
        const items = this.state.items;

        return (

            <div>

                <ListingSelect 
                    
                    items={items}
                    resetOnClose={true}
                    onItemSelect={this.handleValueChange}
                    onQueryChange={this.handleQueryChange}
                    noResults={<MenuItem disabled={true} text="No results." />}
                    itemRenderer={renderListing}>
                    
                    <Button
                        rightIcon="caret-down"
                        text={selected? `${this.getSelectDisplayName(selected)} ` : "(No selection)"}
                    />
                </ListingSelect>
                <Button onClick={()=>this.props.add(this.state.selected)}>Add</Button>
            </div>
        );
    }


    private getSelectDisplayName(listing: Listing): string {
        let i = listing.getInstrument()
        let m = listing.getMarket()
        
        if( i && m ) {
            return i.getDisplaysymbol() + " - " + i.getName() + " - " + m.getMic()
        } else {
            let msg = "listing " + listing + " is missing instrument or market"
            logError(msg)
            return msg
        }
    }

    private handleValueChange = (listing: Listing) => {
        this.setState({
             ...this.state,...{
                selected: listing,
            }
        })
    };

}

