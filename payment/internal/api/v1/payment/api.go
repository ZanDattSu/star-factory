package payment

import (
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
	"payment/internal/service"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer
	service service.PaymentService
}

func NewApi(service service.PaymentService) *api {
	return &api{
		service: service,
	}
}
