package model

import "time"

// OrderPaidEvent - событие "заказ оплачен" (приходит от Order Service)
type OrderPaidEvent struct {
	EventUuid       string
	OrderUuid       string
	UserUuid        string
	PaymentMethod   PaymentMethod
	TransactionUuid string
}

// ShipAssembledEvent - событие "корабль собран"
type ShipAssembledEvent struct {
	EventUuid string
	OrderUuid string
	UserUuid  string
	BuildTime time.Duration
}
