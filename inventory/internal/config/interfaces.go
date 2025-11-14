package config

import "time"

type App interface {
	ShutdownTimeout() time.Duration
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryGRPCConfig interface {
	GRPCAddress() string
	GRPCPort() string
}

type InventoryHTTPConfig interface {
	HTTPAddress() string
	HTTPPort() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
	ConnectTimeout() time.Duration
	ShutdownTimeout() time.Duration
}
