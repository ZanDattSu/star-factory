package order

import (
	gRPCClient "order/internal/client/grpc"
	"order/internal/repository"
	srvc "order/internal/service"
)

// Компиляторная проверка: убеждаемся, что *service реализует интерфейс OrderService.
var _ srvc.OrderService = (*service)(nil)

type service struct {
	repository      repository.OrderRepository
	paymentClient   gRPCClient.PaymentClient
	inventoryClient gRPCClient.InventoryClient
}

func NewService(
	repository repository.OrderRepository,
	payClient gRPCClient.PaymentClient,
	invClient gRPCClient.InventoryClient,
) *service {
	return &service{
		repository:      repository,
		paymentClient:   payClient,
		inventoryClient: invClient,
	}
}
