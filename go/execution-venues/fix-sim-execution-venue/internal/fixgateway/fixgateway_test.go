package fixgateway

import (
	"github.com/ettec/otp-common/model"
	"testing"
)

func toFixString(decimal64 model.Decimal64) string {
	d, s := toFixDecimal(&decimal64)

	str := d.StringFixed(s)
	return str
}

func TestToFixDecimal(t *testing.T) {

	type args struct {
		d model.Decimal64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "basic", args: args{d: model.Decimal64{Mantissa: 1936, Exponent: -2}}, want: "19.36"},
		{name: "basic", args: args{d: model.Decimal64{Mantissa: 1936, Exponent: 0}}, want: "1936"},
		{name: "basic", args: args{d: model.Decimal64{Mantissa: 1936, Exponent: 2}}, want: "193600"},
		{name: "basic", args: args{d: model.Decimal64{Mantissa: -1936, Exponent: -2}}, want: "-19.36"},
		{name: "basic", args: args{d: model.Decimal64{Mantissa: -1936, Exponent: 0}}, want: "-1936"},
		{name: "basic", args: args{d: model.Decimal64{Mantissa: -1936, Exponent: 2}}, want: "-193600"},
		{name: "basic", args: args{d: model.Decimal64{Mantissa: 193600, Exponent: -1}}, want: "19360.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toFixString(tt.args.d)
			if got != tt.want {
				t.Errorf("toFixDecimal() got = %v, want %v", got, tt.want)
			}

		})
	}
}
