package model

import (
	"github.com/shopspring/decimal"
	"math"
	"reflect"
	"testing"
)

func TestToDecimal64(t *testing.T) {

	result := ToDecimal64(decimal.New(10, 2))
	resultAsInt64 := int64(float64(result.Mantissa) * (math.Pow(10, float64(result.Exponent))))

	if resultAsInt64 != 1000 {
		t.Fatalf("expected %v got %v", 1000, resultAsInt64)
	}

}

func TestToDecimal64IrrationalNumber(t *testing.T) {

	one := decimal.New(1, 0)
	three := decimal.New(3, 0)
	irn := one.Div(three)

	result := ToDecimal64(irn)

	resultAsFloat64 := float64(result.Mantissa) * (math.Pow(10, float64(result.Exponent)))

	expectedResult := 0.33333333333333333333
	if !equal(resultAsFloat64, expectedResult, 0.00000001) {
		t.Fatalf("expected %v got %v", expectedResult, resultAsFloat64)
	}

}

func equal(a, b, delta float64) bool {
	return math.Abs(a-b) <= delta
}

func TestTimestamp_After(t1 *testing.T) {
	type fields struct {
		Seconds     int64
		Nanoseconds int32
	}
	type args struct {
		o *Timestamp
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "seconds greater",
			fields: fields{Seconds: 3, Nanoseconds: 1},
			args:   args{o: &Timestamp{Seconds: 2, Nanoseconds: 1}},
			want:   true,
		},

		{
			name:   "seconds less",
			fields: fields{Seconds: 1, Nanoseconds: 1},
			args:   args{o: &Timestamp{Seconds: 2, Nanoseconds: 1}},
			want:   false,
		},

		{
			name:   "nanoseconds greater",
			fields: fields{Seconds: 3, Nanoseconds: 5},
			args:   args{o: &Timestamp{Seconds: 3, Nanoseconds: 1}},
			want:   true,
		},

		{
			name:   "nanoseconds less",
			fields: fields{Seconds: 2, Nanoseconds: 5},
			args:   args{o: &Timestamp{Seconds: 2, Nanoseconds: 10}},
			want:   false,
		},

		{
			name:   "same",
			fields: fields{Seconds: 2, Nanoseconds: 5},
			args:   args{o: &Timestamp{Seconds: 2, Nanoseconds: 5}},
			want:   false,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Timestamp{
				Seconds:     tt.fields.Seconds,
				Nanoseconds: tt.fields.Nanoseconds,
			}
			if got := t.After(tt.args.o); got != tt.want {
				t1.Errorf("After() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_Before(t1 *testing.T) {
	type fields struct {
		Seconds     int64
		Nanoseconds int32
	}
	type args struct {
		o *Timestamp
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "seconds greater",
			fields: fields{Seconds: 3, Nanoseconds: 1},
			args:   args{o: &Timestamp{Seconds: 2, Nanoseconds: 1}},
			want:   false,
		},

		{
			name:   "seconds less",
			fields: fields{Seconds: 1, Nanoseconds: 1},
			args:   args{o: &Timestamp{Seconds: 2, Nanoseconds: 1}},
			want:   true,
		},

		{
			name:   "nanoseconds greater",
			fields: fields{Seconds: 3, Nanoseconds: 5},
			args:   args{o: &Timestamp{Seconds: 3, Nanoseconds: 1}},
			want:   false,
		},

		{
			name:   "nanoseconds less",
			fields: fields{Seconds: 2, Nanoseconds: 5},
			args:   args{o: &Timestamp{Seconds: 2, Nanoseconds: 10}},
			want:   true,
		},

		{
			name:   "same",
			fields: fields{Seconds: 2, Nanoseconds: 5},
			args:   args{o: &Timestamp{Seconds: 2, Nanoseconds: 5}},
			want:   false,
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Timestamp{
				Seconds:     tt.fields.Seconds,
				Nanoseconds: tt.fields.Nanoseconds,
			}
			if got := t.Before(tt.args.o); got != tt.want {
				t1.Errorf("Before() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal64_AsDecimal(t *testing.T) {

	type fields struct {
		Mantissa int64
		Exponent int32
	}

	tests := []struct {
		name   string
		fields fields
		want   decimal.Decimal
	}{
		{"test", fields{Mantissa: 1, Exponent: 0}, decimal.New(1, 0)},
		{"test", fields{Mantissa: 231, Exponent: -7}, decimal.New(231, -7)},
		{"test", fields{Mantissa: 123, Exponent: 5}, decimal.New(123, 5)},
		{"test", fields{Mantissa: 0, Exponent: 0}, decimal.New(0, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Decimal64{
				Mantissa: tt.fields.Mantissa,
				Exponent: tt.fields.Exponent,
			}
			if got := m.AsDecimal(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal64_Add(t *testing.T) {
	type args struct {
		t Decimal64
		o Decimal64
	}
	tests := []struct {
		name string
		args args
		want Decimal64
	}{
		{"test", args{Decimal64{Mantissa: 0, Exponent: 0}, Decimal64{Mantissa: 0, Exponent: 0}},
			Decimal64{Mantissa: 0, Exponent: 0}},
		{"test", args{Decimal64{Mantissa: 1, Exponent: 0}, Decimal64{Mantissa: 0, Exponent: 0}},
			Decimal64{Mantissa: 1, Exponent: 0}},
		{"test", args{Decimal64{Mantissa: 1, Exponent: -1}, Decimal64{Mantissa: 1, Exponent: -1}},
			Decimal64{Mantissa: 2, Exponent: -1}},
		{"test", args{Decimal64{Mantissa: 3, Exponent: 1}, Decimal64{Mantissa: 8, Exponent: 1}},
			Decimal64{Mantissa: 110, Exponent: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.t.Add(&tt.args.o)
			if !tt.args.t.Equal(&tt.want) {
				t.Errorf("Add() = %v, want %v", tt.args.t, tt.want)
			}
		})
	}
}
