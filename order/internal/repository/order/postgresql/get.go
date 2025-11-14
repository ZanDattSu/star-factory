package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/ZanDattSu/star-factory/order/internal/model"
)

func (r *repository) GetOrder(ctx context.Context, uuid string) (*model.Order, error) {
	const query = `
		SELECT
			o.order_uuid,
			o.user_uuid,
			o.part_uuids,
			o.total_price,
			o.transaction_uuid,
			o.payment_method_id,
			o.status_id
		FROM orders o
		WHERE o.order_uuid = $1
	`

	var (
		order           model.Order
		paymentMethodID int
		statusID        int
	)
	err := r.pool.QueryRow(ctx, query, uuid).Scan(
		&order.OrderUUID,
		&order.UserUUID,
		&order.PartUuids,
		&order.TotalPrice,
		&order.TransactionUUID,
		&paymentMethodID,
		&statusID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.NewOrderNotFoundError(uuid)
		}
		return nil, err
	}

	order.PaymentMethod, _ = model.PaymentMethodFromID(paymentMethodID) //nolint:gosec
	order.Status, _ = model.OrderStatusFromID(statusID)                 //nolint:gosec

	return &order, nil
}
