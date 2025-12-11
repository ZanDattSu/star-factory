package config

import "time"

type App interface {
	ShutdownTimeout() time.Duration
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type GRPCConfig interface {
	Address() string
	Host() string
	Port() string
	ShutdownTimeout() time.Duration
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
	MigrationsPath() string
}

type RedisConfig interface {
	Address() string
	Host() string
	Port() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type SessionConfig interface {
	TTL() time.Duration
}
