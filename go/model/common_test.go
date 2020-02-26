package model

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


func TestCompare(t *testing.T) {
	type args struct {
		l Decimal64
		r Decimal64
	}

	tests := []struct {
		name string
		args args
		want int
	}{

		{"l==r",
			args{Decimal64{Mantissa: 0, Exponent:0}, Decimal64{Mantissa:0, Exponent:0}},
			0},
		{"l>r",
			args{Decimal64{Mantissa: 1, Exponent:0}, Decimal64{Mantissa:0, Exponent:0}},
			1},
		{"l<r",
			args{Decimal64{Mantissa: 0, Exponent:0}, Decimal64{Mantissa:1, Exponent:0}},
			-1},
		{"lr",
			args{Decimal64{Mantissa: 1, Exponent:1}, Decimal64{Mantissa:1, Exponent:1}},
			0},
		{"lr",
			args{Decimal64{Mantissa: 1, Exponent:1}, Decimal64{Mantissa:1, Exponent:0}},
			1},

		{"lr",
			args{Decimal64{Mantissa: 1, Exponent:0}, Decimal64{Mantissa:1, Exponent:1}},
			-1},

		{"lr",
			args{Decimal64{Mantissa: 12345600, Exponent:0}, Decimal64{Mantissa:123456, Exponent:2}},
			0},
		{"lr",
			args{Decimal64{Mantissa: 12345600, Exponent:0}, Decimal64{Mantissa:123456, Exponent:0}},
			1},
		{"lr",
			args{Decimal64{Mantissa: 12345600, Exponent:0}, Decimal64{Mantissa:123456, Exponent:3}},
			-1},

		{"lr",
			args{Decimal64{Mantissa: 12345600, Exponent:-2}, Decimal64{Mantissa:123456, Exponent:0}},
			0},
		{"lr",
			args{Decimal64{Mantissa: 12345600, Exponent:-1}, Decimal64{Mantissa:123456, Exponent:0}},
			1},
		{"lr",
			args{Decimal64{Mantissa: 12345600, Exponent:-3}, Decimal64{Mantissa:123456, Exponent:0}},
			-1},

		{"lr",
			args{Decimal64{Mantissa: 12345600, Exponent:-4}, Decimal64{Mantissa:123456, Exponent:-2}},
			0},
		{"lr",
			args{Decimal64{Mantissa: 12345600, Exponent:-3}, Decimal64{Mantissa:123456, Exponent:-2}},
			1},
		{"lr",
			args{Decimal64{Mantissa: 12345600, Exponent:-5}, Decimal64{Mantissa:123456, Exponent:2}},
			-1},





	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compare(tt.args.l, tt.args.r); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}