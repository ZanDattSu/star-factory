package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	inventoryV1 "github.com/ZanDattSu/star-factory/order/internal/client/grpc/inventory/v1"
	"github.com/ZanDattSu/star-factory/order/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

var partsNotFound = "one or more parts not found"

func (s *service) CreateOrder(ctx context.Context, userUUID string, partUuids []string) (string, float64, error) {
	logger.Info(ctx, "Creating new order",
		zap.String("user_uuid", userUUID),
		zap.Int("parts_count", len(partUuids)),
	)

	if len(partUuids) == 0 {
		logger.Warn(ctx, "Failed to create order: empty parts list",
			zap.String("user_uuid", userUUID),
		)
		return "", 0, fmt.Errorf("%s: empty parts list", partsNotFound)
	}

	parts, err := s.inventoryClient.ListParts(
		ctx,
		model.PartsFilter{
			Uuids: partUuids,
		},
	)
	if err != nil {
		notFound := &inventoryV1.PartsNotFoundError{}
		if errors.As(err, &notFound) {
			logger.Error(ctx, "Failed to create order: parts not found",
				zap.String("user_uuid", userUUID),
				zap.Strings("part_uuids", partUuids),
				zap.Error(err),
			)
			return "", 0, fmt.Errorf("%s: %w", partsNotFound, err)
		}
		logger.Error(ctx, "Failed to get parts from inventory",
			zap.String("user_uuid", userUUID),
			zap.Strings("part_uuids", partUuids),
			zap.Error(err),
		)
		return "", 0, err
	}

	if len(parts) != len(partUuids) {
		logger.Error(ctx, "Failed to create order: parts count mismatch",
			zap.String("user_uuid", userUUID),
			zap.Int("requested_parts", len(partUuids)),
			zap.Int("found_parts", len(parts)),
		)
		return "", 0, fmt.Errorf("%s: %w", partsNotFound, err)
	}

	var totalPrice float64
	for _, part := range parts {
		totalPrice += part.Price
	}

	orderUUID := uuid.New().String()
	newOrder := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUuids:  partUuids,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPENDINGPAYMENT,
	}

	err = s.repository.PutOrder(ctx, orderUUID, newOrder)
	if err != nil {
		logger.Error(ctx, "Failed to save order to repository",
			zap.String("order_uuid", orderUUID),
			zap.String("user_uuid", userUUID),
			zap.Error(err),
		)
		return "", 0, fmt.Errorf("failed to put order in repository: %w", err)
	}

	logger.Info(ctx, "Order created successfully",
		zap.String("order_uuid", orderUUID),
		zap.String("user_uuid", userUUID),
		zap.Float64("total_price", totalPrice),
		zap.Int("parts_count", len(partUuids)),
	)

	return orderUUID, totalPrice, nil
}
