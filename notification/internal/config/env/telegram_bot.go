package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type telegramBotEnvConfig struct {
	Token      string        `env:"TELEGRAM_BOT_TOKEN,required"`
	MaxRetries int           `env:"TELEGRAM_BOT_MAX_RETRIES" envDefault:"10"`
	RetryDelay time.Duration `env:"TELEGRAM_BOT_RETRY_DELAY" envDefault:"3s"`
}

type telegramBotConfig struct {
	raw telegramBotEnvConfig
}

func NewTelegramBotConfig() (*telegramBotConfig, error) {
	var raw telegramBotEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &telegramBotConfig{raw: raw}, nil
}

func (cfg *telegramBotConfig) Token() string {
	return cfg.raw.Token
}

func (cfg *telegramBotConfig) MaxRetries() int {
	return cfg.raw.MaxRetries
}

func (cfg *telegramBotConfig) RetryDelay() time.Duration {
	return cfg.raw.RetryDelay
}
