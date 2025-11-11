package inmemory

import (
	"context"

	"order/internal/model"
	"order/internal/repository/converter"
)

func (r *repository) PutOrder(_ context.Context, uuid string, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[uuid] = converter.OrderToRepoModel(order)
	return nil
}
