package converter

import (
	"github.com/ZanDattSu/star-factory/order/internal/model"
	repoModel "github.com/ZanDattSu/star-factory/order/internal/repository/model"
)

// OrderToRepoModel конвертирует *model.Order → *repoModel.Order.
func OrderToRepoModel(o *model.Order) *repoModel.Order {
	if o == nil {
		return nil
	}
	return &repoModel.Order{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartUuids:       o.PartUuids,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   repoModel.PaymentMethod(o.PaymentMethod),
		Status:          repoModel.OrderStatus(o.Status),
	}
}

// OrderToModel конвертирует *repoModel.Order → *model.Order.
func OrderToModel(o *repoModel.Order) *model.Order {
	if o == nil {
		return nil
	}
	return &model.Order{
		OrderUUID:       o.OrderUUID,
		UserUUID:        o.UserUUID,
		PartUuids:       o.PartUuids,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   model.PaymentMethod(o.PaymentMethod),
		Status:          model.OrderStatus(o.Status),
	}
}

// OrdersToRepoModel конвертирует []*model.Order → []*repoModel.Order.
func OrdersToRepoModel(orders []*model.Order) []*repoModel.Order {
	out := make([]*repoModel.Order, 0, len(orders))
	for _, o := range orders {
		out = append(out, OrderToRepoModel(o))
	}
	return out
}

// OrdersToModel конвертирует []*repoModel.Order → []*model.Order.
func OrdersToModel(orders []*repoModel.Order) []*model.Order {
	out := make([]*model.Order, 0, len(orders))
	for _, o := range orders {
		out = append(out, OrderToModel(o))
	}
	return out
}
