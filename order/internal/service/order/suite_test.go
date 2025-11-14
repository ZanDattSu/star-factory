package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/ZanDattSu/star-factory/order/internal/client/grpc/mocks"
	"github.com/ZanDattSu/star-factory/order/internal/repository/mocks"
)

type SuiteService struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	orderRepository *mocks.OrderRepository
	paymentClient   *clientMocks.PaymentClient
	inventoryClient *clientMocks.InventoryClient

	service *service
}

func (s *SuiteService) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.paymentClient = clientMocks.NewPaymentClient(s.T())
	s.inventoryClient = clientMocks.NewInventoryClient(s.T())

	s.service = NewService(s.orderRepository, s.paymentClient, s.inventoryClient)
}

func (s *SuiteService) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(SuiteService))
}
