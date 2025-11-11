package inmemory

import (
	"context"
	"fmt"

	"order/internal/model"
	"order/internal/repository/converter"
)

func (r *repository) GetOrder(_ context.Context, uuid string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[uuid]

	if !ok {
		return nil, fmt.Errorf("order with UUID %s not found", uuid)
	}

	return converter.OrderToModel(order), nil
}
