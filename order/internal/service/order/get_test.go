package order

import (
	"fmt"
	"math/rand/v2"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"order/internal/model"
)

func (s *SuiteService) TestGetOrderSuccess() {
	expectedOrder := RandomOrder()

	s.orderRepository.
		On("GetOrder", s.ctx, expectedOrder.OrderUUID).
		Return(expectedOrder, true).
		Once()

	order, err := s.service.GetOrder(s.ctx, expectedOrder.OrderUUID)
	s.Require().NoError(err)
	s.Require().Equal(expectedOrder, order)
}

func (s *SuiteService) TestGetOrderNotFound() {
	expectedOrder := RandomOrder()

	s.orderRepository.
		On("GetOrder", s.ctx, expectedOrder.OrderUUID).
		Return(nil, false).
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

func RandomOrder() *model.Order {
	return &model.Order{
		OrderUUID:       gofakeit.UUID(),
		UserUUID:        gofakeit.UUID(),
		PartUuids:       RandomPartUuids(),
		TotalPrice:      gofakeit.Price(100, 1000),
		TransactionUUID: lo.ToPtr(gofakeit.UUID()),
		PaymentMethod:   RandomPaymentMethod(),
		Status:          RandomOrderStatus(),
	}
}

func RandomPartUuids() []string {
	countParts := 1 + rand.IntN(9) // [1 - 10]
	partUuids := make([]string, countParts)

	for i := 0; i < countParts; i++ {
		partUuids[i] = gofakeit.UUID()
	}

	return partUuids
}

func RandomPaymentMethod() model.PaymentMethod {
	paymentMethods := []model.PaymentMethod{
		model.PaymentMethodCard,
		model.PaymentMethodCreditCard,
		model.PaymentMethodInvestorMoney,
		model.PaymentMethodSbp,
	}
	return paymentMethods[rand.IntN(len(paymentMethods))]
}

func RandomOrderStatus() model.OrderStatus {
	statuses := []model.OrderStatus{
		model.OrderStatusPendingPayment,
		model.OrderStatusPaid,
		model.OrderStatusCancelled,
	}
	return statuses[rand.IntN(len(statuses))]
}
