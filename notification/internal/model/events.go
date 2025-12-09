package model

import "time"

// OrderPaidEvent - событие "заказ оплачен" (приходит от Order Service)
type OrderPaidEvent struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   PaymentMethod
	TransactionUUID string
}

// ShipAssembledEvent - событие "корабль собран"
type ShipAssembledEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec time.Duration
}
