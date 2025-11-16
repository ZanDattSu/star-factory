package kafka

import "github.com/ZanDattSu/star-factory/assembly/internal/model"

// OrderPaidDecoder - декодер для OrderPaidEvent события
type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
