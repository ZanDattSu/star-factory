package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/ZanDattSu/star-factory/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, paymentMethod model.PaymentMethod, orderUUID string) (string, error) {
	order, err := s.repository.GetOrder(ctx, orderUUID)
	if err != nil {
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

	order.Status = model.OrderStatusPAID
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = paymentMethod

	err = s.repository.UpdateOrder(ctx, orderUUID, order)
	if err != nil {
		return "", fmt.Errorf("failed to put order in repository: %w", err)
	}

	err = s.orderProducerService.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		EventUuid:       uuid.NewString(),
		OrderUuid:       order.OrderUUID,
		UserUuid:        order.UserUUID,
		PaymentMethod:   order.PaymentMethod,
		TransactionUuid: transactionUUID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to produce order: %w", err)
	}

	return transactionUUID, nil
}
