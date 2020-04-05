package quoteaggregator

import (
	"github.com/ettec/open-trading-platform/go/model"
	"reflect"
	"testing"
)

func Test_combineQuotes(t *testing.T) {
	type args struct {
		combinedListingId int32
		quotes            []*model.ClobQuote
	}


	tests := []struct {
		name string
		args args
		want *model.ClobQuote
	}{
		{name: "test combine",
			args: args {combinedListingId: 1, quotes: []*model.ClobQuote{
				{
					ListingId:            2,
					Bids: []*model.ClobLine{
						{ Size: d64(15)  ,  Price: d64 (105)},
						{ Size: d64(13)  ,  Price: d64 (103)},
						{ Size: d64(10)  ,  Price: d64 (100)},
					},
					Offers:               []*model.ClobLine{
						{ Size: d64(10)  ,  Price: d64 (100)},
						{ Size: d64(13)  ,  Price: d64 (103)},
						{ Size: d64(15)  ,  Price: d64 (105)},
					},
					StreamInterrupted:    false,
					StreamStatusMsg:      "",
				},
				{
					ListingId:            3,
					Bids: []*model.ClobLine{
						{ Size: d64(13)  ,  Price: d64 (103)},
						{ Size: d64(12)  ,  Price: d64 (102)},
						{ Size: d64(11)  ,  Price: d64 (101)},
					},
					Offers:               []*model.ClobLine{
						{ Size: d64(11)  ,  Price: d64 (101)},
						{ Size: d64(12)  ,  Price: d64 (102)},
						{ Size: d64(13)  ,  Price: d64 (103)},
					},
					StreamInterrupted:    false,
					StreamStatusMsg:      "",
				},
			}},

			want: &model.ClobQuote{
				ListingId:            1,
				Bids: []*model.ClobLine{
					{ Size: d64(15)  ,  Price: d64 (105), ListingId: 2},
					{ Size: d64(13)  ,  Price: d64 (103), ListingId: 2},
					{ Size: d64(13)  ,  Price: d64 (103), ListingId: 3},
					{ Size: d64(12)  ,  Price: d64 (102), ListingId: 3},
					{ Size: d64(11)  ,  Price: d64 (101), ListingId: 3},
					{ Size: d64(10)  ,  Price: d64 (100), ListingId: 2},
				},
				Offers:
				[]*model.ClobLine{
					{ Size: d64(10)  ,  Price: d64 (100), ListingId: 2},
					{ Size: d64(11)  ,  Price: d64 (101), ListingId: 3},
					{ Size: d64(12)  ,  Price: d64 (102), ListingId: 3},
					{ Size: d64(13)  ,  Price: d64 (103), ListingId: 2},
					{ Size: d64(13)  ,  Price: d64 (103), ListingId: 3},
					{ Size: d64(15)  ,  Price: d64 (105), ListingId: 2},
				},
				StreamInterrupted:    false,
				StreamStatusMsg:      "",
			}},



	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combineQuotes(tt.args.combinedListingId, tt.args.quotes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combineQuotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func d64(mantissa int) *model.Decimal64 {
	return &model.Decimal64{Mantissa: int64(mantissa), Exponent: 0}
}