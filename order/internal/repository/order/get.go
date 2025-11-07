package order

import (
	"context"

	"order/internal/model"
	"order/internal/repository/converter"
)

func (r *repository) GetOrder(_ context.Context, uuid string) (*model.Order, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[uuid]

	return converter.OrderToModel(order), ok
}
