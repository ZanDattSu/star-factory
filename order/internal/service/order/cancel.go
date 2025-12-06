package order

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/order/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID string) error {
	logger.Info(ctx, "Cancelling order",
		zap.String("order_uuid", orderUUID),
	)

	order, err := s.repository.GetOrder(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx, "Failed to get order for cancellation",
			zap.String("order_uuid", orderUUID),
			zap.Error(err),
		)
		return model.NewOrderNotFoundError(orderUUID)
	}

	logger.Debug(ctx, "Order found, checking status",
		zap.String("order_uuid", orderUUID),
		zap.String("current_status", string(order.Status)),
	)

	switch order.Status {
	case model.OrderStatusPENDINGPAYMENT:
		order.Status = model.OrderStatusCANCELLED
		err = s.repository.UpdateOrder(ctx, order.OrderUUID, order)
		if err != nil {
			logger.Error(ctx, "Failed to update order status to cancelled",
				zap.String("order_uuid", orderUUID),
				zap.Error(err),
			)
			return fmt.Errorf("failed to update order status to cancelled: %w", err)
		}
		logger.Info(ctx, "Order cancelled successfully",
			zap.String("order_uuid", orderUUID),
			zap.String("user_uuid", order.UserUUID),
		)
		return nil
	case model.OrderStatusPAID:
		logger.Warn(ctx, "Cannot cancel paid order",
			zap.String("order_uuid", orderUUID),
			zap.String("status", string(order.Status)),
		)
		return model.NewConflictError("cannot cancel a paid order")
	case model.OrderStatusCANCELLED:
		logger.Warn(ctx, "Order already cancelled",
			zap.String("order_uuid", orderUUID),
		)
		return model.NewConflictError("cannot cancel a canceled order")
	}

	return nil
}
