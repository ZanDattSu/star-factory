package order

import (
	"fmt"

	"order/internal/model"
)

func (s *SuiteService) TestGetOrderSuccess() {
	expectedOrder := RandomOrder()

	s.orderRepository.
		On("GetOrder", s.ctx, expectedOrder.OrderUUID).
		Return(expectedOrder, nil).
		Once()

	order, err := s.service.GetOrder(s.ctx, expectedOrder.OrderUUID)
	s.Require().NoError(err)
	s.Require().Equal(expectedOrder, order)
}

func (s *SuiteService) TestGetOrderNotFound() {
	expectedOrder := RandomOrder()

	s.orderRepository.
		On("GetOrder", s.ctx, expectedOrder.OrderUUID).
		Return(nil, &model.OrderNotFoundError{}).
		Once()

	order, err := s.service.GetOrder(s.ctx, expectedOrder.OrderUUID)

	s.Require().Error(err)
	s.Require().Nil(order)

	var notFound *model.OrderNotFoundError
	s.Require().ErrorAs(err, &notFound)
	s.Require().Equal(expectedOrder.OrderUUID, notFound.OrderUUID)
	s.Require().Equal(404, notFound.Code)
	s.Require().Contains(err.Error(), fmt.Sprintf("order with UUID %q not found", expectedOrder.OrderUUID))
}
