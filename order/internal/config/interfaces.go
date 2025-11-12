package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OrderHttpConfig interface {
	OrderAddress() string
	OrderPort() string
	ReadHeaderTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

type PaymentGRPCService interface {
	PaymentAddress() string
	PaymentServicePort() string
}

type InventoryGrpcService interface {
	InventoryAddress() string
	InventoryServicePort() string
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
	MigrationsPath() string
}
