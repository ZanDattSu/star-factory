package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	inventoryV1 "order/internal/client/grpc/inventory/v1"
	"order/internal/model"
)

var partsNotFound = "one or more parts not found"

func (s *service) CreateOrder(ctx context.Context, userUUID string, partUuids []string) (string, float64, error) {
	parts, err := s.inventoryClient.ListParts(
		ctx,
		model.PartsFilter{
			Uuids: partUuids,
		},
	)
	if err != nil {
		partNotFound := &inventoryV1.PartsNotFoundError{}
		if errors.As(err, &partNotFound) {
			return "", 0, fmt.Errorf("%s: %w", partsNotFound, err)
		}
		return "", 0, fmt.Errorf("failed to list parts from inventory service: %w", err)
	}

	if len(parts) != len(partUuids) {
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
		Status:     model.OrderStatusPendingPayment,
	}

	s.repository.PutOrder(ctx, orderUUID, newOrder)

	return orderUUID, totalPrice, nil
}
