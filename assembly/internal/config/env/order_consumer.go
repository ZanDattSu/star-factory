package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderConsumerEnvConfig struct {
	Topic   string `env:"CONSUME_TOPIC_NAME,required"`
	GroupID string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

type orderConsumerConfig struct {
	raw orderConsumerEnvConfig
}

func NewOrderConsumerConfig() (*orderConsumerConfig, error) {
	var raw orderConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderConsumerConfig{raw: raw}, nil
}

func (cfg *orderConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *orderConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *orderConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
