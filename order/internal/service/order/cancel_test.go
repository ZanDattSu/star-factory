package order

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/ZanDattSu/star-factory/order/internal/model"
)

func (s *SuiteService) TestCancelOrderSuccess() {
	order := RandomOrder()
	order.Status = model.OrderStatusPENDINGPAYMENT

	s.orderRepository.
		On("GetOrder", s.ctx, order.OrderUUID).
		Return(order, nil).Once()

	s.orderRepository.
		On("UpdateOrder",
			s.ctx,
			order.OrderUUID,
			mock.MatchedBy(func(o *model.Order) bool {
				return o.Status == model.OrderStatusCANCELLED
			}),
		).Return(nil).Once()

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)
	s.Require().NoError(err)
}

func (s *SuiteService) TestCancelOrderNotFound() {
	order := RandomOrder()

	s.orderRepository.
		On("GetOrder", s.ctx, order.OrderUUID).
		Return(nil, &model.OrderNotFoundError{}).
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
	order.Status = model.OrderStatusPAID

	s.orderRepository.
		On("GetOrder", s.ctx, order.OrderUUID).
		Return(order, nil).
		Once()

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)

	s.Require().Error(err)

	var conflict *model.ConflictError
	s.Require().ErrorAs(err, &conflict)
	s.Require().Equal(409, conflict.Code)
	s.Require().Contains(conflict.Error(), "cannot cancel a paid order")
}

func (s *SuiteService) TestCancelOrderConflictAlreadyCancelled() {
	order := RandomOrder()
	order.Status = model.OrderStatusCANCELLED

	s.orderRepository.
		On("GetOrder", s.ctx, order.OrderUUID).
		Return(order, nil).
		Once()

	err := s.service.CancelOrder(s.ctx, order.OrderUUID)

	s.Require().Error(err)

	var conflict *model.ConflictError
	s.Require().ErrorAs(err, &conflict)
	s.Require().Equal(409, conflict.Code)
	s.Require().Contains(conflict.Error(), "cannot cancel a canceled order")
}
