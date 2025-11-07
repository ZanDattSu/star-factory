package v1

import (
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

type client struct {
	genClient paymentV1.PaymentServiceClient
}

func NewClient(genClient paymentV1.PaymentServiceClient) *client {
	return &client{genClient: genClient}
}
