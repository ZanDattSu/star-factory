package order

import (
	gRPCClient "github.com/ZanDattSu/star-factory/order/internal/client/grpc"
	"github.com/ZanDattSu/star-factory/order/internal/repository"
	srvc "github.com/ZanDattSu/star-factory/order/internal/service"
)

// Компиляторная проверка: убеждаемся, что *service реализует интерфейс OrderService.
var _ srvc.OrderService = (*service)(nil)

type service struct {
	repository           repository.OrderRepository
	paymentClient        gRPCClient.PaymentClient
	inventoryClient      gRPCClient.InventoryClient
	orderProducerService srvc.OrderProducerService
}

func NewService(
	repository repository.OrderRepository,
	payClient gRPCClient.PaymentClient,
	invClient gRPCClient.InventoryClient,
	orderProducerService srvc.OrderProducerService,
) *service {
	return &service{
		repository:           repository,
		paymentClient:        payClient,
		inventoryClient:      invClient,
		orderProducerService: orderProducerService,
	}
}
