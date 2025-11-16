package model

import "time"

type OrderPaidEvent struct {
	EventUuid       string
	OrderUuid       string
	UserUuid        string
	PaymentMethod   PaymentMethod
	TransactionUuid string
}

type ShipAssembledEvent struct {
	EventUuid string
	OrderUuid string
	UserUuid  string
	BuildTime time.Duration
}
