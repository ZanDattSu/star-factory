package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type authGrpcEnvConfig struct {
	AuthGRPCHost string `env:"AUTH_GRPC_HOST,required"`
	AuthGRPCPort string `env:"AUTH_GRPC_PORT,required"`
}

type authGrpcConfig struct {
	raw authGrpcEnvConfig
}

func NewAuthGrpcConfig() (*authGrpcConfig, error) {
	var raw authGrpcEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &authGrpcConfig{raw: raw}, nil
}

func (cfg *authGrpcConfig) AuthServicePort() string {
	return cfg.raw.AuthGRPCPort
}

func (cfg *authGrpcConfig) AuthServiceAddress() string {
	return net.JoinHostPort(cfg.raw.AuthGRPCHost, cfg.raw.AuthGRPCPort)
}
