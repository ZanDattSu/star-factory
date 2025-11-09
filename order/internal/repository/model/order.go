package model

type Order struct {
	OrderUUID       string        `json:"order_uuid"`
	UserUUID        string        `json:"user_uuid"`
	PartUuids       []string      `json:"part_uuids"`
	TotalPrice      float64       `json:"total_price"`
	TransactionUUID *string       `json:"transaction_uuid"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	Status          OrderStatus   `json:"status"`
}

type OrderStatus string

const (
	OrderStatusUnspecified    OrderStatus = "UNSPECIFIED"
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodUnspecified   PaymentMethod = "UNSPECIFIED"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSbp           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)
