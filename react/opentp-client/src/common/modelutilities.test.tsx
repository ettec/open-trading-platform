import { Listing, TickSizeEntry, TickSizeTable } from "../serverapi/listing_pb"
import { toDecimal64 } from "./decimal64Conversion"
import { roundToTick } from "./modelutilities"


test("roundToTick", () => {

    let tickSizeTable: TickSizeTable = new TickSizeTable()
    let entry1: TickSizeEntry = new TickSizeEntry()
    entry1.setLowerpricebound(toDecimal64(0))
    entry1.setUpperpricebound(toDecimal64(9.99))
    entry1.setTicksize(toDecimal64(0.01))



    let entry2: TickSizeEntry = new TickSizeEntry()
    entry2.setLowerpricebound(toDecimal64(10))
    entry2.setUpperpricebound(toDecimal64(99.9))
    entry2.setTicksize(toDecimal64(0.1))

    let entry3: TickSizeEntry = new TickSizeEntry()
    entry3.setLowerpricebound(toDecimal64(100))
    entry3.setUpperpricebound(toDecimal64(999))
    entry3.setTicksize(toDecimal64(1))

    tickSizeTable.addEntries(entry1)
    tickSizeTable.addEntries(entry2)
    tickSizeTable.addEntries(entry3)

    let listing: Listing = new Listing()
    listing.setTicksize(tickSizeTable)


    let toTest: Array<number> = [5.0111, 11.51666, 11.37666, 11.77666, 234.897]

    let expected: Array<number> = [5.01, 11.5, 11.4, 11.8, 235]

    for (var i = 0; i < toTest.length; i++) {
        var result = roundToTick(toTest[i], listing)
        expect(result).toEqual(expected[i])
    }

})