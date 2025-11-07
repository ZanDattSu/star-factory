package order

import (
	gRPCClient "order/internal/client/grpc"
	inventoryService "order/internal/client/grpc/inventory/v1"
	paymentService "order/internal/client/grpc/payment/v1"
	"order/internal/repository"
	srvc "order/internal/service"

	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
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
	payClient paymentV1.PaymentServiceClient,
	invClient inventoryV1.InventoryServiceClient,
) *service {
	return &service{
		repository:      repository,
		paymentClient:   paymentService.NewClient(payClient),
		inventoryClient: inventoryService.NewClient(invClient),
	}
}
