package model

import (
	"math/big"
)

func Compare(l Decimal64, r Decimal64) int {

	var expDiff = l.Exponent - r.Exponent

	if expDiff < 0 {
		diffMultiplier := int64(10^expDiff)
		adj := big.NewInt(0).Mul(big.NewInt(l.Mantissa), big.NewInt(diffMultiplier))
		return adj.Cmp(big.NewInt(r.Mantissa))
	} else if expDiff > 0 {
		diffMultiplier := int64(10^expDiff)
		adj := big.NewInt(0).Mul(big.NewInt(r.Mantissa), big.NewInt(diffMultiplier))
		return adj.Cmp(big.NewInt(l.Mantissa))
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

