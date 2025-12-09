package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/ZanDattSu/star-factory/notification/internal/converter/kafka"
	serv "github.com/ZanDattSu/star-factory/notification/internal/service"
	"github.com/ZanDattSu/star-factory/platform/pkg/kafka"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

type service struct {
	shipAssembledConsumer kafka.Consumer
	shipAssembledDecoder  kafkaConverter.ShipAssembledDecoder
	notificationService   serv.NotificationService
}

func NewService(
	shipAssembledConsumer kafka.Consumer,
	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder,
	notificationService serv.NotificationService,
) *service {
	return &service{
		shipAssembledConsumer: shipAssembledConsumer,
		shipAssembledDecoder:  shipAssembledDecoder,
		notificationService:   notificationService,
	}
}

func (s *service) RunShipAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting consumer for ship.assembled topic")

	err := s.shipAssembledConsumer.Consume(ctx, s.handleShipAssembled)
	if err != nil {
		logger.Error(ctx, "Failed to consume from ship.assembled topic", zap.Error(err))
		return err
	}

	logger.Info(ctx, "ship.assembled consumer stopped")
	return nil
}
