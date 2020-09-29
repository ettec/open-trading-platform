package executionvenue

import (
	"context"
	"fmt"
	api "github.com/ettec/otp-common/api/executionvenue"
	"github.com/ettec/otp-common/model"
	"log"
)

type orderManager interface {
	CancelOrder(id *api.CancelOrderParams) error
	CreateAndRouteOrder(params *api.CreateAndRouteOrderParams) (*api.OrderId, error)
	ModifyOrder(params *api.ModifyOrderParams) error
	SetOrderStatus(orderId string, status model.OrderStatus) error
	SetErrorMsg(orderId string, msg string) error
	AddExecution(orderId string, lastPrice model.Decimal64, lastQty model.Decimal64, execId string) error
	Close()
}

type ExecVenueService struct {
	orderManager orderManager
}

func New(om orderManager) *ExecVenueService {
	service := ExecVenueService{orderManager: om}
	return &service
}

func (s *ExecVenueService) CreateAndRouteOrder(_ context.Context, params *api.CreateAndRouteOrderParams) (*api.OrderId, error) {

	log.Printf("Received  order parameters-> %v", params)

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
		log.Printf("error when creating and routing order:%v", err)
		return nil, err
	}

	log.Printf("created order id:%v", result.OrderId)

	return &api.OrderId{
		OrderId: result.OrderId,
	}, nil
}

func (s *ExecVenueService) CancelOrder(_ context.Context, p *api.CancelOrderParams) (*model.Empty, error) {
	err := s.orderManager.CancelOrder(p)
	if err != nil {
		return nil, err
	}

	return &model.Empty{}, nil
}

func (s *ExecVenueService) ModifyOrder(_ context.Context, params *api.ModifyOrderParams) (*model.Empty, error) {

	err := s.orderManager.ModifyOrder(params)
	if err != nil {
		return nil, err
	}

	return &model.Empty{}, nil
}

func (s *ExecVenueService) GetExecutionParametersMetaData(context.Context, *model.Empty) (*api.ExecParamsMetaDataJson, error) {
	return &api.ExecParamsMetaDataJson{}, nil
}

func (s *ExecVenueService) Close() {
	if s.orderManager != nil {
		s.orderManager.Close()
	}
}
