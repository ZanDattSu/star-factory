package order_consumer

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/platform/pkg/kafka/consumer"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

func (s *service) handleShipAssembled(ctx context.Context, msg consumer.Message) error {
	logger.Info(ctx, "Processing message from ship.assembled topic",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
	)

	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode Ship Assembled event",
			zap.String("topic", msg.Topic),
			zap.Int32("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
			zap.Error(err),
		)
		return err
	}

	if event.OrderUUID == "" {
		logger.Error(ctx, "Invalid event: empty order_uuid",
			zap.String("topic", msg.Topic),
			zap.Int32("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
			zap.String("event_uuid", event.EventUUID),
		)
		return errors.New("invalid event")
	}

	logger.Info(ctx, "Received ShipAssembled event",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int("build_time_sec", int(event.BuildTimeSec.Seconds())),
	)

	err = s.notificationService.SendAssembledNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send assembly telegram notification", zap.Error(err))
		return err
	}

	logger.Info(ctx, "ShipAssembled event processed successfully",
		zap.String("order_uuid", event.OrderUUID),
	)

	return nil
}
