package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type grpcEnvConfig struct {
	Host            string        `env:"GRPC_HOST,required"`
	Port            string        `env:"GRPC_PORT,required"`
	ShutdownTimeout time.Duration `env:"GRPC_SHUTDOWN_TIMEOUT" envDefault:"5s"`
}

type grpcConfig struct {
	raw grpcEnvConfig
}

func NewGRPCConfig() (*grpcConfig, error) {
	var raw grpcEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &grpcConfig{raw: raw}, nil
}

func (g *grpcConfig) Address() string {
	return net.JoinHostPort(g.raw.Host, g.raw.Port)
}

func (g *grpcConfig) Host() string { return g.raw.Host }
func (g *grpcConfig) Port() string { return g.raw.Port }

func (g *grpcConfig) ShutdownTimeout() time.Duration {
	return g.raw.ShutdownTimeout
}
