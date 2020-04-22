package main

import (
	"context"
	"github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/model"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"reflect"
	"testing"
)

type testEvClient struct {
	params []*executionvenue.CreateAndRouteOrderParams
}

func (t *testEvClient) CreateAndRouteOrder(ctx context.Context, in *executionvenue.CreateAndRouteOrderParams, opts ...grpc.CallOption) (*executionvenue.OrderId, error) {
	t.params = append(t.params, in)
	id, _ := uuid.NewUUID()
	return &executionvenue.OrderId{
		OrderId: id.String(),
	}, nil
}

func (t *testEvClient) CancelOrder(ctx context.Context, in *executionvenue.CancelOrderParams, opts ...grpc.CallOption) (*model.Empty, error) {
	panic("implement me")
}

func Test_submitBuyOrders(t *testing.T) {

	orderId := "a"
	evId := "testev"

	q := &model.ClobQuote{
		Offers: []*model.ClobLine{
			{Size: model.IasD(10), Price: model.IasD(100), ListingId: 1},
			{Size: model.IasD(10), Price: model.IasD(110), ListingId: 2},
			{Size: model.IasD(10), Price: model.IasD(120), ListingId: 1},
			{Size: model.IasD(10), Price: model.IasD(130), ListingId: 2},
			{Size: model.IasD(10), Price: model.IasD(140), ListingId: 1},
			{Size: model.IasD(10), Price: model.IasD(150), ListingId: 1},
		},
	}

	mo := newManagedOrder(model.NewOrder(orderId, model.Side_BUY, model.IasD(50), model.IasD(130), 0, "oi", "od"))

	listing1 := &model.Listing{Id: 1}
	listing2 := &model.Listing{Id: 2}
	underlyingListings := map[int32]*model.Listing{
		1: listing1,
		2: listing2,
	}

	client := &testEvClient{}
	submitBuyOrders(q, mo, underlyingListings, evId, client)

	if len(client.params) != 5 {
		t.FailNow()
	}

	expectedParams := []*executionvenue.CreateAndRouteOrderParams{{
			OrderSide:            model.Side_BUY,
			Quantity:             model.IasD(10),
			Price:                model.IasD(100),
			Listing:              listing1,
			OriginatorId:         evId,
			OriginatorRef:        orderId,
		},
		{
			OrderSide:            model.Side_BUY,
			Quantity:             model.IasD(10),
			Price:                model.IasD(110),
			Listing:              listing2,
			OriginatorId:         evId,
			OriginatorRef:        orderId,
		},
		{
			OrderSide:            model.Side_BUY,
			Quantity:             model.IasD(10),
			Price:                model.IasD(120),
			Listing:              listing1,
			OriginatorId:         evId,
			OriginatorRef:        orderId,
		},
		{
			OrderSide:            model.Side_BUY,
			Quantity:             model.IasD(10),
			Price:                model.IasD(130),
			Listing:              listing2,
			OriginatorId:         evId,
			OriginatorRef:        orderId,
		},
		{
			OrderSide:            model.Side_BUY,
			Quantity:             model.IasD(10),
			Price:                model.IasD(130),
			Listing:              listing1,
			OriginatorId:         evId,
			OriginatorRef:        orderId,
		},

	}

	for idx, params := range client.params {
		if !reflect.DeepEqual(expectedParams[idx], params) {
			t.Fatalf("expected params at idx %v do not match", idx)
		}
	}


}
