package config

import (
	"github.com/joho/godotenv"

	"github.com/ZanDattSu/star-factory/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	PaymentGRPC PaymentGRPCConfig
	Auth        AuthGRPCService
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	paymentGrpc, err := env.NewPaymentGrpcConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:      logger,
		PaymentGRPC: paymentGrpc,
		Auth:        paymentGrpc,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
