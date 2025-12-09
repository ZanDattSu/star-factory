package kafka

import "github.com/ZanDattSu/star-factory/order/internal/model"

type ShipAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
