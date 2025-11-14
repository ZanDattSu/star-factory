package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	partV1Api "inventory/internal/api/v1/part"
	"inventory/internal/config"
	partRepo "inventory/internal/repository/part/mongodb"
	"inventory/internal/servers"
	partService "inventory/internal/service/part"

	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

const (
	configPath = "./deploy/compose/inventory/.env"
)

func main() {
	ctx := context.Background()

	if err := config.Load(configPath); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig().Mongo.ConnectTimeout())
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		logger.Error(ctx, "Failed to create MongoDB connect", zap.Error(err))
		return
	}

	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			logger.Error(ctx, "Error disconnecting from MongoDB", zap.Error(cerr))
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Error(ctx, "Failed to ping MongoDB", zap.Error(err))
		return
	}

	logger.Info(ctx, "Connected to MongoDB")

	partsCollection := client.Database("inventory")

	repo := partRepo.NewRepository(partsCollection)
	repo.InitTestData()

	service := partService.NewService(repo)
	api := partV1Api.NewApi(service)

	gRPCServer, err := servers.NewGRPCServer(config.AppConfig().InventoryGRPC.Address(), api)
	if err != nil {
		logger.Error(ctx, "Failed to create gRPC server", zap.Error(err))
		return
	}

	// Запускаем gRPC сервер
	go func() {
		logger.Info(ctx, "GRPC server listening on", zap.String("port", gRPCServer.GetPort()))
		if err := gRPCServer.Serve(); err != nil {
			logger.Error(ctx, "GRPC server failed", zap.Error(err))
			return
		}
	}()

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	gatewayServer, err := servers.NewHTTPServer(ctx,
		config.AppConfig().InventoryGRPC.HttpPort(),
		config.AppConfig().InventoryGRPC.Address(),
	)
	if err != nil {
		logger.Error(ctx, "Failed to create HTTP server", zap.Error(err))
		return
	}

	// Запускаем HTTP сервер с gRPC Gateway
	go func() {
		logger.Info(ctx, "HTTP server with gRPC-Gateway listening on", zap.String("port", gatewayServer.GetPort()))
		if err = gatewayServer.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(ctx, "Failed to serve HTTP", zap.Error(err))
			return
		}
	}()

	gracefulShutdown()

	logger.Info(ctx, "Shutting down servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, config.AppConfig().Mongo.ShutdownTimeout())
	defer shutdownCancel()

	// Сначала останавливаем HTTP сервер
	logger.Info(ctx, "Shutting down HTTP server...")
	if err = gatewayServer.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, "HTTP shutdown error", zap.Error(err))
	}
	logger.Info(ctx, "HTTP server stopped")

	logger.Info(ctx, "Shutting down gRPC server...")
	gRPCServer.Shutdown()
	logger.Info(ctx, "GRPC server stopped")
}

// gracefulShutdown мягко завершает работу программы
// когда в канал quit поступает один из сисколлов ОС
//
// SIGTERM - "вежливая" просьба завершиться,
// SIGINT - прерывание с клавиатуры (Ctrl+C)
func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
