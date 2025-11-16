package order

import (
	"context"
	"errors"
	"fmt"

	api2 "github.com/ZanDattSu/star-factory/order/internal/converter/api"
	"github.com/ZanDattSu/star-factory/order/internal/model"
	orderV1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	if req.PaymentMethod == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "payment method should be not empty",
		}, nil
	}

	if err := req.PaymentMethod.Validate(); err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("payment method %s does not exist", req.PaymentMethod),
		}, nil
	}

	if params.OrderUUID == "" {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "empty order UUID",
		}, nil
	}

	transactionUUID, err := a.orderService.PayOrder(ctx, api2.PaymentMethodToModel(req.PaymentMethod), params.OrderUUID)
	if err != nil {
		notFound := &model.OrderNotFoundError{}
		if errors.As(err, &notFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("one or more parts not found: %s", err),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("payment service internal error: %v", err),
		}, nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
