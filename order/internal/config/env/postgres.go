package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host           string `env:"POSTGRES_HOST" envDefault:"POSTGRES_HOST,required"`
	Port           string `env:"POSTGRES_PORT,required"`
	Database       string `env:"POSTGRES_DB,required"`
	User           string `env:"POSTGRES_USER,required"`
	Password       string `env:"POSTGRES_PASSWORD,required"`
	SslMode        string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
	MigrationsPath string `env:"MIGRATION_DIRECTORY,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &postgresConfig{raw: raw}, nil
}

func (cfg *postgresConfig) URI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.Database,
		cfg.raw.SslMode,
	)
}

func (cfg *postgresConfig) DatabaseName() string {
	return cfg.raw.Database
}

func (cfg *postgresConfig) MigrationsPath() string {
	return cfg.raw.MigrationsPath
}
