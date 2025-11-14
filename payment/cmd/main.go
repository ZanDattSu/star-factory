package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"payment/internal/api/v1/payment"
	"payment/internal/config"
	"payment/internal/servers"
	payService "payment/internal/service/payment"

	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

const (
	configPath = "./deploy/compose/payment/.env"
)

func main() {
	ctx := context.Background()

	if err := config.Load(configPath); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	service := payService.NewService()
	api := payment.NewApi(service)

	gRPCServer, err := servers.NewGRPCServer(config.AppConfig().PaymentGRPC.GRPCAddress(), api)
	if err != nil {
		logger.Error(ctx, "Failed to create gRPC server", zap.Error(err))
		return
	}

	go func() {
		logger.Info(ctx, "GRPC server listening on", zap.String("port", gRPCServer.GetPort()))
		if err := gRPCServer.Serve(); err != nil {
			logger.Error(ctx, "GRPC server failed", zap.Error(err))
			return
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gatewayServer, err := servers.NewHTTPServer(
		ctx,
		config.AppConfig().PaymentGRPC.HttpAddress(),
		config.AppConfig().PaymentGRPC.GRPCAddress())
	if err != nil {
		logger.Error(ctx, "Failed to create HTTP server", zap.Error(err))
		return
	}

	// Запускаем HTTP сервер с gRPC Gateway
	go func() {
		logger.Info(ctx, "HTTP server with gRPC-Gateway listening on", zap.String("port", gatewayServer.GetPort()))
		if err := gatewayServer.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(ctx, "Failed to serve HTTP", zap.Error(err))
			return
		}
	}()

	// Graceful shutdown
	gracefulShutdown()

	logger.Info(ctx, "Shutting down servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), config.AppConfig().PaymentGRPC.ShutdownTimeout())
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
