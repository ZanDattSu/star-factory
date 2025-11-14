package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ZanDattSu/star-factory/order/internal/client/converter"
	"github.com/ZanDattSu/star-factory/order/internal/model"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod model.PaymentMethod) (string, error) {
	transactionUUID, err := c.genClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: converter.PaymentMethodToProto(paymentMethod),
	})
	if err != nil {
		statusCode, ok := status.FromError(err)
		if ok && statusCode.Code() == codes.Internal {
			return "", fmt.Errorf("payment service internal error: %w", err)
		}
		return "", err
	}

	return transactionUUID.TransactionUuid, nil
}
