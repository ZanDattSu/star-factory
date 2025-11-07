package model

type OrderStatus string

const (
	OrderStatusUnspecified    OrderStatus = "UNSPECIFIED"
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)
