package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
