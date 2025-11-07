package order

import (
	"context"
	"errors"
	"fmt"

	"order/internal/model"

	orderV1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	if params.OrderUUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "order UUID should be not empty",
		}, nil
	}

	err := a.orderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {

		notFound := &model.OrderNotFoundError{}
		conflict := &model.ConflictError{}
		switch {
		case errors.As(err, &notFound):
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("CancelOrder err: %s", err),
			}, nil
		case errors.As(err, &conflict):
			return &orderV1.ConflictError{
				Code:    409,
				Message: err.Error(),
			}, nil

		default:
			return &orderV1.InternalServerError{
				Code:    500,
				Message: fmt.Sprintf("internal server error: %s", err),
			}, nil
		}
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
