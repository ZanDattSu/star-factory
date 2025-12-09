package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/ZanDattSu/star-factory/order/internal/converter/kafka"
	"github.com/ZanDattSu/star-factory/order/internal/repository"
	serv "github.com/ZanDattSu/star-factory/order/internal/service"
	"github.com/ZanDattSu/star-factory/platform/pkg/kafka"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

type service struct {
	shipAssembledConsumer kafka.Consumer
	shipAssembledDecoder  kafkaConverter.ShipAssembledDecoder
	orderService          serv.OrderService
	orderRepository       repository.OrderRepository
}

func NewService(
	shipAssembledConsumer kafka.Consumer,
	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder,
	orderService serv.OrderService,
	orderRepository repository.OrderRepository,
) *service {
	return &service{
		shipAssembledConsumer: shipAssembledConsumer,
		shipAssembledDecoder:  shipAssembledDecoder,
		orderService:          orderService,
		orderRepository:       orderRepository,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order consumer for ship.assembled topic")

	err := s.shipAssembledConsumer.Consume(ctx, s.handleShipAssembled)
	if err != nil {
		logger.Error(ctx, "Failed to consume from ship.assembled topic", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order consumer stopped")
	return nil
}
