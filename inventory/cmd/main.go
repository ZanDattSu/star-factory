package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	partV1Api "inventory/internal/api/v1/part"
	"inventory/internal/config"
	partRepo "inventory/internal/repository/part/mongodb"
	"inventory/internal/servers"
	partService "inventory/internal/service/part"
)

const (
	configPath = "./deploy/compose/inventory/.env"
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	logger := setupLogger()

	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig().Mongo.ConnectTimeout())
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		logger.Error("Failed to create MongoDB connect", slogErr(err))
		return
	}

	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			logger.Error("Error disconnecting from MongoDB", slogErr(cerr))
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Error("Failed to ping MongoDB", slogErr(err))
		return
	}

	logger.Info("Connected to MongoDB")

	partsCollection := client.Database("inventory")

	repo := partRepo.NewRepository(partsCollection)
	repo.InitTestData()

	service := partService.NewService(repo)
	api := partV1Api.NewApi(service)

	gRPCServer, err := servers.NewGRPCServer(config.AppConfig().InventoryGRPC.Address(), api)
	if err != nil {
		logger.Error("Failed to create gRPC server", slogErr(err))
		return
	}

	// Запускаем gRPC сервер
	go func() {
		logger.Info("GRPC server listening on", slog.String("port", gRPCServer.GetPort()))
		if err := gRPCServer.Serve(); err != nil {
			logger.Error("GRPC server failed", slogErr(err))
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
		logger.Error("Failed to create HTTP server", slogErr(err))
		return
	}

	// Запускаем HTTP сервер с gRPC Gateway
	go func() {
		logger.Info("HTTP server with gRPC-Gateway listening on", slog.String("port", gatewayServer.GetPort()))
		if err = gatewayServer.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Failed to serve HTTP", slogErr(err))
			return
		}
	}()

	gracefulShutdown()

	logger.Info("Shutting down servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, config.AppConfig().Mongo.ShutdownTimeout())
	defer shutdownCancel()

	// Сначала останавливаем HTTP сервер
	logger.Info("Shutting down HTTP server...")
	if err = gatewayServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP shutdown error", slogErr(err))
	}
	logger.Info("HTTP server stopped")

	logger.Info("Shutting down gRPC server...")
	gRPCServer.Shutdown()
	logger.Info("GRPC server stopped")
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
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
