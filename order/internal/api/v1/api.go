package v1

import (
	"context"

	orderV1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
)

type OrderApi interface {
	CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error)
	PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error)
	GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error)
	CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error)
	NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode
}
