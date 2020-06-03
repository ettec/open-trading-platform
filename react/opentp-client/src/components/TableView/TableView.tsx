import { IRegion } from "@blueprintjs/table";
import "@blueprintjs/table/lib/css/table.css";
import React from 'react';
import { logError } from "../../logging/Logging";
import '../TableView/TableCommon.css';
import { ColumnChooserController } from "../Container";




export interface TableViewProperties {
    colsChooser : ColumnChooserController
}


export interface TableViewState {

    columns: Array<JSX.Element>
    columnWidths: Array<number>
}


export interface TableViewConfig {
    columnOrder: string[]
    columnWidths: number[]
}


export default abstract class TableView<P extends TableViewProperties, S extends TableViewState> extends React.Component<P, S>{


    private colsChooser : ColumnChooserController

    protected abstract getTableName() : string
    protected abstract getColumns() : JSX.Element[]

    constructor( props: P) {
        super(props)
        this.colsChooser = props.colsChooser
    }


    protected editVisibleColumns = () => {

        this.colsChooser.open(this.getTableName(), this.state.columns, this.state.columnWidths,
        this.getColumns(), (newCols, newWidths)=> {
    
          if( newCols && newWidths ) {
            let newState: TableViewState = {
              ...this.state, ...{
                columns : newCols,
                columnWidths : newWidths
              }
            }
      
            this.setState(newState)
          }
          
        })
        
      }


    columnResized = (index: number, size: number) => {
        let newColWidths = this.state.columnWidths.slice();
        newColWidths[index] = size
        let tableViewState = {
            ...this.state, ...{
                columnWidths: newColWidths
            }
        }

        this.setState(tableViewState)

    }

    onColumnsReordered = (oldIndex: number, newIndex: number, length: number) => {

        let newCols = reorderColumnData(oldIndex, newIndex, length, this.state.columns)
        let newColWidths = reorderColumnData(oldIndex, newIndex, length, this.state.columnWidths)

        let tableViewState = {
            ...this.state, ...{
                columns: newCols,
                columnWidths: newColWidths
            }
        }

        this.setState(tableViewState)
    }





    getSelectedItems<T>(selectedRegions: IRegion[], items: T[]) {
        let selectedOrderArray: Array<T> = new Array<T>();
        for (let region of selectedRegions) {
            let firstRowIdx: number;
            let lastRowIdx: number;
            if (region.rows) {
                firstRowIdx = region.rows[0];
                lastRowIdx = region.rows[1];
            }
            else {
                firstRowIdx = 0;
                lastRowIdx = items.length - 1;
            }
            for (let i = firstRowIdx; i <= lastRowIdx; i++) {
                let item = items[i];
                if (item) {
                    selectedOrderArray.push(item);
                }
            }
        }
        return selectedOrderArray;
    }

    



}






export function reorderColumnData<T>(oldIndex: number, newIndex: number, length: number, cols: Array<T>): Array<T> {

    let colSegment = cols.slice(oldIndex, oldIndex + length)
    let left = cols.slice(0, oldIndex)
    let right = cols.slice(oldIndex + length, cols.length)
    let colsWithoutSeg = left.concat(right)

    let newLeft = colsWithoutSeg.slice(0, newIndex)
    let newRight = colsWithoutSeg.slice(newIndex, colsWithoutSeg.length)

    return newLeft.concat(colSegment).concat(newRight)
}

export function getColIdsInOrder(cols: JSX.Element[]) {

    let colOrderIds = new Array<string>();
    for (let col of cols) {
        let colId = col.props["id"];
        if (colId) {
            colOrderIds.push(colId);
        } else {
            logError("column is missing id field, column:" + col)
        }
    }
    return colOrderIds;
}


export const defaultColumnWidth = 100;

export function getConfiguredColumns(columns: JSX.Element[], config?: TableViewConfig) :  [Array<JSX.Element>,  Array<number>] {
    let colMap = new Map<string, JSX.Element>();
    for (let col of columns) {
        colMap.set(col.props["id"], col);
    }
    
    if (config) {
        
            let cols = new Array<JSX.Element>();
            let widths =  new Array<number>()
            let idx=0
            for (let id of config.columnOrder) {
                let col = colMap.get(id);
                if (col) {
                    cols.push(col);
                    widths.push(config.columnWidths[idx])
                }
                
                idx++
            }

            return [cols, widths]

    
    }

    return [ new Array<JSX.Element>() , new Array<number>() ];
}

