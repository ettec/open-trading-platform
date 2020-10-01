import { toDecimal64, toNumber } from "./decimal64Conversion"


test("convert numbers to decimal64", () => {

    let toTest: Array<number> = [6.45334, 0.00056, 10021.43523, 5, 34, 4000.345430, 0.0000034, 230000, 0, -1, -1000.3254,
        -0.00006309]

    for( var num of toTest) {
        let result = toDecimal64(num)

        let resultAsNumber = toNumber(result);
 
        expect(num).toEqual(resultAsNumber)
    }

})
