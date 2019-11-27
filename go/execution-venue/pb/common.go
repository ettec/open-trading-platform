package pb

import (
	"fmt"
	"github.com/shopspring/decimal"
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
