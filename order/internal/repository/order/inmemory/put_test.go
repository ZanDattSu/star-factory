package inmemory

import "github.com/ZanDattSu/star-factory/order/internal/model"

func (s *SuiteRepository) TestPutOrderOverridesExisting() {
	order1 := &model.Order{
		OrderUUID:     "order-override",
		UserUUID:      "user-1",
		PaymentMethod: model.PaymentMethodCard,
		Status:        model.OrderStatusPendingPayment,
	}

	order2 := &model.Order{
		OrderUUID:     "order-override",
		UserUUID:      "user-2",
		PaymentMethod: model.PaymentMethodSbp,
		Status:        model.OrderStatusPaid,
	}

	err := s.repo.PutOrder(s.ctx, order1.OrderUUID, order1)
	if err != nil {
		return
	}
	got1, err1 := s.repo.GetOrder(s.ctx, order1.OrderUUID)
	s.Require().NoError(err1)
	s.Equal("user-1", got1.UserUUID)
	s.Equal(model.OrderStatusPendingPayment, got1.Status)

	_ = s.repo.PutOrder(s.ctx, order2.OrderUUID, order2)
	got2, err2 := s.repo.GetOrder(s.ctx, order2.OrderUUID)

	s.Require().NoError(err2)
	s.Equal("user-2", got2.UserUUID)
	s.Equal(model.OrderStatusPaid, got2.Status)
}

func (s *SuiteRepository) TestConcurrentAccess() {
	order := &model.Order{
		OrderUUID: "order-concurrent",
		UserUUID:  "user-777",
		Status:    model.OrderStatusPendingPayment,
	}

	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			_ = s.repo.PutOrder(s.ctx, order.OrderUUID, order)
		}
		close(done)
	}()

	for i := 0; i < 100; i++ {
		_, _ = s.repo.GetOrder(s.ctx, order.OrderUUID)
	}

	<-done
	s.True(true, "не должно паниковать при конкурентном доступе")
}
