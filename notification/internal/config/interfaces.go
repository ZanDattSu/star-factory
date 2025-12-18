package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type ShipAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type TelegramBotConfig interface {
	Token() string
	MaxRetries() int
	RetryDelay() time.Duration
}

type AuthGRPCService interface {
	AuthServiceAddress() string
	AuthServicePort() string
}
