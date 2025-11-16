package ship_assembled_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/ZanDattSu/star-factory/assembly/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/kafka"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	eventsV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/events/v1"
)

type service struct {
	shipAssembledProducer kafka.Producer
}

func NewService(shipAssembledProducer kafka.Producer) *service {
	return &service{
		shipAssembledProducer: shipAssembledProducer,
	}
}

func (s *service) PublishShipAssembled(ctx context.Context, event *model.ShipAssembledEvent) error {
	msg := &eventsV1.ShipAssembledEvent{
		EventUuid:    event.EventUuid,
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: int64(event.BuildTime.Seconds()),
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Failed to marshal ShipAssembled event", zap.Error(err))
		return err
	}

	err = s.shipAssembledProducer.Send(ctx, []byte(event.OrderUuid), payload)
	if err != nil {
		logger.Error(ctx, "Failed to publish ShipAssembled event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "ShipAssembled event published",
		zap.String("event_uuid", event.EventUuid),
		zap.String("order_uuid", event.OrderUuid),
		zap.Int64("build_time_sec", int64(event.BuildTime.Seconds())),
	)

	return nil
}
