package app

import (
	"context"

	payApi "payment/internal/api/v1/payment"
	"payment/internal/service"
	payService "payment/internal/service/payment"

	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1Api   paymentV1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1Api(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1Api == nil {
		d.paymentV1Api = payApi.NewApi(d.PaymentService(ctx))
	}

	return d.paymentV1Api
}

func (d *diContainer) PaymentService(_ context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = payService.NewService()
	}

	return d.paymentService
}
