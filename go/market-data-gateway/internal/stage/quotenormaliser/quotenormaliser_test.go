package quotenormaliser

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/fix"
	md "github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/stage"
	"log"
	"reflect"
	"strconv"
	"testing"
)

type testQuoteSink struct {
	out chan *model.ClobQuote
}

func (s testQuoteSink) send(quote *model.ClobQuote) {
	s.out <- quote
}

func Test_quoteNormaliser_close(t *testing.T) {

	fromNormalise := make(chan *model.ClobQuote, 100)

	n := newQuoteNormaliser(testQuoteSink{fromNormalise})
	log.Println("normaliser:", n)

	lIds := stage.ListingIdSymbol{ListingId: 1, Symbol: "A"}

	done := make(chan bool)
	go func() {
		n.registerMapping(lIds)
		done<-true
	}()
	n.readInputChannel()
	<-done

	n.sendRefresh(&stage.Refresh{
		MdIncGrp: []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_BID, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 10, 5, "A")},
	})

	n.readInputChannel()

	n.close()

	close := n.readInputChannel()

	if !close {
		t.Errorf("expected close to be true")
	}

}

func Test_quoteNormaliser_processUpdates(t *testing.T) {

	fromNormalise := make(chan *model.ClobQuote, 100)
	n := newQuoteNormaliser( testQuoteSink{fromNormalise})
	defer n.close()
	log.Println("normaliser:", n)

	lIds := stage.ListingIdSymbol{ListingId: 1, Symbol: "A"}

	done := make(chan bool)
	go func() {
		n.registerMapping(lIds)
		done<-true
	}()
	n.readInputChannel()
	<-done

	entries := []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_BID, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 10, 5, "A")}

	n.sendRefresh(&stage.Refresh{
		MdIncGrp: entries,
	})

	entries2 := []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_OFFER, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 12, 5, "A")}

	n.sendRefresh(&stage.Refresh{
		MdIncGrp: entries2,
	})


	entries3 := []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_OFFER, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 11, 2, "A")}
	n.sendRefresh(&stage.Refresh{
		MdIncGrp: entries3,
	})


	invoke(n.readInputChannel, 3)

	err := testEqual(getLastSnapshot(fromNormalise), [5][4]int64{{5, 10, 11, 2}, {0, 0, 12, 5}}, lIds.ListingId)
	if err != nil {
		t.Errorf("Books not equal %v", err)
	}
}

func invoke(f func() bool, times int) {

	for i := 0; i < times; i++ {
		f()
	}
}

func getLastSnapshot(fromNormalise chan *model.ClobQuote) *model.ClobQuote {
	var quote *model.ClobQuote

loop:
	for {
		select {
		case s := <-fromNormalise:
			quote = s
		default:
			break loop
		}
	}
	return quote
}

func testEqual(quote *model.ClobQuote, book [5][4]int64, listingId int) error {

	if quote.ListingId != int32(listingId) {
		return fmt.Errorf("quote listing id and listing id are not the same")
	}

	var compare [5][4]int64

	for idx, line := range quote.Bids {
		compare[idx][0] = line.Size.Mantissa
		compare[idx][1] = line.Price.Mantissa
	}

	for idx, line := range quote.Offers {
		compare[idx][3] = line.Size.Mantissa
		compare[idx][2] = line.Price.Mantissa
	}

	if book != compare {
		return fmt.Errorf("Expected book %v does not match book create from quote %v", book, compare)
	}

	return nil
}

var id int = 0

func getNextId() string {
	id++
	return strconv.Itoa(id)
}

func getEntry(mt md.MDEntryTypeEnum, ma md.MDUpdateActionEnum, price int64, size int64, symbol string) *md.MDIncGrp {
	instrument := &common.Instrument{Symbol: symbol}
	entry := &md.MDIncGrp{
		MdEntryId:      getNextId(),
		MdEntryType:    mt,
		MdUpdateAction: ma,
		MdEntryPx:      &fix.Decimal64{Mantissa: price, Exponent: 0},
		MdEntrySize:    &fix.Decimal64{Mantissa: size, Exponent: 0},
		Instrument:     instrument,
	}
	return entry
}

