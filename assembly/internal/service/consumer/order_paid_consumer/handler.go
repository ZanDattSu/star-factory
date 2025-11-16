package order_paid_consumer

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/platform/pkg/kafka/consumer"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

func (s *service) handleOrderPaid(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid event", zap.Error(err))
		return err
	}

	if event.OrderUuid == "" {
		logger.Error(ctx, "Invalid event: empty order_uuid")
		return errors.New("invalid event")
	}

	logger.Info(ctx, "Received OrderPaid event",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUuid),
		zap.String("order_uuid", event.OrderUuid),
		zap.String("payment_method", string(event.PaymentMethod)),
		zap.String("transaction_uuid", event.TransactionUuid),
	)

	err = s.assemblyService.ProcessOrderPaid(ctx, &event)
	if err != nil {
		logger.Error(ctx, "Failed to process OrderPaid event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "OrderPaid event processed successfully",
		zap.String("order_uuid", event.OrderUuid),
	)

	return nil
}
