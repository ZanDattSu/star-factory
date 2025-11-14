package order

import (
	"context"
	"net/http"

	"github.com/ZanDattSu/star-factory/order/internal/service"
	orderV1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
)

type api struct {
	orderService service.OrderService
}

func NewApi(orderService service.OrderService) *api {
	return &api{orderService: orderService}
}

func (a *api) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}
