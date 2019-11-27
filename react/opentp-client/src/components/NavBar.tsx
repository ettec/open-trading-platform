import React, { PureComponent } from 'react';
import { Navbar, Alignment, Button } from '@blueprintjs/core';
import "flexlayout-react/style/dark.css";

export interface NavBarProps {
    
}
 
export interface NavBarState {
    
}
 
class NavBarComponent extends React.Component<NavBarProps, NavBarState> {

    state = { }

    constructor() {
        super({}, {}); 

        this.onSave = this.onSave.bind(this);
    }

    onSave() {
     //   var jsonStr = JSON.stringify(this.state!.toJson(), null, "\t");
     //   console.log("JSON IS:" + jsonStr);
    }

    render() { 
        return (
            
        <Navbar className="bp3-dark">
            <Navbar.Group align={Alignment.LEFT}>
                <Navbar.Heading>Open OMS</Navbar.Heading>
                <Navbar.Divider />
                <Button className="bp3-minimal" icon="home" text="Home" />
                <Button className="bp3-minimal" icon="floppy-disk" text="Save" onClick={this.onSave}/>
            </Navbar.Group>
        </Navbar>
         );
        
    }
}
 
export default NavBarComponent;