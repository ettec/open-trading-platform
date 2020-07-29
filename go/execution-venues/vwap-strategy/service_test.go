package main

import (
	"github.com/ettec/otp-common/model"
	"reflect"
	"testing"
)

func Test_getBuckets(t *testing.T) {

	testListing := &model.Listing{}

	type args struct {
		listing          *model.Listing
		utcStartTimeSecs int
		utcEndTimeSecs   int
		buckets          int
		quantity         *model.Decimal64
	}
	tests := []struct {
		name       string
		args       args
		wantResult []bucket
	}{
		{
			"first",
			args{
				testListing,
				0,
				15,
				5,
				&model.Decimal64{
					Mantissa: 10,
				},
			},

			[]bucket{bucket{
				quantity: model.Decimal64{Mantissa: 2}, utcStartTimeSecs: 0, utcEndTimeSecs: 3,
			},
				{
					quantity: model.Decimal64{Mantissa: 2}, utcStartTimeSecs: 3, utcEndTimeSecs: 6,
				},
				{
					quantity: model.Decimal64{Mantissa: 2}, utcStartTimeSecs: 6, utcEndTimeSecs: 9,
				},
				{
					quantity: model.Decimal64{Mantissa: 2}, utcStartTimeSecs: 9, utcEndTimeSecs: 12,
				},
				{
					quantity: model.Decimal64{Mantissa: 2}, utcStartTimeSecs: 12, utcEndTimeSecs: 15,
				}},
		},
		{
			"oddQnt",
			args{
				testListing,
				0,
				15,
				5,
				&model.Decimal64{
					Mantissa: 11,
				},
			},

			[]bucket{bucket{
				quantity: model.Decimal64{Mantissa: 2}, utcStartTimeSecs: 0, utcEndTimeSecs: 3,
			},
				{
					quantity: model.Decimal64{Mantissa: 2}, utcStartTimeSecs: 3, utcEndTimeSecs: 6,
				},
				{
					quantity: model.Decimal64{Mantissa: 2}, utcStartTimeSecs: 6, utcEndTimeSecs: 9,
				},
				{
					quantity: model.Decimal64{Mantissa: 2}, utcStartTimeSecs: 9, utcEndTimeSecs: 12,
				},
				{
					quantity: model.Decimal64{Mantissa: 3}, utcStartTimeSecs: 12, utcEndTimeSecs: 15,
				}},
		},
		{
			"notDivisableTime",
			args{
				testListing,
				0,
				8,
				3,
				&model.Decimal64{
					Mantissa: 11,
				},
			},

			[]bucket{bucket{
				quantity: model.Decimal64{Mantissa: 4}, utcStartTimeSecs: 0, utcEndTimeSecs: 2,
			},
				{
					quantity: model.Decimal64{Mantissa: 4}, utcStartTimeSecs: 2, utcEndTimeSecs: 4,
				},
				{
					quantity: model.Decimal64{Mantissa: 3}, utcStartTimeSecs: 4, utcEndTimeSecs: 6,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := getBuckets(tt.args.listing, tt.args.utcStartTimeSecs, tt.args.utcEndTimeSecs, tt.args.buckets, tt.args.quantity); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("getBuckets() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
