package main

import (
	"context"
	"github.com/ettec/otp-common/api/executionvenue"
)

type service struct {
}

func (v service) CreateAndRouteOrder(ctx context.Context, params *executionvenue.CreateAndRouteOrderParams) (*executionvenue.OrderId, error) {
	panic("implement me")
}

func (v service) CancelOrder(ctx context.Context, params *executionvenue.CancelOrderParams) (*interface{}, error) {
	panic("implement me")
}

func (v service) ModifyOrder(ctx context.Context, params *executionvenue.ModifyOrderParams) (*interface{}, error) {
	panic("implement me")
}

func (v service) GetExecutionParametersMetaData(ctx context.Context, m *interface{}) (*executionvenue.ExecParamsMetaDataJson, error) {
	panic("implement me")
}

func main() {

}
