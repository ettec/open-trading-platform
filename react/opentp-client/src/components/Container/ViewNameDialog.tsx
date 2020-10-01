import { AnchorButton, Classes, Dialog, Intent, InputGroup } from '@blueprintjs/core';
import * as React from "react";
import { ViewNameDialogController } from "./Controllers";
import { Layout } from 'flexlayout-react';



export interface ViewNameDialogProps {
    controller:  ViewNameDialogController
    }


interface ViewNameDialogState {
   isOpen: boolean
   usePortal: boolean
   question: string
   title: string
   viewName: string
}


export default class ViewNameDialog extends React.Component<ViewNameDialogProps, ViewNameDialogState> {


    layout?: Layout
    component: string
    displayName: string

    constructor(props: ViewNameDialogProps) {
        super(props)

        this.open = this.open.bind(this);

        this.layout = undefined
        this.component = ""
        this.displayName = ""

        props.controller.setDialog(this)

        this.state = {
            isOpen: false,
            usePortal: false,
            question: "",
            title: "",
            viewName: "",
        }

        this.onViewNameChange = this.onViewNameChange.bind(this);
        this.handleOK = this.handleOK.bind(this);
    }

   

    render() {
        return (
            <Dialog
                icon="help"
                onClose={this.handleCancel}
             //   title={this.state.title}
                {...this.state}
                className="bp3-dark">
                <div className={Classes.DIALOG_BODY} >
                  <h1>{this.state.question}</h1>
                  <InputGroup style={{marginBottom: 30}}  placeholder="View name..." onChange={this.onViewNameChange} />
                </div>
                <div className={Classes.DIALOG_FOOTER}>
                    <div className={Classes.DIALOG_FOOTER_ACTIONS}>
                    <AnchorButton onClick={this.handleOK} disabled={this.state.viewName.length===0}
                            intent={Intent.PRIMARY}>OK
                        </AnchorButton>
                        <AnchorButton onClick={this.handleCancel}
                            intent={Intent.PRIMARY}>Cancel
                        </AnchorButton>
                    </div>
                </div>


            </Dialog>


        )
    }

    onViewNameChange(e:any) {
        if( e.target && e.target.value) {
            this.setState({
                ...this.state, ...{
                    viewName: e.target.value,
                }
            })
          }
      }
   

    open(component : string, displayName : string, layout : Layout) {


        this.component = component
        this.displayName = displayName
        this.layout = layout
        let title = "Enter name for " + displayName
        let state =
        {
            isOpen: true,
            usePortal: false,
            question: title,
            title: title,
        }

        this.setState(state)
    }

    private handleOK()  {
        this.setState({
            ...this.state, ...{
                isOpen: false,
            }
        })
        
        if( this.layout ) {
            this.layout.addTabWithDragAndDrop("Add " + this.displayName + "<br>(Drag to location)", {
                component: this.component,
                name: this.state.viewName
            });
        }

    };


    private handleCancel = () => {
        this.setState({
            ...this.state, ...{
                isOpen: false,
            }
        })
    };

}
