package config

import (
	"time"

	"github.com/IBM/sarama"
)

type App interface {
	ShutdownTimeout() time.Duration
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OrderHTTPConfig interface {
	OrderAddress() string
	OrderPort() string
	ReadHeaderTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

type PaymentGRPCService interface {
	PaymentAddress() string
	PaymentServicePort() string
}

type InventoryGrpcService interface {
	InventoryAddress() string
	InventoryServicePort() string
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
	MigrationsPath() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type AssemblyConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
