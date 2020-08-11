package main

import (
	"encoding/json"
	"fmt"
	"github.com/ettec/otp-common/model"
	"github.com/ettec/otp-common/strategy"
	"time"
)


type vwapParameters struct {
	utcStartTimeSecs int64
	utcEndTimeSecs   int64
	buckets          int
}


func executeAsVwapStrategy(om *strategy.Strategy, buckets []bucket,  listing *model.Listing) {

	go func() {

		om.Log.Println("order initialised")

		ticker := time.NewTicker(1 * time.Second)

		for {
			done, err := om.CheckIfDone()
			if err != nil {
				om.ErrLog.Printf("failed to check if done, cancelling order:%v", err)
				om.Cancel()
			}

			if done {
				break
			}

			select {
			case <-ticker.C:
				nowUtc := time.Now().Unix()
				shouldHaveSentQty := &model.Decimal64{}
				for  i :=0; i< len(buckets); i++ {
					if buckets[i].utcStartTimeSecs <= nowUtc {
						shouldHaveSentQty.Add(&buckets[i].quantity)
					}
				}

				sentQty  := &model.Decimal64{}
				sentQty.Add(om.ParentOrder.GetTradedQuantity())
				sentQty.Add(om.ParentOrder.GetExposedQuantity())

				if sentQty.LessThan(shouldHaveSentQty) {
					shouldHaveSentQty.Sub(sentQty)
					err := om.SendChildOrder(om.ParentOrder.Side,  shouldHaveSentQty, om.ParentOrder.Price, listing)
					if err != nil {
						om.CancelOrderWithErrorMsg(fmt.Sprintf("failed to send child order:%v" , err))
					}
				}

			case errMsg := <-om.CancelChan:
				if errMsg != "" {
					om.ParentOrder.ErrorMessage = errMsg
				}
				err := om.CancelParentOrder(func(listingId int32) *model.Listing {
					return listing
				})
				if err != nil {
					om.ErrLog.Printf("failed to cancel order:%v", err)
				}
			case co, ok := <-om.ChildOrderUpdateChan:
				om.OnChildOrderUpdate(ok, co)
			}

		}

		ticker.Stop()
	}()
}

func getBucketsFromParamsString(vwapParamsJson string,  quantity *model.Decimal64, listing *model.Listing) ([]bucket, error) {
	vwapParameters := &vwapParameters{}
	err := json.Unmarshal([]byte(vwapParamsJson), vwapParameters)
	if err != nil {
		return nil, err
	}

	numBuckets := vwapParameters.buckets
	if numBuckets == 0 {
		if quantity.ToFloat() > 100 {
			numBuckets = 100
		} else {
			numBuckets = int(quantity.ToFloat())
		}
	}

	buckets := getBuckets(listing, vwapParameters.utcStartTimeSecs, vwapParameters.utcEndTimeSecs, numBuckets, quantity)
	return buckets, nil
}

type bucket struct {
	quantity         model.Decimal64
	utcStartTimeSecs int64
	utcEndTimeSecs   int64
}

func getBuckets(listing *model.Listing, utcStartTimeSecs int64, utcEndTimeSecs int64, buckets int, quantity *model.Decimal64) (result []bucket) {
	// need historical traded volume data, for now use a TWAP profile
	bucketInterval := (utcEndTimeSecs - utcStartTimeSecs) / int64(buckets)

	fBuckets := float64(buckets)
	fQuantity := quantity.ToFloat()
	bucketQnt := fQuantity / fBuckets

	startTime := utcStartTimeSecs
	endTime := startTime + bucketInterval

	for i := 0; i < buckets; i++ {
		bucket := bucket{
			quantity:         *listing.RoundToLotSize(bucketQnt),
			utcStartTimeSecs: startTime,
			utcEndTimeSecs:   endTime,
		}
		result = append(result, bucket)

		startTime = endTime
		endTime = endTime + bucketInterval
	}

	var totalQnt model.Decimal64
	for _, bucket := range result {
		totalQnt.Add(&bucket.quantity)
	}

	quantity.Sub(&totalQnt)
	if result != nil {
		result[len(result)-1].quantity.Add(quantity)
	}

	return result
}
