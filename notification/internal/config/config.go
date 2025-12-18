package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/ZanDattSu/star-factory/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                LoggerConfig
	Kafka                 KafkaConfig
	OrderPaidConsumer     OrderPaidConsumerConfig
	ShipAssembledConsumer ShipAssembledConsumerConfig
	TelegramBot           TelegramBotConfig
	AuthService           AuthGRPCService
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

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	shipAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	telegramBotCfg, err := env.NewTelegramBotConfig()
	if err != nil {
		return err
	}

	authGrpcConfig, err := env.NewAuthGrpcConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                loggerCfg,
		Kafka:                 kafkaCfg,
		OrderPaidConsumer:     orderPaidConsumerCfg,
		ShipAssembledConsumer: shipAssembledConsumerCfg,
		TelegramBot:           telegramBotCfg,
		AuthService:           authGrpcConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
