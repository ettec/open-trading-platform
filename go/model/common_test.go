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