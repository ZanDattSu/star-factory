package kafka

import (
	"context"

	"github.com/ZanDattSu/star-factory/platform/pkg/kafka/consumer"
)

// MessageHandler — обработчик сообщений.
type MessageHandler func(ctx context.Context, msg consumer.Message) error

type Consumer interface {
	Consume(ctx context.Context, handler consumer.MessageHandler) error
}

type Producer interface {
	Send(ctx context.Context, key, value []byte) error
}
