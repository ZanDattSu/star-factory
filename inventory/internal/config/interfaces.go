package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryGRPCConfig interface {
	Address() string
	GRPCPort() string
	HttpPort() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
	ConnectTimeout() time.Duration
	ShutdownTimeout() time.Duration
}
