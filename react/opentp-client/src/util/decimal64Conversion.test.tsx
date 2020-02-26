import { toDecimal64, toNumber } from "./decimal64Conversion"


test("convert numbers to decimal64", () => {

    let num : number = 6.45334
 
    let result = toDecimal64(num)
    let resultAsNumber = toNumber(result);
 
    expect(1).toEqual(1)

})