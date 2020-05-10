import { logError } from "../../logging/Logging";


export default interface TableViewConfig {
    columnOrder: string[]
    columnWidths: number[]
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

export function getConfiguredColumns(columns: JSX.Element[], config: TableViewConfig) :  [Array<JSX.Element>,  Array<number>] {
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