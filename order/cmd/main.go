package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
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

	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

const (
	configPath = "./deploy/compose/order/.env"
)

func main() {
	ctx := context.Background()

	if err := config.Load(configPath); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	dbURI := config.AppConfig().Postgres.URI()

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		logger.Error(ctx, "Failed to create pgxpool connect", zap.Error(err))
		return
	}

	defer func() {
		pool.Close()
	}()

	err = pool.Ping(ctx)
	if err != nil {
		logger.Error(ctx, "Database is unavailable", zap.Error(err))
		return
	}

	migratorDir := config.AppConfig().Postgres.MigrationsPath()
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*pool.Config().ConnConfig), migratorDir)

	err = migratorRunner.Up()
	if err != nil {
		logger.Error(ctx, "Database migration error", zap.Error(err))
		return
	}

	logger.Info(ctx, "Creating payment gRPC client...")

	paymentConn, err := newGRPCConnectWithoutSecure(config.AppConfig().Payment.PaymentAddress())
	if err != nil {
		logger.Error(ctx,
			"Failed to connect to payment gRPC service",
			zap.String("port", config.AppConfig().Payment.PaymentServicePort()),
			zap.Error(err))
		return
	}
	defer func() {
		if closeErr := paymentConn.Close(); closeErr != nil {
			logger.Warn(ctx, "Failed to close payment gRPC connection", zap.Error(closeErr))
		}
	}()
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)
	logger.Info(ctx,
		"Payment gRPC client created successfully",
		zap.String("port", config.AppConfig().Payment.PaymentServicePort()))

	logger.Info(ctx, "Creating inventory gRPC client...")
	inventoryConn, err := newGRPCConnectWithoutSecure(config.AppConfig().Inventory.InventoryAddress())
	if err != nil {
		logger.Error(ctx,
			"Failed to connect to inventory gRPC service",
			zap.String("port", config.AppConfig().Inventory.InventoryServicePort()),
			zap.Error(err))
		return
	}
	defer func() {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			logger.Warn(ctx, "Failed to close inventory gRPC connection", zap.Error(closeErr))
		}
	}()
	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	logger.Info(ctx,
		"Inventory gRPC client created successfully",
		zap.String("port", config.AppConfig().Inventory.InventoryServicePort()))

	logger.Info(ctx, "Creating order API handler...")
	orderRepository := ordRepo.NewRepository(pool)
	orderService := ordService.NewService(
		orderRepository,
		paymentService.NewClient(paymentClient),
		inventoryService.NewClient(inventoryClient),
	)
	orderApi := ordApi.NewApi(orderService)

	logger.Info(ctx, "Creating HTTP server...")
	server, err := httpServer.NewHTTPServer(config.AppConfig().OrderHttp.OrderAddress(), orderApi)
	if err != nil {
		logger.Error(ctx, "Failed to create HTTP server", zap.Error(err))
		return
	}
	logger.Info(ctx, "HTTP server created successfully")

	go func() {
		logger.Info(ctx, "HTTP server listening", zap.String("port", server.GetPort()))
		if err := server.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(ctx, "HTTP server error", zap.Error(err))
			return
		}
	}()

	gracefulShutdown()

	logger.Info(ctx, "Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, config.AppConfig().OrderHttp.ShutdownTimeout())
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, "server shutdown error", zap.Error(err))
		return
	}

	logger.Info(ctx, "HTTP server stopped successfully")
}

func newGRPCConnectWithoutSecure(port string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		port,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // отключаем TLS
	)
	return conn, err
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
