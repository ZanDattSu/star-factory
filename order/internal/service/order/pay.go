package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/order/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

func (s *service) PayOrder(ctx context.Context, paymentMethod model.PaymentMethod, orderUUID string) (string, error) {
	logger.Info(ctx, "Processing order payment",
		zap.String("order_uuid", orderUUID),
		zap.String("payment_method", string(paymentMethod)),
	)

	order, err := s.repository.GetOrder(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx, "Failed to get order for payment",
			zap.String("order_uuid", orderUUID),
			zap.Error(err),
		)
		return "", model.NewOrderNotFoundError(orderUUID)
	}

	logger.Debug(ctx, "Order found, initiating payment",
		zap.String("order_uuid", orderUUID),
		zap.String("user_uuid", order.UserUUID),
		zap.String("order_status", string(order.Status)),
	)

	transactionUUID, err := s.paymentClient.PayOrder(
		ctx,
		order.OrderUUID,
		order.UserUUID,
		paymentMethod,
	)
	if err != nil {
		logger.Error(ctx, "Payment failed",
			zap.String("order_uuid", orderUUID),
			zap.String("user_uuid", order.UserUUID),
			zap.String("payment_method", string(paymentMethod)),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to pay order %s: %w", orderUUID, err)
	}

	logger.Info(ctx, "Payment successful, updating order status",
		zap.String("order_uuid", orderUUID),
		zap.String("transaction_uuid", transactionUUID),
	)

	order.Status = model.OrderStatusPAID
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = paymentMethod

	err = s.repository.UpdateOrder(ctx, orderUUID, order)
	if err != nil {
		logger.Error(ctx, "Failed to update order status after payment",
			zap.String("order_uuid", orderUUID),
			zap.String("transaction_uuid", transactionUUID),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to put order in repository: %w", err)
	}

	err = s.orderProducerService.ProduceOrderPaid(ctx,
		model.OrderPaidEvent{
			EventUuid:       uuid.NewString(),
			OrderUuid:       order.OrderUUID,
			UserUuid:        order.UserUUID,
			PaymentMethod:   order.PaymentMethod,
			TransactionUuid: transactionUUID,
		},
	)
	if err != nil {
		logger.Error(ctx, "Failed to produce order paid event",
			zap.String("order_uuid", orderUUID),
			zap.String("transaction_uuid", transactionUUID),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to produce order: %w", err)
	}

	logger.Info(ctx, "Order payment completed successfully",
		zap.String("order_uuid", orderUUID),
		zap.String("transaction_uuid", transactionUUID),
		zap.String("user_uuid", order.UserUUID),
	)

	return transactionUUID, nil
}
