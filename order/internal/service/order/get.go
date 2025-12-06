package order

import (
	"context"

	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/order/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

func (s *service) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	logger.Debug(ctx, "Getting order",
		zap.String("order_uuid", orderUUID),
	)

	order, err := s.repository.GetOrder(ctx, orderUUID)
	if err != nil {
		logger.Warn(ctx, "Order not found",
			zap.String("order_uuid", orderUUID),
			zap.Error(err),
		)
		return nil, model.NewOrderNotFoundError(orderUUID)
	}

	logger.Debug(ctx, "Order retrieved successfully",
		zap.String("order_uuid", orderUUID),
		zap.String("user_uuid", order.UserUUID),
		zap.String("status", string(order.Status)),
	)

	return order, nil
}
