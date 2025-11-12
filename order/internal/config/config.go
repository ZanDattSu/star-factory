package config

import (
	"github.com/joho/godotenv"
	"order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger    LoggerConfig
	OrderHttp OrderHttpConfig
	Payment   PaymentGRPCService
	Inventory InventoryGrpcService
	Postgres  PostgresConfig
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

	order, err := env.NewOrderHttpConfig()
	if err != nil {
		return err
	}

	postgres, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:    logger,
		OrderHttp: order,
		Payment:   order,
		Inventory: order,
		Postgres:  postgres,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
