package service

import (
	"context"

	"github.com/ZanDattSu/star-factory/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userUUID string, partUuids []string) (string, float64, error)
	PayOrder(ctx context.Context, paymentMethod model.PaymentMethod, orderUUID string) (string, error)
	GetOrder(ctx context.Context, orderUUID string) (*model.Order, error)
	CancelOrder(ctx context.Context, orderUUID string) error
}
