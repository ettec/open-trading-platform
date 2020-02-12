package fixsim

import (
	"fmt"

	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/fix"
	md "github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/model"
	"reflect"
	"strconv"
	"testing"
)

type testMarketDataClient struct {

	refreshChan chan *md.MarketDataIncrementalRefresh
	errorChan chan error
	subscribeChan chan string
	closeSignalChan chan bool

}

func newTestMarketDataClient() *testMarketDataClient {
	t := &testMarketDataClient{

		refreshChan : make(chan *md.MarketDataIncrementalRefresh, 100),
		errorChan : make(chan error, 100),
		subscribeChan : make(chan string, 100),
		closeSignalChan : make(chan bool, 100),
	}
	return t
}

func (t *testMarketDataClient) connect(connectionId string) (receiveIncRefreshFn, error) {
	return func() (refresh *md.MarketDataIncrementalRefresh, err error) {
		select {
		case e := <-t.errorChan:
			return nil, e
		case r := <-t.refreshChan:
			return r, nil
		}

	}, nil

}

func (t *testMarketDataClient) subscribe(symbol string, subscriberId string) error {
	t.subscribeChan<-symbol
	return nil
}

func (t *testMarketDataClient) close() error {
	t.closeSignalChan<-true
	return nil
}




func Test_fixSimConnection_close(t *testing.T) {



	listingIdToSym := map[int]string{1:"A", 2:"B"}
	tmd := newTestMarketDataClient()

	out := make(chan *model.ClobQuote, 100)
	n := NewFixSimConnection( tmd, "testName", toLookupFunc(listingIdToSym), out )


	n.Subscribe(1)

	if sym := <-tmd.subscribeChan; sym != "A" {
		t.Errorf("expected symbol A")
	}

	tmd.refreshChan <- &md.MarketDataIncrementalRefresh{
		MdIncGrp: []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_BID, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 10, 5, "A")},
	}

	q := <-out
	if q.ListingId != 1 {
		t.Errorf("expected quote for listing id 1")
	}

	n.Close()

	<-tmd.closeSignalChan
	tmd.errorChan <- fmt.Errorf("an error")

	if _, ok := <-out; ok {
		t.Errorf("expected channel t be closed")
	}

}





func Test_quoteNormaliser_processUpdates(t *testing.T) {

	tmd := newTestMarketDataClient()

	out := make(chan *model.ClobQuote, 100)
	listingIdToSym := map[int]string{1:"A", 2:"B"}
	n := NewFixSimConnection( tmd, "testName", toLookupFunc(listingIdToSym), out )


	n.Subscribe(1)

	symbol := <-tmd.subscribeChan
	if symbol != "A" {
		t.Errorf("exepcted subscribe call for symbol A")
	}

	entries := []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_BID, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 10, 5, "A")}

	n.refreshChan <- &md.MarketDataIncrementalRefresh{
		MdIncGrp: entries,
	}

	entries2 := []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_OFFER, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 12, 5, "A")}

	n.refreshChan <- &md.MarketDataIncrementalRefresh{
		MdIncGrp: entries2,
	}

	entries3 := []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_OFFER, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 11, 2, "A")}
	n.refreshChan <- &md.MarketDataIncrementalRefresh{
		MdIncGrp: entries3,
	}

	q := <-out
	q = <-out
	q = <-out

	err := testEqual(q, [5][4]int64{{5, 10, 11, 2}, {0, 0, 12, 5}}, 1)
	if err != nil {
		t.Errorf("Books not equal %v", err)
	}
}

func toLookupFunc(listingIdToSym map[int]string) func(listingId int) (s string, err error) {
	return func(listingId int) (s string, err error) {
		if sym, ok := listingIdToSym[listingId]; ok {
			return sym, nil
		} else {
			return "", fmt.Errorf("no symbol for listing id %v", listingId)
		}
	}
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
