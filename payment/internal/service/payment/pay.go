package payment

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/payment/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) string {
	logger.Info(ctx, "Processing payment",
		zap.String("order_uuid", orderUUID),
		zap.String("user_uuid", userUUID),
		zap.String("payment_method", string(paymentMethod)),
	)

	transactionUUID := uuid.New().String()

	logger.Info(ctx, "Payment transaction created",
		zap.String("order_uuid", orderUUID),
		zap.String("user_uuid", userUUID),
		zap.String("transaction_uuid", transactionUUID),
		zap.String("payment_method", string(paymentMethod)),
	)

	return transactionUUID
}
