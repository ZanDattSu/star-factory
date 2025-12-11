package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryGRPCEnvConfig struct {
	Host         string `env:"GRPC_HOST,required"`
	GRPCPort     string `env:"GRPC_PORT,required"`
	HTTPPort     string `env:"HTTP_GATEWAY_PORT,required"`
	AuthGRPCHost string `env:"AUTH_GRPC_HOST,required"`
	AuthGRPCPort string `env:"AUTH_GRPC_PORT,required"`
}

type inventoryGRPCConfig struct {
	raw inventoryGRPCEnvConfig
}

func NewInventoryGRPCConfig() (*inventoryGRPCConfig, error) {
	var raw inventoryGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &inventoryGRPCConfig{raw: raw}, nil
}

func (cfg *inventoryGRPCConfig) GRPCAddress() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.GRPCPort)
}

func (cfg *inventoryGRPCConfig) GRPCPort() string {
	return cfg.raw.GRPCPort
}

func (cfg *inventoryGRPCConfig) HTTPAddress() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.HTTPPort)
}

func (cfg *inventoryGRPCConfig) HTTPPort() string {
	return cfg.raw.HTTPPort
}

func (cfg *inventoryGRPCConfig) AuthServicePort() string {
	return cfg.raw.AuthGRPCPort
}

func (cfg *inventoryGRPCConfig) AuthServiceAddress() string {
	return net.JoinHostPort(cfg.raw.AuthGRPCHost, cfg.raw.AuthGRPCPort)
}
