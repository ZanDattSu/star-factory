package config

import (
	"github.com/joho/godotenv"
	"inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	InventoryGRPC InventoryGRPCConfig
	Mongo         MongoConfig
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

	inventory, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	mongo, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        logger,
		InventoryGRPC: inventory,
		Mongo:         mongo,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
