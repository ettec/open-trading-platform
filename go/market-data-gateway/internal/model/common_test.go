package model

import "testing"

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
			


	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Compare(tt.args.l, tt.args.r); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}