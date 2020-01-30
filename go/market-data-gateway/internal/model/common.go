package model

import (
	"math/big"
)


// Decimal64{Mantissa: 1, Exponent:0}, Decimal64{Mantissa:1, Exponent:1}},

func Compare(l Decimal64, r Decimal64) int {

	var expDiff = int64(l.Exponent - r.Exponent)

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


	if expDiff < 0 {
		diffMultiplier := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(-expDiff), nil)
		adjustedRight := big.NewInt(0).Mul(big.NewInt(r.Mantissa), diffMultiplier )
		return adjustedRight.Cmp(big.NewInt(l.Mantissa))
	} else if expDiff > 0 {
		diffMultiplier := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(expDiff), nil)
		adjustedLeft := big.NewInt(0).Mul(big.NewInt(l.Mantissa), diffMultiplier )
		return adjustedLeft.Cmp(big.NewInt(r.Mantissa))
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

