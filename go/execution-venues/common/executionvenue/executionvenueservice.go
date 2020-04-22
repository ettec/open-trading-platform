package executionvenue

import (
	"context"
	"fmt"
	api "github.com/ettec/open-trading-platform/go/common/api/executionvenue"
	"github.com/ettec/open-trading-platform/go/execution-venues/common/ordermanager"
	"github.com/ettec/open-trading-platform/go/model"
	"log"
)

type ExecVenueService struct {
	orderManager ordermanager.OrderManager
}

func New(om ordermanager.OrderManager) *ExecVenueService {
	service := ExecVenueService{orderManager: om}
	return &service
}

func (s *ExecVenueService) CreateAndRouteOrder(context context.Context, params *api.CreateAndRouteOrderParams) (*api.OrderId, error) {

	log.Printf("Received  order parameters-> %v", params)

	if params.GetQuantity() == nil {
		return nil, fmt.Errorf("quantity required on params:%v", params)
	}

	if params.GetPrice() == nil {
		return nil, fmt.Errorf("price required on params:%v", params)
	}

	if params.GetListing() == nil {
		return nil, fmt.Errorf("listing required on params:%v", params)
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

func (s *ExecVenueService) CancelOrder(ctx context.Context, p *api.CancelOrderParams) (*model.Empty, error) {
	return &model.Empty{}, s.orderManager.CancelOrder(p)
}

func (s *ExecVenueService) Close() {
	if s.orderManager != nil {
		s.orderManager.Close()
	}
}
