package order

import (
	"context"
	"errors"
	"fmt"

	orderV1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
	"order/internal/converter"
	"order/internal/model"
)

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	if params.OrderUUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "order UUID should be not empty",
		}, nil
	}

	order, err := a.orderService.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		notFound := &model.OrderNotFoundError{}
		if errors.As(err, &notFound) {
			return nil, fmt.Errorf("order not found error: %w", err)
		}
		return nil, fmt.Errorf("get order error: %w", err)
	}

	return converter.OrderToAPI(order), nil
}
