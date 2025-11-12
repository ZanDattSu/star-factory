package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PaymentGRPCConfig interface {
	GRPCAddress() string
	GRPCPort() string
	HttpAddress() string
	HttpPort() string
	ShutdownTimeout() time.Duration
}
