package order

import (
	"context"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/ZanDattSu/star-factory/order/internal/client/grpc/mocks"
	"github.com/ZanDattSu/star-factory/order/internal/repository/mocks"
	serviceMocks "github.com/ZanDattSu/star-factory/order/internal/service/mocks"
)

type SuiteService struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	orderRepository      *mocks.OrderRepository
	paymentClient        *clientMocks.PaymentClient
	inventoryClient      *clientMocks.InventoryClient
	orderProducerService *serviceMocks.OrderProducerService

	service *service
}

func (s *SuiteService) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.paymentClient = clientMocks.NewPaymentClient(s.T())
	s.inventoryClient = clientMocks.NewInventoryClient(s.T())
	s.orderProducerService = serviceMocks.NewOrderProducerService(s.T())

	s.service = NewService(
		s.orderRepository,
		s.paymentClient,
		s.inventoryClient,
		s.orderProducerService,
	)
	logger.SetNopLogger()
}

func (s *SuiteService) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(SuiteService))
}
