package grpc

import (
	"context"

	"github.com/ZanDattSu/star-factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, partsFilter model.PartsFilter) ([]*model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod model.PaymentMethod) (string, error)
}
