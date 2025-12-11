package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type paymentGrpcEnvConfig struct {
	Host                string        `env:"GRPC_HOST,required"`
	GRPCPort            string        `env:"GRPC_PORT,required"`
	HTTPPort            string        `env:"HTTP_GATEWAY_PORT,required"`
	HttpShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT,required"`
	AuthGRPCHost        string        `env:"AUTH_GRPC_HOST,required"`
	AuthGRPCPort        string        `env:"AUTH_GRPC_PORT,required"`
}

type paymentGrpcConfig struct {
	raw paymentGrpcEnvConfig
}

func NewPaymentGrpcConfig() (*paymentGrpcConfig, error) {
	var raw paymentGrpcEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &paymentGrpcConfig{raw: raw}, nil
}

func (cfg *paymentGrpcConfig) GRPCAddress() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.GRPCPort)
}

func (cfg *paymentGrpcConfig) GRPCPort() string {
	return cfg.raw.GRPCPort
}

func (cfg *paymentGrpcConfig) HTTPAddress() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.HTTPPort)
}

func (cfg *paymentGrpcConfig) HTTPPort() string {
	return cfg.raw.HTTPPort
}

func (cfg *paymentGrpcConfig) ShutdownTimeout() time.Duration {
	return cfg.raw.HttpShutdownTimeout
}

func (cfg *paymentGrpcConfig) AuthServicePort() string {
	return cfg.raw.AuthGRPCPort
}

func (cfg *paymentGrpcConfig) AuthServiceAddress() string {
	return net.JoinHostPort(cfg.raw.AuthGRPCHost, cfg.raw.AuthGRPCPort)
}
