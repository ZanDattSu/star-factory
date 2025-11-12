package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ordApi "order/internal/api/v1/order"
	inventoryService "order/internal/client/grpc/inventory/v1"
	paymentService "order/internal/client/grpc/payment/v1"
	"order/internal/config"
	"order/internal/migrator"
	ordRepo "order/internal/repository/order/postgresql"
	httpServer "order/internal/server"
	ordService "order/internal/service/order"

	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

const (
	configPath = "./deploy/compose/order/.env"
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic("failed to load config file: " + err.Error())
	}

	ctx := context.Background()
	logger := setupLogger()

	dbURI := config.AppConfig().Postgres.URI()

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		logger.Error("Failed to create pgxpool connect", slogErr(err))
		return
	}

	defer func() {
		pool.Close()
	}()

	err = pool.Ping(ctx)
	if err != nil {
		logger.Error("Database is unavailable", slogErr(err))
		return
	}

	migratorDir := config.AppConfig().Postgres.MigrationsPath()
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*pool.Config().ConnConfig), migratorDir)

	err = migratorRunner.Up()
	if err != nil {
		logger.Error("Database migration error", slogErr(err))
		return
	}

	logger.Info("Creating payment gRPC client...")

	paymentConn, err := newGRPCConnectWithoutSecure(config.AppConfig().Payment.PaymentAddress())
	if err != nil {
		logger.Error(
			"Failed to connect to payment gRPC service",
			slog.String("port", config.AppConfig().Payment.PaymentServicePort()),
			slogErr(err))
		return
	}
	defer func() {
		if closeErr := paymentConn.Close(); closeErr != nil {
			logger.Warn("Failed to close payment gRPC connection", slogErr(closeErr))
		}
	}()
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)
	logger.Info("Payment gRPC client created successfully", slog.String("port", config.AppConfig().Payment.PaymentServicePort()))

	logger.Info("======================================")

	logger.Info("Creating inventory gRPC client...")
	inventoryConn, err := newGRPCConnectWithoutSecure(config.AppConfig().Inventory.InventoryAddress())
	if err != nil {
		logger.Error(
			"Failed to connect to inventory gRPC service",
			slog.String("port", config.AppConfig().Inventory.InventoryServicePort()),
			slogErr(err))
		return
	}
	defer func() {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			logger.Warn("Failed to close inventory gRPC connection", slogErr(closeErr))
		}
	}()
	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	logger.Info(
		"Inventory gRPC client created successfully",
		slog.String("port", config.AppConfig().Inventory.InventoryServicePort()))

	logger.Info("======================================")

	logger.Info("Creating order API handler...")
	orderRepository := ordRepo.NewRepository(pool)
	orderService := ordService.NewService(
		orderRepository,
		paymentService.NewClient(paymentClient),
		inventoryService.NewClient(inventoryClient),
	)
	orderApi := ordApi.NewApi(orderService)

	logger.Info("Creating HTTP server...")
	server, err := httpServer.NewHTTPServer(config.AppConfig().OrderHttp.OrderAddress(), orderApi)
	if err != nil {
		logger.Error("Failed to create HTTP server", slogErr(err))
		return
	}
	logger.Info("HTTP server created successfully")

	go func() {
		logger.Info("HTTP server listening", slog.String("port", server.GetPort()))
		if err := server.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", slogErr(err))
			return
		}
	}()

	gracefulShutdown()

	logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, config.AppConfig().OrderHttp.ShutdownTimeout())
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", slogErr(err))
		return
	}

	logger.Info("HTTP server stopped successfully")
}

func newGRPCConnectWithoutSecure(port string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		port,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // отключаем TLS
	)
	return conn, err
}

func slogErr(err error) slog.Attr {
	return slog.Attr{
		Key:   "Error",
		Value: slog.StringValue(err.Error()),
	}
}

func setupLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)
}

// gracefulShutdown мягко завершает работу программы
// когда в канал quit поступает один из сисколлов ОС
//
// SIGTERM - "вежливая" просьба завершиться,
// SIGINT - прерывание с клавиатуры (Ctrl+C)
func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