func Test_updateAsksWithInserts(t *testing.T) {
	type args struct {
		asks   []*model.ClobLine
		update md.MDIncGrp
	}

	tests := []struct {
		name string
		args args
		want []*model.ClobLine
	}{

		{
			"insert ask into empty book",
			args{
				asks: []*model.ClobLine{},
				update: md.MDIncGrp{MdEntryId: "A", MdEntrySize: f64(20), MdEntryPx: f64(6),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW},
			},
			[]*model.ClobLine{{EntryId: "A", Size: d64(20), Price: d64(6)}},
		},

		{
			"insert ask into middle of book",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(2)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(6)}},
				update: md.MDIncGrp{MdEntryId: "X", MdEntrySize: f64(20), MdEntryPx: f64(3),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(2)},
				{EntryId: "X", Size: d64(20), Price: d64(3)},
				{EntryId: "B", Size: d64(20), Price: d64(4)},
				{EntryId: "C", Size: d64(20), Price: d64(6)}},
		},


		{
			"insert ask at same price",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(2)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(6)}},
				update: md.MDIncGrp{MdEntryId: "X", MdEntrySize: f64(20), MdEntryPx: f64(4),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(2)},
				{EntryId: "B", Size: d64(20), Price: d64(4)},
				{EntryId: "X", Size: d64(20), Price: d64(4)},
				{EntryId: "C", Size: d64(20), Price: d64(6)}},
		},

		{
			"insert ask at top of book ",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(2)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(6)}},
				update: md.MDIncGrp{MdEntryId: "X", MdEntrySize: f64(20), MdEntryPx: f64(1),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW},
			},
			[]*model.ClobLine{
				{EntryId: "X", Size: d64(20), Price: d64(1)},
				{EntryId: "A", Size: d64(20), Price: d64(2)},
				{EntryId: "B", Size: d64(20), Price: d64(4)},
				{EntryId: "C", Size: d64(20), Price: d64(6)}},
		},

		{
			"insert ask at bottom of book ",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(2)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(6)}},
				update: md.MDIncGrp{MdEntryId: "X", MdEntrySize: f64(20), MdEntryPx: f64(8),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(2)},
				{EntryId: "B", Size: d64(20), Price: d64(4)},
				{EntryId: "C", Size: d64(20), Price: d64(6)},
				{EntryId: "X", Size: d64(20), Price: d64(8)}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateClobLines(tt.args.asks, &tt.args.update, false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateClobLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateAsksWithUpdates(t *testing.T) {
	type args struct {
		asks   []*model.ClobLine
		update md.MDIncGrp
	}

	tests := []struct {
		name string
		args args
		want []*model.ClobLine
	}{

		{
			"update ask quantity",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(2)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(6)}},
				update: md.MDIncGrp{MdEntryId: "B", MdEntrySize: f64(10), MdEntryPx: f64(4),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(2)},
				{EntryId: "B", Size: d64(10), Price: d64(4)},
				{EntryId: "C", Size: d64(20), Price: d64(6)}},
		},

		{
			"update ask price - no order change",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(2)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(6)}},
				update: md.MDIncGrp{MdEntryId: "B", MdEntrySize: f64(20), MdEntryPx: f64(3),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(2)},
				{EntryId: "B", Size: d64(20), Price: d64(3)},
				{EntryId: "C", Size: d64(20), Price: d64(6)}},
		},

		{
			"update ask price down to same as other - order change",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(2)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(6)}},

				update: md.MDIncGrp{MdEntryId: "B", MdEntrySize: f64(20), MdEntryPx: f64(6),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(2)},
				{EntryId: "C", Size: d64(20), Price: d64(6)},
				{EntryId: "B", Size: d64(20), Price: d64(6)}},
		},

		{
			"update ask price up to same as other - order change",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(2)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(6)}},

				update: md.MDIncGrp{MdEntryId: "B", MdEntrySize: f64(20), MdEntryPx: f64(2),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(2)},
				{EntryId: "B", Size: d64(20), Price: d64(2)},
				{EntryId: "C", Size: d64(20), Price: d64(6)}},
		},

		{
			"update ask price up to top of book",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(2)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(6)}},

				update: md.MDIncGrp{MdEntryId: "B", MdEntrySize: f64(20), MdEntryPx: f64(1),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE},
			},
			[]*model.ClobLine{
				{EntryId: "B", Size: d64(20), Price: d64(1)},
				{EntryId: "A", Size: d64(20), Price: d64(2)},
				{EntryId: "C", Size: d64(20), Price: d64(6)}},
		},

		{
			"update ask price up to bottom of book",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(2)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(6)}},

				update: md.MDIncGrp{MdEntryId: "B", MdEntrySize: f64(20), MdEntryPx: f64(8),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(2)},
				{EntryId: "C", Size: d64(20), Price: d64(6)},
				{EntryId: "B", Size: d64(20), Price: d64(8)}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateClobLines(tt.args.asks, &tt.args.update, false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateClobLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateBidsWithInserts(t *testing.T) {
	type args struct {
		bids   []*model.ClobLine
		update md.MDIncGrp
	}

	tests := []struct {
		name string
		args args
		want []*model.ClobLine
	}{

		{
			"insert bid into empty book",
			args{
				bids: []*model.ClobLine{},
				update: md.MDIncGrp{MdEntryId: "A", MdEntrySize: f64(20), MdEntryPx: f64(6),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW},
			},
			[]*model.ClobLine{{EntryId: "A", Size: d64(20), Price: d64(6)}},
		},

		{
			"insert bid into middle of book",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "X", MdEntrySize: f64(20), MdEntryPx: f64(3),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "B", Size: d64(20), Price: d64(4)},
				{EntryId: "X", Size: d64(20), Price: d64(3)},
				{EntryId: "C", Size: d64(20), Price: d64(2)}},
		},

		{
			"insert bid into middle of book",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "X", MdEntrySize: f64(20), MdEntryPx: f64(3),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "B", Size: d64(20), Price: d64(4)},
				{EntryId: "X", Size: d64(20), Price: d64(3)},
				{EntryId: "C", Size: d64(20), Price: d64(2)}},
		},

		{
			"insert bid at same price",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "X", MdEntrySize: f64(20), MdEntryPx: f64(4),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "B", Size: d64(20), Price: d64(4)},
				{EntryId: "X", Size: d64(20), Price: d64(4)},
				{EntryId: "C", Size: d64(20), Price: d64(2)}},
		},

		{
			"insert bid at top of book ",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "X", MdEntrySize: f64(20), MdEntryPx: f64(8),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW},
			},
			[]*model.ClobLine{
				{EntryId: "X", Size: d64(20), Price: d64(8)},
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "B", Size: d64(20), Price: d64(4)},
				{EntryId: "C", Size: d64(20), Price: d64(2)}},
		},

		{
			"insert bid at bottom of book ",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "X", MdEntrySize: f64(20), MdEntryPx: f64(1),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "B", Size: d64(20), Price: d64(4)},
				{EntryId: "C", Size: d64(20), Price: d64(2)},
				{EntryId: "X", Size: d64(20), Price: d64(1)}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateClobLines(tt.args.bids, &tt.args.update, true); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateClobLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateBidsWithUpdates(t *testing.T) {
	type args struct {
		bids   []*model.ClobLine
		update md.MDIncGrp
	}

	tests := []struct {
		name string
		args args
		want []*model.ClobLine
	}{

		{
			"update bid quantity",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "B", MdEntrySize: f64(10), MdEntryPx: f64(4),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "B", Size: d64(10), Price: d64(4)},
				{EntryId: "C", Size: d64(20), Price: d64(2)}},
		},

		{
			"update bid price - no order change",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "B", MdEntrySize: f64(10), MdEntryPx: f64(3),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "B", Size: d64(10), Price: d64(3)},
				{EntryId: "C", Size: d64(20), Price: d64(2)}},
		},

		{
			"update bid price down to same as other - order change",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(3)}},

				update: md.MDIncGrp{MdEntryId: "B", MdEntrySize: f64(20), MdEntryPx: f64(3),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "C", Size: d64(20), Price: d64(3)},
				{EntryId: "B", Size: d64(20), Price: d64(3)}},
		},

		{
			"update bid price up to same as other - order change",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(3)}},

				update: md.MDIncGrp{MdEntryId: "B", MdEntrySize: f64(20), MdEntryPx: f64(6),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "B", Size: d64(20), Price: d64(6)},
				{EntryId: "C", Size: d64(20), Price: d64(3)}},
		},

		{
			"update bid price up to top of book",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(3)}},

				update: md.MDIncGrp{MdEntryId: "B", MdEntrySize: f64(20), MdEntryPx: f64(8),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE},
			},
			[]*model.ClobLine{
				{EntryId: "B", Size: d64(20), Price: d64(8)},
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "C", Size: d64(20), Price: d64(3)}},
		},

		{
			"update bid price up to bottom of book",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(3)}},

				update: md.MDIncGrp{MdEntryId: "B", MdEntrySize: f64(20), MdEntryPx: f64(2),
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_CHANGE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "C", Size: d64(20), Price: d64(3)},
				{EntryId: "B", Size: d64(20), Price: d64(2)}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateClobLines(tt.args.bids, &tt.args.update, true); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateClobLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateBidsWithDelete(t *testing.T) {
	type args struct {
		bids   []*model.ClobLine
		update md.MDIncGrp
	}

	tests := []struct {
		name string
		args args
		want []*model.ClobLine
	}{

		{
			"delete from middle of book",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "B",
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_DELETE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "C", Size: d64(20), Price: d64(2)}},
		},
		{
			"delete from top of book",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "A",
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_DELETE},
			},
			[]*model.ClobLine{

				{EntryId: "B", Size: d64(20), Price: d64(4)},
				{EntryId: "C", Size: d64(20), Price: d64(2)}},
		},
		{
			"delete from bottom of book",
			args{
				bids: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "C",
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_DELETE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "B", Size: d64(20), Price: d64(4)},},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateClobLines(tt.args.bids, &tt.args.update, true); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateClobLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateAsksWithDelete(t *testing.T) {
	type args struct {
		asks   []*model.ClobLine
		update md.MDIncGrp
	}

	tests := []struct {
		name string
		args args
		want []*model.ClobLine
	}{

		{
			"delete from middle of book",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "B",
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_DELETE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "C", Size: d64(20), Price: d64(2)}},
		},
		{
			"delete from top of book",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "A",
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_DELETE},
			},
			[]*model.ClobLine{

				{EntryId: "B", Size: d64(20), Price: d64(4)},
				{EntryId: "C", Size: d64(20), Price: d64(2)}},
		},
		{
			"delete from bottom of book",
			args{
				asks: []*model.ClobLine{
					{EntryId: "A", Size: d64(20), Price: d64(6)},
					{EntryId: "B", Size: d64(20), Price: d64(4)},
					{EntryId: "C", Size: d64(20), Price: d64(2)}},
				update: md.MDIncGrp{MdEntryId: "C",
					MdUpdateAction: md.MDUpdateActionEnum_MD_UPDATE_ACTION_DELETE},
			},
			[]*model.ClobLine{
				{EntryId: "A", Size: d64(20), Price: d64(6)},
				{EntryId: "B", Size: d64(20), Price: d64(4)},},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateClobLines(tt.args.asks, &tt.args.update, false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateClobLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func d64(mantissa int) *model.Decimal64 {
	return &model.Decimal64{Mantissa: int64(mantissa), Exponent: 0}
}

func f64(mantissa int) *fix.Decimal64 {
	return &fix.Decimal64{Mantissa: int64(mantissa), Exponent: 0}
}
