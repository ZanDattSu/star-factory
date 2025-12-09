package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/ZanDattSu/star-factory/notification/internal/converter/kafka"
	serv "github.com/ZanDattSu/star-factory/notification/internal/service"
	"github.com/ZanDattSu/star-factory/platform/pkg/kafka"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

type service struct {
	orderPaidConsumer   kafka.Consumer
	orderPaidDecoder    kafkaConverter.OrderPaidDecoder
	notificationService serv.NotificationService
}

func NewService(
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	notificationService serv.NotificationService,
) *service {
	return &service{
		orderPaidConsumer:   orderPaidConsumer,
		orderPaidDecoder:    orderPaidDecoder,
		notificationService: notificationService,
	}
}

func (s *service) RunOrderPaidConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting consumer for order.paid topic")

	err := s.orderPaidConsumer.Consume(ctx, s.handleOrderPaid)
	if err != nil {
		logger.Error(ctx, "Failed to consume from order.paid topic", zap.Error(err))
		return err
	}

	logger.Info(ctx, "order.paid consumer stopped")
	return nil
}
