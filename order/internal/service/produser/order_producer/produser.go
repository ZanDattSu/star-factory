package order_producer

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/ZanDattSu/star-factory/order/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/kafka"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	eventsV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/events/v1"
)

type service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{orderPaidProducer: orderPaidProducer}
}

func (s *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	logger.Info(ctx, "Producing OrderPaid event",
		zap.String("event_uuid", event.EventUuid),
		zap.String("order_uuid", event.OrderUuid),
		zap.String("user_uuid", event.UserUuid),
		zap.String("transaction_uuid", event.TransactionUuid),
		zap.String("payment_method", string(event.PaymentMethod)),
	)

	id, err := event.PaymentMethod.ID()
	if err != nil {
		logger.Error(ctx, "Invalid payment method",
			zap.String("event_uuid", event.EventUuid),
			zap.String("order_uuid", event.OrderUuid),
			zap.String("payment_method", string(event.PaymentMethod)),
			zap.Error(err),
		)
		return fmt.Errorf("invalid payment method: %w", err)
	}

	msg := &eventsV1.OrderPaid{
		OrderUuid:       event.OrderUuid,
		UserUuid:        event.UserUuid,
		PaymentMethod:   eventsV1.PaymentMethod(int32(id)), //nolint:gosec // unavoidable proto enum cast
		TransactionUuid: event.TransactionUuid,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Failed to marshal OrderPaid event",
			zap.String("event_uuid", event.EventUuid),
			zap.String("order_uuid", event.OrderUuid),
			zap.Error(err),
		)
		return err
	}

	err = s.orderPaidProducer.Send(ctx, []byte(event.OrderUuid), payload)
	if err != nil {
		logger.Error(ctx, "Failed to publish OrderPaid event",
			zap.String("event_uuid", event.EventUuid),
			zap.String("order_uuid", event.OrderUuid),
			zap.String("transaction_uuid", event.TransactionUuid),
			zap.Error(err),
		)
		return err
	}

	logger.Info(ctx, "OrderPaid event published successfully",
		zap.String("event_uuid", event.EventUuid),
		zap.String("order_uuid", event.OrderUuid),
		zap.String("transaction_uuid", event.TransactionUuid),
	)

	return nil
}
