import TableView, { getColIdsInOrder, reorderColumnData, getConfiguredColumns, TableViewConfig } from "./TableView"
import React from 'react';
import { Column } from "@blueprintjs/table";


test("test get cold ids in order", () => {


    let cols = [<Column key="id" id="id" name="Id" />,
    <Column key="side" id="side" name="Side" />,
    <Column key="symbol" id="symbol" name="Symbol" />,]

    let colIds = getColIdsInOrder(cols)

    expect(colIds[0]).toEqual("id")
    expect(colIds[1]).toEqual("side")
    expect(colIds[2]).toEqual("symbol")
})

test("reorder column data", () => {


    let cols = [<Column key="id" id="id" name="Id" />,
    <Column key="side" id="side" name="Side" />,
    <Column key="symbol" id="symbol" name="Symbol" />,]

    let newCols = reorderColumnData(1, 0, 1, cols)

    expect(newCols[0].props["id"]).toEqual("side")
    expect(newCols[1].props["id"]).toEqual("id")
    expect(newCols[2].props["id"]).toEqual("symbol")
})

class TestConfig implements TableViewConfig {
    columnOrder: string[]
    columnWidths: number[]

    constructor(
        columnOrder: string[],
        columnWidths: number[]) {

            this.columnOrder = columnOrder
            this.columnWidths = columnWidths
        }
}