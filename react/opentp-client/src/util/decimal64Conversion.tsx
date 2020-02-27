import { Decimal64 } from "../serverapi/common_pb"


export function toNumber(dec?: Decimal64): number | undefined {

    

    if (dec) {
      if( dec.getExponent() === 0 ) {
        return dec.getMantissa()
      }

      if( dec.getExponent() > 0 ) {
        return dec.getMantissa() * Math.pow(10, dec.getExponent())
      } else {

          let sign=1
          let mantissa = dec.getMantissa()
          if(mantissa < 0 ) {
              mantissa =  -mantissa
              sign=-1
          }

          let manStr = mantissa.toString()
          if( manStr.length < -1*dec.getExponent()) {
            manStr = manStr.padStart(-1*dec.getExponent()  , "0")
          }
          let decPos = manStr.length + dec.getExponent()
          manStr = manStr.slice(0,decPos) + "." + manStr.slice(decPos, manStr.length)
          if( sign === -1 ) {
            manStr = "-" + manStr
          }

          return parseFloat(manStr)
      }

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
        result.setExponent(dpIdx - numStr.length)


    } else {
        result.setMantissa(num)
        result.setExponent(0)
    }
  
  
  
    
  
    return result
  }