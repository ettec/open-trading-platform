package model

import (
	"fmt"
	"math"
)

func (m *Listing) RoundToLotSize(qty float64) *Decimal64 {
	return roundToDecimal(qty, &Decimal64{Mantissa: 1})
}

func (m *Listing) RoundToNearestTick(price float64) (*Decimal64, error) {

	for _, entry := range m.TickSize.Entries {
		delta := entry.TickSize.ToFloat() / 1000
		lowerBound := entry.LowerPriceBound.ToFloat()
		upperBound := entry.UpperPriceBound.ToFloat()

		if compare(price, lowerBound, delta) >= 0 &&
			compare(price, upperBound, delta) <= 0 {
			return roundToDecimal(price, entry.TickSize), nil
		}
	}

	return nil, fmt.Errorf("no tick table entry for price:%v", price)
}

func compare(f1 float64, f2 float64, delta float64) int {

	diff := f1 - f2

	if math.Abs(diff) < delta {
		return 0
	}

	if diff < 0 {
		return -1
	}

	return 1
}

func roundToDecimal(price float64, tick *Decimal64) *Decimal64 {

	floatTicks := math.Round(price / tick.ToFloat())
	numTicks := int64(floatTicks)

	return &Decimal64{Mantissa: numTicks * tick.Mantissa, Exponent: tick.Exponent}
}
