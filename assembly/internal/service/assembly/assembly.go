package assembly

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/assembly/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

func (s *service) ProcessOrderPaid(ctx context.Context, event *model.OrderPaidEvent) error {
	logger.Info(ctx, "Starting ship assembly",
		zap.String("order_uuid", event.OrderUuid),
		zap.String("user_uuid", event.UserUuid),
	)

	randomSeconds, err := rand.Int(rand.Reader, big.NewInt(10))
	if err != nil {
		logger.Error(ctx, "Failed to generate random build time", zap.Error(err))
		return err
	}
	buildTime := time.Duration(randomSeconds.Int64()+1) * time.Second

	logger.Info(ctx, "Assembling ship...",
		zap.String("order_uuid", event.OrderUuid),
		zap.Int("build_time_sec", int(buildTime.Seconds())),
	)

	timer := time.NewTimer(buildTime)
	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
	}

	logger.Info(ctx, "Ship assembled successfully",
		zap.String("order_uuid", event.OrderUuid),
		zap.Int("build_time_sec", int(buildTime.Seconds())),
	)

	shipAssembledEvent := &model.ShipAssembledEvent{
		EventUuid: uuid.New().String(),
		OrderUuid: event.OrderUuid,
		UserUuid:  event.UserUuid,
		BuildTime: buildTime,
	}

	if err = s.shipAssembledProducer.PublishShipAssembled(ctx, shipAssembledEvent); err != nil {
		logger.Error(ctx, "Failed to publish ShipAssembled event", zap.Error(err))
		return err
	}

	return nil
}
