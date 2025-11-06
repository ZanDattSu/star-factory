package service

import (
	"context"

	"payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod model.PaymentMethod) string
}
