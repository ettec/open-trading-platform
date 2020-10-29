import { AnchorButton, Checkbox, Classes, Dialog, Intent } from '@blueprintjs/core';
import { Column, ColumnState } from 'ag-grid-community';
import * as React from "react";
import { AgGridColumnChooserController } from "../Container/Controllers";



export interface ColumnChooserAgGridProps {

    controller: AgGridColumnChooserController

}

interface ColumnChooserAgGridState {
    isOpen: boolean
    usePortal: boolean
    tableName: string
    columns: ColumnState[]
    callback: (columns: ColumnState[] | undefined) => void
    idToHeader: Map<string, string>

}



export default class ColumnChooserAgGrid extends React.Component<ColumnChooserAgGridProps, ColumnChooserAgGridState>{

    controller: AgGridColumnChooserController


    constructor(props: ColumnChooserAgGridProps) {
        super(props)

        this.controller = props.controller
        this.controller.setDialog(this)
        

        this.state = {
            isOpen: false,
            usePortal: false,
            columns: new Array<ColumnState>(),
            tableName: "",
            callback: (columns: ColumnState[] | undefined)  => { },
            idToHeader: new Map<string,string>()
        }

    }

    getTitle(): string {
        return "Edit " + this.state.tableName + " visible columns"
    }

   

    handleChecked(id:string | undefined,  checked:boolean) {
        let cols = this.state.columns
        for( let col of cols) {
            if( col.colId === id) {
                col.hide=!checked
            }
        }

        let state = {
            ...this.state, ...{
                columns: cols,
            }
        }

        this.setState(state)

    }

    render() {

        const cbs =[]
        for( let col of this.state.columns) {
            cbs.push(<Checkbox key={col.colId} checked={!col.hide} label={col.colId?this.state.idToHeader.get(col.colId):col.colId}  onChange={(e:React.FormEvent<HTMLInputElement>)=>{
                let checked = e.currentTarget.checked
                this.handleChecked(col.colId, checked)
            }} />)
        }

        return (  
            <Dialog
                icon="bring-data"
                onClose={this.handleCancel}
                title={this.getTitle()}
                {...this.state}
                className="bp3-dark">
                <div className={Classes.DIALOG_BODY} >
                    {cbs}
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







    open(tableName: string, colStates: ColumnState[], cols: Column[], callback: (columns: ColumnState[] | undefined) => void) {


        let idToHeader = new Map<string,string>()
        for( let col of cols) {
            col.getColId()
            let header = col.getDefinition().headerName
            if( header) {
                idToHeader.set(col.getColId(), header)
            }
        }

        
        let state = {
            ...this.state, ...{
                isOpen: true,
                columns: colStates,
                callback: callback,
                tableName: tableName,
                idToHeader: idToHeader
            }
        }

        this.setState(state)
    }


    private handleOK = () => {
        this.setState({
            ...this.state, ...{
                isOpen: false,
            }
        })



        this.state.callback(this.state.columns)
    };


    private handleCancel = () => {
        this.setState({
            ...this.state, ...{
                isOpen: false,
            }
        })
        this.state.callback(undefined)
    };





}
