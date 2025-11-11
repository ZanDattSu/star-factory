package order

import (
	"context"

	"order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, orderUUID string) (*model.Order, error) {
	order, err := s.repository.GetOrder(ctx, orderUUID)
	if err != nil {
		return nil, model.NewOrderNotFoundError(orderUUID)
	}

	return order, nil
}
