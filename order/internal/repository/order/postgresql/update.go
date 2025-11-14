package postgresql

import (
	"context"
	"fmt"

	"github.com/ZanDattSu/star-factory/order/internal/model"
)

func (r *repository) UpdateOrder(ctx context.Context, _ string, order *model.Order) error {
	paymentMethodID, err := order.PaymentMethod.ID()
	if err != nil {
		return fmt.Errorf("invalid payment method: %w", err)
	}

	statusID, err := order.Status.ID()
	if err != nil {
		return fmt.Errorf("invalid order status: %w", err)
	}

	const query = `
		UPDATE orders o
		SET user_uuid = ($2),
		    part_uuids = ($3),
		    total_price = ($4),
		    transaction_uuid = ($5),
		    payment_method_id = ($6),
		    status_id = ($7)
		WHERE order_uuid = ($1)
	`

	cmdTag, err := r.pool.Exec(ctx, query,
		order.OrderUUID,
		order.UserUUID,
		order.PartUuids,
		order.TotalPrice,
		order.TransactionUUID,
		paymentMethodID,
		statusID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order %s: %w", order.OrderUUID, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("order %s not found", order.OrderUUID)
	}

	return nil
}
