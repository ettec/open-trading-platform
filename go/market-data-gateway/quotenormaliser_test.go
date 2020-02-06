package main

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/common"
	"github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/fix"
	md "github.com/ettec/open-trading-platform/go/market-data-gateway/internal/fix/marketdata"
	"log"
	"strconv"
	"testing"
	"time"
)

func Test_newQuoteNormaliser(t *testing.T) {

}

func Test_quoteNormaliser_close(t *testing.T) {

	toNormaliser := make(chan mdupdate)
	fromNormalise := make(chan *snapshot, 100)

	n := newQuoteNormaliser(toNormaliser, fromNormalise)
	log.Println("normaliser:", n)

	lIds := &listingIdSymbol{1, "A"}
	toNormaliser <- mdupdate{listingIdToSymbol: lIds}

	toNormaliser <- mdupdate{refresh: &refresh{
		MdIncGrp: []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_BID, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 10, 5, "A")},
	}}

	n.close()

	toNormaliser <- mdupdate{refresh: &refresh{
		MdIncGrp: []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_OFFER, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 12, 5, "A")},
	}}

}

func Test_quoteNormaliser_processUpdates(t *testing.T) {
	toNormaliser := make(chan mdupdate)
	fromNormalise := make(chan *snapshot, 100)

	n := newQuoteNormaliser(toNormaliser, fromNormalise)
	log.Println("normaliser:", n)

	lIds := &listingIdSymbol{1, "A"}
	toNormaliser <- mdupdate{listingIdToSymbol: lIds}

	entries := []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_BID, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 10, 5, "A")}

	toNormaliser <- mdupdate{refresh: &refresh{
		MdIncGrp: entries,
	}}

	entries2 := []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_OFFER, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 12, 5, "A")}
	toNormaliser <- mdupdate{refresh: &refresh{
		MdIncGrp: entries2,
	}}

	entries3 := []*md.MDIncGrp{getEntry(md.MDEntryTypeEnum_MD_ENTRY_TYPE_OFFER, md.MDUpdateActionEnum_MD_UPDATE_ACTION_NEW, 11, 2, "A")}
	toNormaliser <- mdupdate{refresh: &refresh{
		MdIncGrp: entries2,
	}}

	toNormaliser <- mdupdate{refresh: &refresh{
		MdIncGrp: entries3,
	}}

	time.Sleep(2 * time.Second)
	var snapshot *snapshot

l:
	for {
		select {
		case s := <-fromNormalise:
			snapshot = s
		default:
			break l
		}
	}

	err := testEqual(snapshot, [5][4]int64{{5, 10, 12, 5}, {0, 0, 11, 2}}, lIds.listingId)
	if err != nil {
		t.Errorf("Books not equal %v", err)
	}

	n.close()
}

func testEqual(snapshot *snapshot, book [5][4]int64, listingId int) error {

	if snapshot.Instrument.Symbol != strconv.Itoa(listingId) {
		return fmt.Errorf("snapshot listing id and listing id are not the same")
	}

	var compare [5][4]int64
	bidIdx := 0
	askIdx := 0
	for _, grp := range snapshot.MdFullGrp {

		if grp.MdEntryType == md.MDEntryTypeEnum_MD_ENTRY_TYPE_BID {
			compare[bidIdx][0] = grp.MdEntrySize.Mantissa
			compare[bidIdx][1] = grp.MdEntryPx.Mantissa
			bidIdx++
		} else if grp.MdEntryType == md.MDEntryTypeEnum_MD_ENTRY_TYPE_OFFER {
			compare[askIdx][3] = grp.MdEntrySize.Mantissa
			compare[askIdx][2] = grp.MdEntryPx.Mantissa
			askIdx++
		}
	}

	if book != compare {
		return fmt.Errorf("Expected book %v does not match book create from snapshot %v", book, compare)
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
