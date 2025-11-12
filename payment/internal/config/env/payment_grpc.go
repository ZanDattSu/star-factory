package env

import (
	"net"
	"time"

	"github.com/caarlos0/env/v11"
)

type paymentGrpcEnvConfig struct {
	Host                string        `env:"GRPC_HOST,required"`
	GRPCPort            string        `env:"GRPC_PORT,required"`
	HttpPort            string        `env:"HTTP_GATEWAY_PORT,required"`
	HttpShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT,required"`
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

func (p *paymentGrpcConfig) GRPCAddress() string {
	return net.JoinHostPort(p.raw.Host, p.raw.GRPCPort)
}

func (p *paymentGrpcConfig) GRPCPort() string {
	return p.raw.GRPCPort
}

func (p *paymentGrpcConfig) HttpAddress() string {
	return net.JoinHostPort(p.raw.Host, p.raw.HttpPort)
}

func (p *paymentGrpcConfig) HttpPort() string {
	return p.raw.HttpPort
}

func (p *paymentGrpcConfig) ShutdownTimeout() time.Duration {
	return p.raw.HttpShutdownTimeout
}
