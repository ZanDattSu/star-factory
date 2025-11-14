package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	inventoryV1 "github.com/ZanDattSu/star-factory/order/internal/client/grpc/inventory/v1"
	"github.com/ZanDattSu/star-factory/order/internal/model"
)

var partsNotFound = "one or more parts not found"

func (s *service) CreateOrder(ctx context.Context, userUUID string, partUuids []string) (string, float64, error) {
	if len(partUuids) == 0 {
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
			return "", 0, fmt.Errorf("%s: %w", partsNotFound, err)
		}
		return "", 0, err
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

	err = s.repository.PutOrder(ctx, orderUUID, newOrder)
	if err != nil {
		return "", 0, fmt.Errorf("failed to put order in repository: %w", err)
	}

	return orderUUID, totalPrice, nil
}
