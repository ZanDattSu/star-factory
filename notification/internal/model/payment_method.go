package model

import "fmt"

type PaymentMethod string

const (
	PaymentMethodUnspecified   PaymentMethod = "Неизвестный метод"
	PaymentMethodCard          PaymentMethod = "Банковская карта"
	PaymentMethodSbp           PaymentMethod = "Система быстрых платежей"
	PaymentMethodCreditCard    PaymentMethod = "Кредитная карта" //nolint:gosec
	PaymentMethodInvestorMoney PaymentMethod = "Деньги инвестора"
)

// Маппинг код <-> ID
var paymentMethodToID = map[PaymentMethod]int{
	PaymentMethodUnspecified:   1,
	PaymentMethodCard:          2,
	PaymentMethodSbp:           3,
	PaymentMethodCreditCard:    4,
	PaymentMethodInvestorMoney: 5,
}

var idToPaymentMethod = map[int]PaymentMethod{
	1: PaymentMethodUnspecified,
	2: PaymentMethodCard,
	3: PaymentMethodSbp,
	4: PaymentMethodCreditCard,
	5: PaymentMethodInvestorMoney,
}

func (m PaymentMethod) ID() (int, error) {
	id, ok := paymentMethodToID[m]
	if !ok {
		return 0, fmt.Errorf("unknown payment method: %s", m)
	}
	return id, nil
}

func PaymentMethodFromID(id int) (PaymentMethod, error) {
	m, ok := idToPaymentMethod[id]
	if !ok {
		return PaymentMethodUnspecified, fmt.Errorf("unknown payment method id: %d", id)
	}
	return m, nil
}
