package order_consumer

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/order/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/kafka/consumer"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

func (s *service) orderHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled event")
		return err
	}

	if event.OrderUuid == "" {
		logger.Error(ctx, "Invalid event: empty order_uuid")
		return errors.New("invalid event")
	}

	logger.Info(ctx, "Received OrderPaidEvent",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUuid),
		zap.String("order_uuid", event.OrderUuid),
		zap.String("user_uuid", event.UserUuid),
		zap.Int("build_time_sec", int(event.BuildTime.Seconds())),
	)

	order, err := s.orderRepository.GetOrder(ctx, event.OrderUuid)
	if err != nil {
		logger.Error(ctx, "Failed to get order", zap.Error(err))
		return err
	}

	order.Status = model.OrderStatusASSEMBLED

	err = s.orderRepository.UpdateOrder(ctx, order.OrderUUID, order)
	if err != nil {
		logger.Error(ctx, "Failed to update order status to ASSEMBLED", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order status updated to ASSEMBLED",
		zap.String("order_uuid", event.OrderUuid),
	)

	return nil
}
