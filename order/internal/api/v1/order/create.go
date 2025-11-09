package order

import (
	"context"
	"errors"
	"fmt"

	orderV1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "order/internal/client/grpc/inventory/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if req.UserUUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "user UUID should be not empty",
		}, nil
	}

	if len(req.PartUuids) == 0 {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "parts UUIDs should be contains at least 1 part",
		}, nil
	}

	orderUUID, totalPrice, err := a.orderService.CreateOrder(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		partNotFound := &inventoryV1.PartsNotFoundError{}
		if errors.As(err, &partNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("one or more parts not found: %s", err),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("inventory service internal error: %v", err),
		}, nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}
