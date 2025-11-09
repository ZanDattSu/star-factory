package order

import (
	"fmt"

	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

func (s *SuiteService) TestCancelOrderSuccess() {
	order := RandomOrder()
	order.Status = model.OrderStatusPendingPayment

	s.orderRepository.
		On("GetOrder", s.ctx, order.OrderUUID).
		Return(order, true).Once()

	s.orderRepository.
		On("PutOrder",
			s.ctx,
			order.OrderUUID,
			mock.MatchedBy(func(o *model.Order) bool {
				return o.Status == model.OrderStatusCancelled
			}),
		).Once()

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.Require().NoError(err)
}

func (s *SuiteService) TestCancelOrderNotFound() {
	order := RandomOrder()

	s.orderRepository.
		On("GetOrder", s.ctx, order.OrderUUID).
		Return(nil, false).
		Once()

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)

	s.Require().Error(err)

	var notFound *model.OrderNotFoundError
	s.Require().ErrorAs(err, &notFound)
	s.Require().Equal(order.OrderUUID, notFound.OrderUUID)
	s.Require().Equal(404, notFound.Code)
	s.Require().Contains(err.Error(), fmt.Sprintf("order with UUID %q not found", order.OrderUUID))
}

func (s *SuiteService) TestCancelOrderConflictPaid() {
	order := RandomOrder()
	order.Status = model.OrderStatusPaid

	s.orderRepository.
		On("GetOrder", s.ctx, order.OrderUUID).
		Return(order, true).
		Once()

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)

	s.Require().Error(err)

	var conflict *model.ConflictError
	s.Require().ErrorAs(err, &conflict)
	s.Require().Equal(409, conflict.Code)
	s.Require().Contains(conflict.Error(), "cannot cancel a paid order")
}

func (s *SuiteService) TestCancelOrder_ConflictAlreadyCancelled() {
	order := RandomOrder()
	order.Status = model.OrderStatusCancelled

	s.orderRepository.
		On("GetOrder", s.ctx, order.OrderUUID).
		Return(order, true).
		Once()

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)

	s.Require().Error(err)

	var conflict *model.ConflictError
	s.Require().ErrorAs(err, &conflict)
	s.Require().Equal(409, conflict.Code)
	s.Require().Contains(conflict.Error(), "cannot cancel a canceled order")
}
