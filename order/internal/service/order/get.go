package order

import (
	"context"

	"order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	order, ok := s.repository.GetOrder(ctx, orderUUID)
	if !ok {
		return nil, model.NewOrderNotFoundError(orderUUID)
	}

	return order, nil
}
