package payment

import (
	"context"
	"log"

	"payment/internal/model"

	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUuid := a.service.PayOrder(ctx, req.OrderUuid, req.UserUuid, model.PaymentMethod(req.PaymentMethod))
	log.Printf("Оплата прошла успешно, transaction_uuid:%s", transactionUuid)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}
