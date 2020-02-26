package model

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math/big"
)


func (m *Decimal64) AsDecimal() decimal.Decimal {
	if m == nil {
		return decimal.New(0, 0)
	}

	return decimal.New(m.Mantissa, m.Exponent)
}

func ToDecimal64(d decimal.Decimal) *Decimal64 {

	if !d.Coefficient().IsInt64() {
		panic(fmt.Sprintf("unable to convert decimal coefficient to int64: %v", d.Coefficient()))
	}

	return &Decimal64{
		Mantissa: d.Coefficient().Int64(),
		Exponent: d.Exponent(),
	}
}

func IntToDecimal64(val int64) *Decimal64 {
	return &Decimal64{
		Mantissa: val,
		Exponent: 0,
	}
}


func Compare(l Decimal64, r Decimal64) int {

	if l.Exponent > r.Exponent {
		expDiff := int64(l.Exponent-r.Exponent)
		diffMultiplier := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(expDiff), nil)
		adjustedLeft := big.NewInt(0).Mul(big.NewInt(l.Mantissa), diffMultiplier )
		return adjustedLeft.Cmp(big.NewInt(r.Mantissa))
	} else if r.Exponent > l.Exponent {
		expDiff := int64(r.Exponent-l.Exponent)
		diffMultiplier := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(expDiff), nil)
		adjustedRight := big.NewInt(0).Mul(big.NewInt(r.Mantissa), diffMultiplier )
		return big.NewInt(l.Mantissa).Cmp(adjustedRight)
	} else {
		diff := l.Mantissa - r.Mantissa
		switch  {
		case diff < 0:
			return -1
		case diff > 0:
			return 1
		default :
			return 0
		}
	}


}

