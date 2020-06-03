import { AnchorButton, Classes, Dialog, Intent, Checkbox } from '@blueprintjs/core';
import * as React from "react";
import { logError } from '../../logging/Logging';
import { ColumnChooserController } from '../Container';
import * as table from "../TableView/TableView";



export interface ColumnChooserProps {

    controller: ColumnChooserController

}

interface ColumnChooserState {
    isOpen: boolean
    usePortal: boolean
    tableName: string
    columns: Array<column>
    callback: (newVisibleCols: JSX.Element[] | undefined, widths: number[] | undefined) => void


}



export default class ColumnChooser extends React.Component<ColumnChooserProps, ColumnChooserState>{

    controller: ColumnChooserController

    allColumns: Map<string, JSX.Element>


    constructor(props: ColumnChooserProps) {
        super(props)

        this.controller = props.controller
        this.controller.setDialog(this)
        this.allColumns = new Map()



        this.state = {
            isOpen: false,
            usePortal: false,
            columns: new Array<column>(),
            tableName: "",
            callback: (newVisibleCols: JSX.Element[] | undefined) => { }

        }

    }

    getTitle(): string {
        return "Edit " + this.state.tableName + " visible columns"
    }

   

    handleChecked(id:string,  checked:boolean) {
        let cols = this.state.columns
        for( let col of cols) {
            if( col.id === id) {
                col.visible=checked
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
            cbs.push(<Checkbox key={col.id} checked={col.visible} label={col.name}  onChange={(e:React.FormEvent<HTMLInputElement>)=>{
                let checked = e.currentTarget.checked
                this.handleChecked(col.id, checked)
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







    open(tableName: string, visibleColumns: JSX.Element[], widths: number[], allColumns: JSX.Element[], callback: (newVisibleCols: JSX.Element[] | undefined,
        widths: number[] | undefined) => void) {

        this.allColumns.clear()
        for (let col of allColumns) {
            let id = col.props["id"]
            this.allColumns.set(id, col)
        }



        let visibleColIds = new Set<string>()

        let columns = new Array<column>()

        let idx = 0
        for (let col of visibleColumns) {
            let id = col.props["id"]
            let colView = new column(true, id, col.props["name"],
                widths[idx])
            columns.push(colView)
            visibleColIds.add(id)
            idx++
        }

        for (let col of allColumns) {
            let id = col.props["id"]

            if (!visibleColIds.has(id)) {
                let colView = new column(false, id, col.props["name"],
                    table.defaultColumnWidth)
                columns.push(colView)
            }



        }


        let state = {
            ...this.state, ...{
                isOpen: true,
                columns: columns,
                callback: callback,
                tableName: tableName
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

        let result = new Array<JSX.Element>()
        let widths = new Array<number>()
        for (let col of this.state.columns) {
            if (col.visible) {
                let jsxCol = this.allColumns.get(col.id)
                if (jsxCol) {
                    result.push(jsxCol)
                    widths.push(col.width)
                } else {
                    logError("column nout found for id " + col.id)
                }
            }
        }


        this.state.callback(result, widths)
    };


    private handleCancel = () => {
        this.setState({
            ...this.state, ...{
                isOpen: false,
            }
        })
        this.state.callback(undefined, undefined)
    };





}

class column {
    visible: boolean
    id: string
    name: string
    width: number

    constructor(visible: boolean,
        id: string,
        name: string,
        width: number) {
        this.visible = visible
        this.id = id
        this.name = name
        this.width = width
    }

}