package inmemory

import (
	"context"

	"github.com/ZanDattSu/star-factory/order/internal/model"
	"github.com/ZanDattSu/star-factory/order/internal/repository/converter"
)

func (r *repository) PutOrder(_ context.Context, uuid string, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[uuid] = converter.OrderToRepoModel(order)
	return nil
}
