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
	id, err := event.PaymentMethod.ID()
	if err != nil {
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
		logger.Error(ctx, "failed to marshal UFORecorded", zap.Error(err))
		return err
	}

	err = s.orderPaidProducer.Send(ctx, []byte(event.OrderUuid), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish UFORecorded", zap.Error(err))
		return err
	}

	return nil
}
