package payment

import (
	"github.com/ZanDattSu/star-factory/payment/internal/service"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
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
