package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type orderHttpEnvConfig struct {
	HttpHost              string        `env:"HTTP_HOST,required"`
	HttpPort              string        `env:"HTTP_PORT,required"`
	InventoryGRPCHost     string        `env:"INVENTORY_GRPC_HOST,required"`
	InventoryGRPCPort     string        `env:"INVENTORY_GRPC_PORT,required"`
	PaymentGRPCHost       string        `env:"PAYMENT_GRPC_HOST,required"`
	PaymentGRPCPort       string        `env:"PAYMENT_GRPC_PORT,required"`
	HttpReadHeaderTimeout time.Duration `env:"HTTP_READ_TIMEOUT,required"`
	HttpShutdownTimeout   time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT,required"`
}

type orderHttpConfig struct {
	raw orderHttpEnvConfig
}

func NewOrderHttpConfig() (*orderHttpConfig, error) {
	var raw orderHttpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &orderHttpConfig{raw: raw}, nil
}

func (cfg *orderHttpConfig) OrderAddress() string {
	return net.JoinHostPort(cfg.raw.HttpHost, cfg.raw.HttpPort)
}

func (cfg *orderHttpConfig) OrderPort() string {
	return cfg.raw.HttpPort
}

func (cfg *orderHttpConfig) PaymentAddress() string {
	return net.JoinHostPort(cfg.raw.PaymentGRPCHost, cfg.raw.PaymentGRPCPort)
}

func (cfg *orderHttpConfig) PaymentServicePort() string {
	return cfg.raw.PaymentGRPCPort
}

func (cfg *orderHttpConfig) InventoryAddress() string {
	return net.JoinHostPort(cfg.raw.InventoryGRPCHost, cfg.raw.InventoryGRPCPort)
}

func (cfg *orderHttpConfig) InventoryServicePort() string {
	return cfg.raw.InventoryGRPCPort
}

func (cfg *orderHttpConfig) ReadHeaderTimeout() time.Duration {
	return cfg.raw.HttpReadHeaderTimeout
}

func (cfg *orderHttpConfig) ShutdownTimeout() time.Duration {
	return cfg.raw.HttpShutdownTimeout
}
