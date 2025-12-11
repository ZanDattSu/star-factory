package config

import (
	"github.com/joho/godotenv"

	"github.com/ZanDattSu/star-factory/auth/internal/config/env"
)

var appConfig *config

type config struct {
	App      App
	Logger   LoggerConfig
	GRPC     GRPCConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Session  SessionConfig
}

func Load(path ...string) error {
	if err := godotenv.Load(path...); err != nil {
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

	grpcCfg, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	postgres, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	redis, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	session, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		App:      app,
		Logger:   logger,
		GRPC:     grpcCfg,
		Postgres: postgres,
		Redis:    redis,
		Session:  session,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
