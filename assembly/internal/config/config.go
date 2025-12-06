package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/ZanDattSu/star-factory/assembly/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	Kafka         KafkaConfig
	OrderConsumer AssemblyConsumerConfig
	OrderProducer AssemblyProducerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	consumerCfg, err := env.NewOrderConsumerConfig()
	if err != nil {
		return err
	}

	producerCfg, err := env.NewOrderProduceConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		Kafka:         kafkaCfg,
		OrderConsumer: consumerCfg,
		OrderProducer: producerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
