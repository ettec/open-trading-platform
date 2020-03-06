package model

import (
	"fmt"
	"github.com/shopspring/decimal"

)




func(t *Timestamp) After(o *Timestamp) bool {
	if t.Seconds > o.Seconds {
		return true
	}

	if t.Seconds == o.Seconds && t.Nanoseconds > o.Nanoseconds {
		return true
	}

	return false
}

func(t *Timestamp) Before(o *Timestamp) bool {
	if t.Seconds < o.Seconds {
		return true
	}

	if t.Seconds == o.Seconds && t.Nanoseconds < o.Nanoseconds {
		return true
	}

	return false
}




func (m *Decimal64) AsDecimal() decimal.Decimal {
	if m == nil {
		return decimal.New(0, 0)
	}

	return decimal.New(m.Mantissa, m.Exponent)
}

func IasD(i int) *Decimal64  {
	return &Decimal64{
		Mantissa:             int64(i),
		Exponent:             0,
	}
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

// Todo implement more efficient version of these operations that does not require interim conversion to Decimal/Rat
// (though allocation of the interim type will be on the stack so cost is minimal)

func NewFromFloat(val float64) *Decimal64 {
	return ToDecimal64(decimal.NewFromFloat(val))
}


func (m *Decimal64) Add(o *Decimal64) {
	r := ToDecimal64(m.AsDecimal().Add(o.AsDecimal()))

	m.Mantissa = r.Mantissa
	m.Exponent = r.Exponent
}

func (m *Decimal64) Equal(o *Decimal64) bool {
	return m.AsDecimal().Equal(o.AsDecimal())
}

func (m *Decimal64) LessThan(o *Decimal64) bool {
	return m.AsDecimal().LessThan(o.AsDecimal())
}

func (m *Decimal64) GreaterThan(o *Decimal64) bool {
	return m.AsDecimal().GreaterThan(o.AsDecimal())
}

