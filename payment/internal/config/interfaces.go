package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PaymentGRPCConfig interface {
	GRPCAddress() string
	GRPCPort() string
	HTTPAddress() string
	HTTPPort() string
	ShutdownTimeout() time.Duration
}

type AuthGRPCService interface {
	AuthServiceAddress() string
	AuthServicePort() string
}
