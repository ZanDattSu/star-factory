package kafka

import "github.com/ZanDattSu/star-factory/notification/internal/model"

// OrderPaidDecoder - декодер для OrderPaidEvent события
type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}

type ShipAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
