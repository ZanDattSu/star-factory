package order

import (
	"context"
	"fmt"

	"order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, paymentMethod model.PaymentMethod, orderUUID string) (string, error) {
	order, ok := s.repository.GetOrder(ctx, orderUUID)
	if !ok {
		return "", model.NewOrderNotFoundError(orderUUID)
	}

	transactionUUID, err := s.paymentClient.PayOrder(
		ctx,
		order.OrderUUID,
		order.UserUUID,
		paymentMethod,
	)
	if err != nil {
		return "", fmt.Errorf("failed to pay order %s: %w", orderUUID, err)
	}

	order.Status = model.OrderStatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = paymentMethod

	s.repository.PutOrder(ctx, orderUUID, order)

	return transactionUUID, nil
}
