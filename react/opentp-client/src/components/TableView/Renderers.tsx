import { Colors, Icon } from "@blueprintjs/core";
import React, { Component } from 'react';
import { GlobalColours } from "../Container/Colours";
import ReactCountryFlag from "react-country-flag";


export class DirectionalPriceRenderer extends Component<any, any> {
  constructor(props: any) {
    super(props);

    this.state = {
      value: this.props.value,
    };
  }


  render() {

    let price = this.state?.value?.price
    let direction = this.state?.value?.direction

    if (price) {
      if (direction) {
        if (direction > 0) {
          return <span><Icon icon="arrow-up" style={{ color: GlobalColours.UPTICK }} />{price}</span>;
        }

        if (direction < 0) {
          return <span><Icon icon="arrow-down" style={{ color: GlobalColours.DOWNTICK }} />{price}</span>;
        }

        return <span>{price}</span>;
      } else {
        return <span>{price}</span>;
      }
    } else {
      return <span></span>
    }
  }
}


export class CountryFlagRenderer extends Component<any, any> {
    constructor(props: any) {
      super(props);
  
      this.state = {
        value: this.props.value,
      };
    }
  
  
    render() {
      return <span><ReactCountryFlag countryCode={this.state.value} /> {this.state.value}</span>;
    }
  }
  
  
  
  export class TargetStatusRenderer extends Component<any, any> {
    constructor(props: any) {
      super(props);
  
      this.state = {
        value: this.props.value,
      };
    }
  
  
    render() {
      if (this.state.value === "None") {
        return <span  >{this.state.value}</span>;
      } else {
        return <span style={{ color: Colors.ORANGE3 }}  ><b>{this.state.value}</b></span>;
      }
    }
  }
  
  
