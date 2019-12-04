import { Decimal64 } from "../serverapi/common_pb"


export function toNumber(dec?: Decimal64): number | undefined {
    if (dec) {
      return dec.getMantissa() * Math.pow(10, dec.getExponent())
    }
  
    return undefined
  }
  
  export function toDecimal64(num: number): Decimal64 {
    
  
    let result = new Decimal64()
  
    let numStr : string =  num.toString()
    if( numStr.indexOf('e') >= 0) {
        throw new Error("Unable to convert number " + numStr + " to Decimal64")
    }

    let dpIdx = numStr.indexOf('.')
    if( dpIdx >= 0) {
        numStr = numStr.replace('.','')
        result.setMantissa(parseInt(numStr))
        result.setExponent(dpIdx)


    } else {
        result.setMantissa(num)
        result.setExponent(0)
    }
  
  
  
    
  
    return result
  }