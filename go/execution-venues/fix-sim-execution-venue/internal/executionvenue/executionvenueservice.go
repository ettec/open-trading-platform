package executionvenue

import (
	"context"
	"fmt"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/model"
	"log/slog"
)

type orderManager interface {
	CancelOrder(id *api.CancelOrderParams) error
	CreateAndRouteOrder(params *api.CreateAndRouteOrderParams) (*api.OrderId, error)
	ModifyOrder(params *api.ModifyOrderParams) error
	SetOrderStatus(orderId string, status model.OrderStatus) error
	SetErrorMsg(orderId string, msg string) error
	AddExecution(orderId string, lastPrice model.Decimal64, lastQty model.Decimal64, execId string) error
}

type ExecVenueService struct {
	orderManager orderManager
}

func New(om orderManager) *ExecVenueService {
	service := ExecVenueService{orderManager: om}
	return &service
}

func (s *ExecVenueService) CreateAndRouteOrder(_ context.Context, params *api.CreateAndRouteOrderParams) (*api.OrderId, error) {

	slog.Info("Received  order parameters", "params", params)

	if params.GetQuantity() == nil {
		return nil, fmt.Errorf("quantity required on params:%v", params)
	}

	if params.GetPrice() == nil {
		return nil, fmt.Errorf("price required on params:%v", params)
	}

	if params.GetListingId() == 0 {
		return nil, fmt.Errorf("listing id required on params:%v", params)
	}

	if params.GetOriginatorId() == "" {
		return nil, fmt.Errorf("originator id required on params:%v", params)
	}

	if params.GetOriginatorRef() == "" {
		return nil, fmt.Errorf("originator ref required on params:%v", params)
	}

	result, err := s.orderManager.CreateAndRouteOrder(params)
	if err != nil {
		return nil, fmt.Errorf("error when creating and routing order:%w", err)
	}

	slog.Info("created order", "orderId", result.OrderId)

	return &api.OrderId{
		OrderId: result.OrderId,
	}, nil
}

func (s *ExecVenueService) CancelOrder(_ context.Context, p *api.CancelOrderParams) (*model.Empty, error) {
	if err := s.orderManager.CancelOrder(p); err != nil {
		return nil, fmt.Errorf("error when cancelling order:%w", err)
	}

	return &model.Empty{}, nil
}

func (s *ExecVenueService) ModifyOrder(_ context.Context, params *api.ModifyOrderParams) (*model.Empty, error) {

	if err := s.orderManager.ModifyOrder(params); err != nil {
		return nil, fmt.Errorf("error when modifying order:%w", err)
	}

	return &model.Empty{}, nil
}

func (s *ExecVenueService) GetExecutionParametersMetaData(context.Context, *model.Empty) (*api.ExecParamsMetaDataJson, error) {
	return &api.ExecParamsMetaDataJson{}, nil
}
