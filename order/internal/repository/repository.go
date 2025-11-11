package repository

import (
	"context"

	"order/internal/model"
)

type OrderRepository interface {
	GetOrder(ctx context.Context, uuid string) (*model.Order, error)
	PutOrder(ctx context.Context, uuid string, order *model.Order) error
	UpdateOrder(ctx context.Context, uuid string, order *model.Order) error
}
