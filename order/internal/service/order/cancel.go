package order

import (
	"context"

	"order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID string) error {
	order, ok := s.repository.GetOrder(ctx, orderUUID)
	if !ok {
		return model.NewOrderNotFoundError(orderUUID)
	}

	switch order.Status {
	case model.OrderStatusPendingPayment:
		order.Status = model.OrderStatusCancelled
		s.repository.PutOrder(ctx, order.OrderUUID, order)
		return nil
	case model.OrderStatusPaid:
		return model.NewConflictError("cannot cancel a paid order")
	case model.OrderStatusCancelled:
		return model.NewConflictError("cannot cancel a canceled order")
	}

	return nil
}
