import { IRegion } from "@blueprintjs/table";
import "@blueprintjs/table/lib/css/table.css";
import React from 'react';
import { logError } from "../../logging/Logging";
import '../TableView/TableCommon.css';


export interface TableViewState {

    columns: Array<JSX.Element>
    columnWidths: Array<number>
}

export interface TableViewConfig {
    columnOrder: string[]
    columnWidths: number[]
}


export default class TableView<P, S extends TableViewState> extends React.Component<P, S>{


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
            logError("columnd is missing id field, column:" + col)
        }
    }
    return colOrderIds;
}

export function getConfiguredColumns(columns: JSX.Element[], config?: TableViewConfig) :  [Array<JSX.Element>,  Array<number>] {
    let colMap = new Map<string, JSX.Element>();
    for (let col of columns) {
        colMap.set(col.props["id"], col);
    }
    let defaultCols = Array.from(colMap.values());
    let defaultColWidths = new Array<number>();
    for (let i: number = 0; i < defaultCols.length; i++) {
        defaultColWidths.push(100);
    }
    if (config) {
        let pc: TableViewConfig = config;
        if (pc.columnWidths && pc.columnWidths.length > 0) {
            defaultColWidths = pc.columnWidths;
        }
        if (pc.columnOrder && pc.columnOrder.length > 0) {
            let cols = new Array<JSX.Element>();
            for (let id of pc.columnOrder) {
                let col = colMap.get(id);
                if (col) {
                    cols.push(col);
                }
                colMap.delete(id)
            }

            for( let newCol of colMap.values() ) {
                cols.push(newCol);
            }

            defaultCols = cols;
        }
        while (defaultColWidths.length < defaultCols.length) {
            defaultColWidths.push(100);
        }
    }
    return [ defaultCols, defaultColWidths ];
}

