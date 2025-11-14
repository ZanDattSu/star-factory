package order

import (
	"context"

	orderV1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
)

func (a *api) Health(_ context.Context) (orderV1.HealthRes, error) {
	return &orderV1.HealthRequest{
		Status:  "ok",
		Service: "order-service",
	}, nil
}
