package service

import (
	"context"

	"github.com/ZanDattSu/star-factory/assembly/internal/model"
)

// AssemblyService - бизнес-логика сборки корабля
type AssemblyService interface {
	ProcessOrderPaid(ctx context.Context, event *model.OrderPaidEvent) error
}

// OrderPaidConsumerService - слушает "order.paid" топик
type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

// ShipAssembledProducerService - отправляет в "ship.assembled" топик
type ShipAssembledProducerService interface {
	PublishShipAssembled(ctx context.Context, event *model.ShipAssembledEvent) error
}
