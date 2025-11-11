package order

import (
	"math/rand/v2"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"order/internal/model"
)

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
