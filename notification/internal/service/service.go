package service

import (
	"context"
	"github.com/ZanDattSu/star-factory/notification/internal/model"
)

// NotificationService - отправка уведомления в телеграмм
type NotificationService interface {
	SendPaidNotification(ctx context.Context, paidEvent model.OrderPaidEvent) error
	SendAssembledNotification(ctx context.Context, shipAssembledEvent model.ShipAssembledEvent) error
}

// OrderPaidConsumerService - слушает "order.paid" топик
type OrderPaidConsumerService interface {
	RunOrderPaidConsumer(ctx context.Context) error
}

// ShipAssembledConsumerService - слушает "ship.assembled" топик
type ShipAssembledConsumerService interface {
	RunShipAssembledConsumer(ctx context.Context) error
}
