package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-faster/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ordApi "order/internal/api/v1/order"
	inventoryService "order/internal/client/grpc/inventory/v1"
	paymentService "order/internal/client/grpc/payment/v1"
	ordRepo "order/internal/repository/order"
	httpServer "order/internal/server"
	ordService "order/internal/service/order"

	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

const (
	orderServicePort     = "8080"
	paymentServicePort   = "50052"
	inventoryServicePort = "50051"

	shutdownTimeout = 10 * time.Second
)

func newGRPCConnectWithoutSecure(port string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		getAddress(port),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // отключаем TLS
	)
	return conn, err
}

func getAddress(port string) string {
	return net.JoinHostPort("localhost", port)
}

func main() {
	logger := setupLogger()

	logger.Info("Creating payment gRPC client...")

	paymentConn, err := newGRPCConnectWithoutSecure(paymentServicePort)
	if err != nil {
		logger.Error("Failed to connect to payment gRPC service", slog.String("port", paymentServicePort), slogErr(err))
		// logger.Error("Failed to connect to payment gRPC service (%s): %v", inventoryServicePort, err)
		return
	}
	defer func() {
		if closeErr := paymentConn.Close(); closeErr != nil {
			logger.Warn("Failed to close payment gRPC connection", slogErr(closeErr))
		}
	}()

	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)
	logger.Info("Payment gRPC client created successfully", slog.String("port", paymentServicePort))

	logger.Info("======================================")

	logger.Info("Creating inventory gRPC client...")

	inventoryConn, err := newGRPCConnectWithoutSecure(inventoryServicePort)
	if err != nil {
		logger.Error("Failed to connect to inventory gRPC service", slog.String("port", inventoryServicePort), slogErr(err))
		return
	}
	defer func() {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			logger.Warn("Failed to close inventory gRPC connection", slogErr(closeErr))
		}
	}()

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	logger.Info("Inventory gRPC client created successfully", slog.String("port", inventoryServicePort))

	logger.Info("======================================")

	logger.Info("Creating order API handler...")
	orderStorage := ordRepo.NewRepository()
	orderService := ordService.NewService(
		orderStorage,
		paymentService.NewClient(paymentClient),
		inventoryService.NewClient(inventoryClient),
	)
	api := ordApi.NewApi(orderService)

	logger.Info("Creating HTTP server...")
	server, err := httpServer.NewHTTPServer(orderServicePort, api)
	if err != nil {
		logger.Error("Failed to create HTTP server", slogErr(err))
		return
	}
	logger.Info("HTTP server created successfully")

	go func() {
		logger.Info("HTTP server listening", slog.String("port", orderServicePort))
		if err := server.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", slogErr(err))
			return
		}
	}()

	gracefulShutdown()

	logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", slogErr(err))
		return
	}

	logger.Info("HTTP server stopped successfully")
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
