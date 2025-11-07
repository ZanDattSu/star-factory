package order

import (
	"context"

	"order/internal/model"
	"order/internal/repository/converter"
)

func (r *repository) PutOrder(_ context.Context, uuid string, order *model.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[uuid] = converter.OrderToRepoModel(order)
}
