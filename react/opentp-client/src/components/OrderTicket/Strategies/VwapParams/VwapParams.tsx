import * as React from "react";
import { Label, NumericInput } from "@blueprintjs/core";
import { TimePicker, TimePrecision } from "@blueprintjs/datetime";

export interface Props {
    children?: React.ReactNode
}

export interface State {
}

export interface VwapParamsState {
    autoFocus: boolean;
    precision?: TimePrecision;
    selectAllOnFocus?: boolean;
    showArrowButtons?: boolean;
    disabled?: boolean;
    minTime?: Date;
    maxTime?: Date;
    useAmPm?: boolean;
}


export default class VwapParams extends React.Component<Props, State> {

    public state = {
        autoFocus: true,
        disabled: false,
        precision: TimePrecision.MINUTE,
        selectAllOnFocus: false,
        showArrowButtons: false,
        useAmPm: false,
    };

    constructor(props: Props) {
        super(props)

        this.state = {
            autoFocus: true,
            disabled: false,
            precision: TimePrecision.MINUTE,
            selectAllOnFocus: false,
            showArrowButtons: false,
            useAmPm: false,
        };
    }

    render() {
        return (
            <div  style={{display: 'flex',  flexDirection: 'row', paddingTop:25}}>
                
                <Label htmlFor="input-b" style={{paddingRight:5}}>Start</Label>
                <TimePicker  precision={TimePrecision.MINUTE} ></TimePicker>
                <Label htmlFor="input-b" style={{paddingRight:5, paddingLeft:25}}>  End </Label>
                <TimePicker {...this.state}></TimePicker>
                <Label htmlFor="input-b" style={{paddingRight:5, paddingLeft:25}}>Buckets</Label>
                <NumericInput min={1} style={{maxWidth: 60}}/>
            </div>
        )
    }
}
