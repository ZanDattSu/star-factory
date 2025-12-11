package v1

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ZanDattSu/star-factory/order/internal/client/converter"
	"github.com/ZanDattSu/star-factory/order/internal/model"
	grpcAuth "github.com/ZanDattSu/star-factory/platform/pkg/grpc/interceptor"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod model.PaymentMethod) (string, error) {
	logger.Info(ctx, "Requesting payment from payment service",
		zap.String("order_uuid", orderUuid),
		zap.String("user_uuid", userUuid),
		zap.String("payment_method", string(paymentMethod)),
	)

	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)

	transactionUUID, err := c.genClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: converter.PaymentMethodToProto(paymentMethod),
	})
	if err != nil {
		statusCode, ok := status.FromError(err)
		if ok && statusCode.Code() == codes.Internal {
			logger.Error(ctx, "Payment service internal error",
				zap.String("order_uuid", orderUuid),
				zap.String("user_uuid", userUuid),
				zap.String("payment_method", string(paymentMethod)),
				zap.String("grpc_code", statusCode.Code().String()),
				zap.Error(err),
			)
			return "", fmt.Errorf("payment service internal error: %w", err)
		}

		logger.Error(ctx, "Payment request failed",
			zap.String("order_uuid", orderUuid),
			zap.String("user_uuid", userUuid),
			zap.String("payment_method", string(paymentMethod)),
			zap.Error(err),
		)
		return "", err
	}

	logger.Info(ctx, "Payment successful",
		zap.String("order_uuid", orderUuid),
		zap.String("user_uuid", userUuid),
		zap.String("transaction_uuid", transactionUUID.TransactionUuid),
		zap.String("payment_method", string(paymentMethod)),
	)

	return transactionUUID.TransactionUuid, nil
}
