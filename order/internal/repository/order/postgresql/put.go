package postgresql

import (
	"context"
	"fmt"

	"github.com/ZanDattSu/star-factory/order/internal/model"
)

func (r *repository) PutOrder(ctx context.Context, _ string, order *model.Order) error {
	paymentMethodID, err := order.PaymentMethod.ID()
	if err != nil {
		return fmt.Errorf("invalid payment method: %w", err)
	}

	statusID, err := order.Status.ID()
	if err != nil {
		return fmt.Errorf("invalid order status: %w", err)
	}

	const query = `
		INSERT INTO orders(order_uuid,
		                   user_uuid,
		                   part_uuids,
		                   total_price,
		                   transaction_uuid,
		                   payment_method_id,
		                   status_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.pool.Exec(ctx, query,
		order.OrderUUID,
		order.UserUUID,
		order.PartUuids,
		order.TotalPrice,
		order.TransactionUUID,
		paymentMethodID,
		statusID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order %s: %w", order.OrderUUID, err)
	}

	return nil
}
