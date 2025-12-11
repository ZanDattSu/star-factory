package config

import (
	"github.com/joho/godotenv"

	"github.com/ZanDattSu/star-factory/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	App           App
	Logger        LoggerConfig
	InventoryGRPC InventoryGRPCConfig
	InventoryHTTP InventoryHTTPConfig
	Auth          AuthGRPCService
	Mongo         MongoConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	app, err := env.NewAppConfig()
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
		App:           app,
		Logger:        logger,
		InventoryGRPC: inventory,
		InventoryHTTP: inventory,
		Auth:          inventory,
		Mongo:         mongo,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
