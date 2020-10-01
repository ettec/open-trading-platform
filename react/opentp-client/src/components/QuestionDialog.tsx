import { AnchorButton, Classes, Dialog, Intent } from '@blueprintjs/core';
import * as React from "react";
import { QuestionDialogController } from './Container/Container';



export interface QuestionDialogProps {
    controller:  QuestionDialogController
}


interface QuestionDialogState {
   isOpen: boolean
   usePortal: boolean
   question: string
   title: string
   callback: (response: boolean)=>void
}


export default class QuestionDialog extends React.Component<QuestionDialogProps, QuestionDialogState> {


    constructor(props: QuestionDialogProps) {
        super(props)

        props.controller.setDialog(this)

        this.state = {
            isOpen: false,
            usePortal: false,
            question: "",
            title: "",
            callback: (response) => {},
        }

    }

   

    render() {
        return (
            <Dialog
                icon="help"
                onClose={this.handleCancel}
              //  title={this.state.title}
                {...this.state}
                className="bp3-dark">
                <div className={Classes.DIALOG_BODY} >
                  <h1>{this.state.question}</h1>
                </div>
                <div className={Classes.DIALOG_FOOTER}>
                    <div className={Classes.DIALOG_FOOTER_ACTIONS}>
                    <AnchorButton onClick={this.handleOK}
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

   

   

    open(question : string, title : string, callback: (response: boolean)=>void) {



        let state =
        {
            isOpen: true,
            usePortal: false,
            question: question,
            title: title,
            callback: callback
        }

        this.setState(state)
    }

    private handleOK = () => {
        this.setState({
            ...this.state, ...{
                isOpen: false,
            }
        })
        this.state.callback(true)
    };


    private handleCancel = () => {
        this.setState({
            ...this.state, ...{
                isOpen: false,
            }
        })
        this.state.callback(true)
    };

}
