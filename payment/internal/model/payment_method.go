package model

type PaymentMethod int

const (
	PaymentMethodUnspecified PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSbp
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

func (c PaymentMethod) String() string {
	return [...]string{
		"Unspecified",
		"Card",
		"SBP",
		"CreditCard",
		"InvestorMoney",
	}[c]
}
