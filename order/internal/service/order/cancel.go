package order

import (
	"context"
	"fmt"

	"github.com/ZanDattSu/star-factory/order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID string) error {
	order, err := s.repository.GetOrder(ctx, orderUUID)
	if err != nil {
		return model.NewOrderNotFoundError(orderUUID)
	}

	switch order.Status {
	case model.OrderStatusPENDINGPAYMENT:
		order.Status = model.OrderStatusCANCELLED
		err = s.repository.UpdateOrder(ctx, order.OrderUUID, order)
		if err != nil {
			return fmt.Errorf("failed to put order in repository: %w", err)
		}
		return nil
	case model.OrderStatusPAID:
		return model.NewConflictError("cannot cancel a paid order")
	case model.OrderStatusCANCELLED:
		return model.NewConflictError("cannot cancel a canceled order")
	}

	return nil
}
