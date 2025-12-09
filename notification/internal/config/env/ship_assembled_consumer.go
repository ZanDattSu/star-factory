package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type shipAssembledConsumerEnvConfig struct {
	Topic   string `env:"SHIP_ASSEMBLED_TOPIC_NAME,required"`
	GroupID string `env:"SHIP_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type shipAssembledConsumerConfig struct {
	raw shipAssembledConsumerEnvConfig
}

func NewOrderAssembledConsumerConfig() (*shipAssembledConsumerConfig, error) {
	var raw shipAssembledConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &shipAssembledConsumerConfig{raw: raw}, nil
}

func (cfg *shipAssembledConsumerConfig) Topic() string {
	return cfg.raw.Topic
}

func (cfg *shipAssembledConsumerConfig) GroupID() string {
	return cfg.raw.GroupID
}

func (cfg *shipAssembledConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
