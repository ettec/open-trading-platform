package pb

import (
	"github.com/shopspring/decimal"
	"math"
	"testing"
)

func TestToDecimal64(t *testing.T) {

	result := ToDecimal64( decimal.New(10,2))
	resultAsInt64 := int64(float64(result.Mantissa) * (math.Pow(10 , float64(result.Exponent))))


	if resultAsInt64 != 1000 {
		t.Fatalf("expected %v got %v", 1000, resultAsInt64)
	}


}

func TestToDecimal64IrrationalNumber(t *testing.T) {

	one := decimal.New(1,0)
	three := decimal.New(3,0)
	irn := one.Div(three)

	result := ToDecimal64(irn)

	resultAsFloat64 := float64(result.Mantissa) * (math.Pow(10 , float64(result.Exponent)))



	expectedResult := 0.33333333333333333333
	if !equal(resultAsFloat64,expectedResult, 0.00000001) {
		t.Fatalf("expected %v got %v", expectedResult, resultAsFloat64)
	}


}



func equal(a, b, delta float64) bool {
	return math.Abs(a - b) <= delta
}