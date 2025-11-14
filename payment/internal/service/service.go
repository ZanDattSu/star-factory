package service

import (
	"context"

	"github.com/ZanDattSu/star-factory/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod model.PaymentMethod) string
}
