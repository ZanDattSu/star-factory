package order

import (
	"github.com/samber/lo"
	"order/internal/model"
)

func (s *SuiteRepository) TestPutAndGetOrderSuccess() {
	order := &model.Order{
		OrderUUID:       "order-123",
		UserUUID:        "user-999",
		PartUuids:       []string{"part-1", "part-2"},
		TotalPrice:      2500.50,
		PaymentMethod:   model.PaymentMethodCard,
		Status:          model.OrderStatusPaid,
		TransactionUUID: lo.ToPtr("txn-777"),
	}

	s.repo.PutOrder(s.ctx, order.OrderUUID, order)

	got, ok := s.repo.GetOrder(s.ctx, order.OrderUUID)

	s.Require().True(ok, "ожидалось, что заказ будет найден")
	s.Require().NotNil(got)

	s.Equal(order.OrderUUID, got.OrderUUID)
	s.Equal(order.UserUUID, got.UserUUID)
	s.Equal(order.TotalPrice, got.TotalPrice)
	s.Equal(order.PaymentMethod, got.PaymentMethod)
	s.Equal(order.Status, got.Status)
	s.NotNil(got.TransactionUUID)
	s.Equal("txn-777", *got.TransactionUUID)
}

func (s *SuiteRepository) TestGetOrderNotFound() {
	got, ok := s.repo.GetOrder(s.ctx, "nonexistent-order")
	s.False(ok)
	s.Nil(got)
}
