package order

import (
	"errors"
	"slices"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"order/internal/model"
)

func (s *SuiteService) TestPayOrderSuccess() {
	order := &model.Order{
		OrderUUID:  gofakeit.UUID(),
		UserUUID:   gofakeit.UUID(),
		PartUuids:  RandomPartUuids(),
		TotalPrice: gofakeit.Price(100, 1000),
		Status:     model.OrderStatusPendingPayment,
	}
	paymentMethod := RandomPaymentMethod()
	expectedTransactionUUID := gofakeit.UUID()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).
		Return(order, true).Once()

	s.orderRepository.On("PutOrder",
		s.ctx,
		order.OrderUUID,
		mock.MatchedBy(func(o *model.Order) bool {
			return o.Status == model.OrderStatusPaid &&
				o.PaymentMethod == paymentMethod &&
				o.TransactionUUID != nil &&
				*o.TransactionUUID == expectedTransactionUUID &&
				o.TotalPrice == order.TotalPrice &&
				o.UserUUID == order.UserUUID &&
				slices.Equal(o.PartUuids, order.PartUuids)
		}),
	).Once()

	s.paymentClient.On("PayOrder", s.ctx, order.OrderUUID, order.UserUUID, paymentMethod).
		Return(expectedTransactionUUID, nil).Once()

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod, order.OrderUUID)

	s.Require().NoError(err)
	s.Require().Equal(expectedTransactionUUID, transactionUUID)
}

func (s *SuiteService) TestPayOrderOrderNotFound() {
	orderUUID := gofakeit.UUID()
	paymentMethod := RandomPaymentMethod()

	s.orderRepository.
		On("GetOrder", s.ctx, orderUUID).
		Return((*model.Order)(nil), false).
		Once()

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod, orderUUID)

	s.Require().Error(err)
	s.Require().Empty(transactionUUID)

	var notFound *model.OrderNotFoundError
	s.Require().ErrorAs(err, &notFound)
	s.Require().Equal(orderUUID, notFound.OrderUUID)

	s.orderRepository.AssertNotCalled(s.T(), "PutOrder", mock.Anything)
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder", mock.Anything)
}

func (s *SuiteService) TestPayOrderPaymentFailed() {
	order := &model.Order{
		OrderUUID:  gofakeit.UUID(),
		UserUUID:   gofakeit.UUID(),
		PartUuids:  RandomPartUuids(),
		TotalPrice: gofakeit.Price(100, 1000),
		Status:     model.OrderStatusPendingPayment,
	}
	paymentMethod := RandomPaymentMethod()

	internalErr := status.Error(codes.Internal, "database timeout")

	s.orderRepository.
		On("GetOrder", s.ctx, order.OrderUUID).
		Return(order, true).Once()

	s.paymentClient.
		On("PayOrder", s.ctx, order.OrderUUID, order.UserUUID, paymentMethod).
		Return("", internalErr).Once()

	transactionUUID, err := s.service.PayOrder(s.ctx, paymentMethod, order.OrderUUID)

	s.Require().Error(err)
	s.Require().Empty(transactionUUID)
	s.Require().Contains(err.Error(), "failed to pay order")
	s.Require().Contains(err.Error(), "database timeout")

	s.orderRepository.AssertNotCalled(s.T(), "PutOrder", mock.Anything)
}

func (s *SuiteService) TestPayOrderPaymentFailedNoMutation() {
	order := &model.Order{
		OrderUUID:       gofakeit.UUID(),
		UserUUID:        gofakeit.UUID(),
		Status:          model.OrderStatusPendingPayment,
		PaymentMethod:   "",
		TransactionUUID: nil,
	}

	initialCopy := *order // копируем состояние
	paymentMethod := RandomPaymentMethod()

	s.orderRepository.On("GetOrder", s.ctx, order.OrderUUID).
		Return(order, true).Once()

	s.paymentClient.On("PayOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", errors.New("failed")).Once()

	_, _ = s.service.PayOrder(s.ctx, paymentMethod, order.OrderUUID)

	s.Require().Equal(initialCopy, *order, "order should remain unchanged after failed payment")
}
