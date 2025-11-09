package payment

import (
	"github.com/brianvoe/gofakeit/v7"
	"payment/internal/model"
)

func (s *ServiceSuite) TestPaySuccess() {
	var (
		orderUuid     = gofakeit.UUID()
		userUuid      = gofakeit.UUID()
		paymentMethod = randomPaymentMethod()
	)

	transactionUuid := s.service.PayOrder(s.ctx, orderUuid, userUuid, paymentMethod)

	s.Require().NotNil(transactionUuid)
	s.IsType("", transactionUuid, "uuid должен быть строкой")
}

func randomPaymentMethod() model.PaymentMethod {
	methods := []model.PaymentMethod{
		model.PaymentMethodCard,
		model.PaymentMethodSbp,
		model.PaymentMethodCreditCard,
		model.PaymentMethodInvestorMoney,
	}
	return methods[gofakeit.Number(0, len(methods)-1)]
}
