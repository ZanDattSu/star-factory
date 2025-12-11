package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type AppEnvConfig struct {
	Shutdown time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
}

type appConfig struct {
	raw AppEnvConfig
}

func NewAppConfig() (*appConfig, error) {
	var raw AppEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &appConfig{raw: raw}, nil
}

func (a *appConfig) ShutdownTimeout() time.Duration {
	return a.raw.Shutdown
}
