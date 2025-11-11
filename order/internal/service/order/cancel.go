package order

import (
	"context"
	"fmt"

	"order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID string) error {
	order, err := s.repository.GetOrder(ctx, orderUUID)
	if err != nil {
		return model.NewOrderNotFoundError(orderUUID)
	}

	switch order.Status {
	case model.OrderStatusPendingPayment:
		order.Status = model.OrderStatusCancelled
		err = s.repository.UpdateOrder(ctx, order.OrderUUID, order)
		if err != nil {
			return fmt.Errorf("failed to put order in repository: %w", err)
		}
		return nil
	case model.OrderStatusPaid:
		return model.NewConflictError("cannot cancel a paid order")
	case model.OrderStatusCancelled:
		return model.NewConflictError("cannot cancel a canceled order")
	}

	return nil
}
