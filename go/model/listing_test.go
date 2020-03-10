package model

import "testing"

func TestListing_RoundToTickSize(t *testing.T) {

	tst := 		&TickSizeTable {
		Entries: []*TickSizeEntry{
			{
				LowerPriceBound:      &Decimal64{Mantissa: -100, Exponent:0},
				UpperPriceBound:      &Decimal64{Mantissa: -10, Exponent:0},
				TickSize:             &Decimal64{Mantissa: 1, Exponent:-1},
			},
			{
				LowerPriceBound:      &Decimal64{Mantissa: -10, Exponent:0},
				UpperPriceBound:      &Decimal64{Mantissa: 10, Exponent:0},
				TickSize:             &Decimal64{Mantissa: 1, Exponent:-2},
			},
			{
				LowerPriceBound:      &Decimal64{Mantissa: 10, Exponent:0},
				UpperPriceBound:      &Decimal64{Mantissa: 100, Exponent:0},
				TickSize:             &Decimal64{Mantissa: 1, Exponent:-1},
			},
			{
				LowerPriceBound:      &Decimal64{Mantissa: 100, Exponent:0},
				UpperPriceBound:      &Decimal64{Mantissa: 1000, Exponent:0},
				TickSize:             &Decimal64{Mantissa: 1, Exponent:0},
			},
		},
	}



	tests := []struct {
		name    string
		tst  *TickSizeTable
		price float64
		result *Decimal64
		wantErr bool
	}{
		{"test", tst, 9.13211,
			&Decimal64{Mantissa:913, Exponent:-2,},false },
		{"test", tst, 19.132,
			&Decimal64{Mantissa:191, Exponent:-1,},false },
		{"test", tst, 200.1,
			&Decimal64{Mantissa:200, Exponent:0,},false },
		{"test", tst, -2.1163,
			&Decimal64{Mantissa:-212, Exponent:-2,},false },

		{"test", tst, 2116.3,
			nil,true },

		{"test", tst, 9.997,
			&Decimal64{Mantissa:10, Exponent:0,},false },

		{"test", tst, 9.9945,
			&Decimal64{Mantissa:999, Exponent:-2,},false },
		{"test", tst, 9.999945,
			&Decimal64{Mantissa:10, Exponent:0,},false },

		{"test", tst, 10.00001,
			&Decimal64{Mantissa:10, Exponent:0,},false },


		{"test", tst, -9.997,
			&Decimal64{Mantissa:-10, Exponent:0,},false },

		{"test", tst, -9.9945,
			&Decimal64{Mantissa:-999, Exponent:-2,},false },
		{"test", tst, -9.999945,
			&Decimal64{Mantissa:-10, Exponent:0,},false },

		{"test", tst, -10.00001,
			&Decimal64{Mantissa:-10, Exponent:0,},false },





	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Listing{
				TickSize:             tt.tst,
			}

			d, err := m.RoundToNearestTick(tt.price);
			if  (err != nil) != tt.wantErr {
				t.Errorf("RoundToNearestTick() error = %v, wantErr %v", err, tt.wantErr)
			}

			if  !d.Equal(tt.result) {
				t.Errorf("RoundToNearestTick() price = %v, wanted %v", d, tt.result)
			}

		})
	}
}

func Test_compare(t *testing.T) {


	type args struct {
		f1    float64
		f2    float64
		delta float64
	}


	tests := []struct {
		name string
		args args
		want int
	}{
		{"test", args{-5, 10, 0.001}, -1},
		{"test", args{-5, -10, 0.001}, 1},
		{"test", args{-5, -5, 0.001}, 0},
		{"test", args{-2.1163, 10, 0.00001}, -1},


	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compare(tt.args.f1, tt.args.f2, tt.args.delta); got != tt.want {
				t.Errorf("compare() = %v, want %v", got, tt.want)
			}
		})
	}
}