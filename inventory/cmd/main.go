package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	partV1Api "inventory/internal/api/v1/part"
	partRepo "inventory/internal/repository/part/mongodb"
	"inventory/internal/servers"
	partService "inventory/internal/service/part"
)

const (
	grpcPort = 50051
	httpPort = 8081

	shutdownTimeout = 10 * time.Second
)

func main() {
	ctx := context.Background()
	logger := setupLogger()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Failed to load .env file", slogErr(err))
		return
	}

	dbURI := os.Getenv("MONGO_URI")
	if dbURI == "" {
		logger.Error("Error: environment variable MONGO_URI is not set")
		return
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
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

	gRPCServer, err := servers.NewGRPCServer(grpcPort, api)
	if err != nil {
		logger.Error("Failed to create gRPC server", slogErr(err))
		return
	}

	// Запускаем gRPC сервер
	go func() {
		logger.Info("GRPC server listening on", slog.Int("port", gRPCServer.GetPort()))
		if err := gRPCServer.Serve(); err != nil {
			logger.Error("GRPC server failed", slogErr(err))
			return
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gatewayServer, err := servers.NewHTTPServer(ctx, httpPort, grpcPort)
	if err != nil {
		logger.Error("Failed to create HTTP server", slogErr(err))
		return
	}

	// Запускаем HTTP сервер с gRPC Gateway
	go func() {
		logger.Info("HTTP server with gRPC-Gateway listening on", slog.Int("port", gatewayServer.GetPort()))
		if err := gatewayServer.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Failed to serve HTTP", slogErr(err))
			return
		}
	}()

	gracefulShutdown()

	logger.Info("Shutting down servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, shutdownTimeout)
	defer shutdownCancel()

	// Сначала останавливаем HTTP сервер
	logger.Info("Shutting down HTTP server...")
	if err := gatewayServer.Shutdown(shutdownCtx); err != nil {
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
