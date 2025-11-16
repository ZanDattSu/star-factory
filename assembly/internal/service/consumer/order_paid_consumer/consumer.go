package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/ZanDattSu/star-factory/assembly/internal/converter/kafka"
	serv "github.com/ZanDattSu/star-factory/assembly/internal/service"
	"github.com/ZanDattSu/star-factory/platform/pkg/kafka"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
	assemblyService   serv.AssemblyService
}

func NewService(
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	assemblyService serv.AssemblyService,
) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		assemblyService:   assemblyService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting OrderPaid consumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.handleOrderPaid)
	if err != nil {
		logger.Error(ctx, "Failed to consume from order.paid topic", zap.Error(err))
		return err
	}

	return nil
}
