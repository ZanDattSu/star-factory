package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"payment/internal/servers"

	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort = 8082
	grpcPort = 50052

	shutdownTimeout = 10 * time.Second
)

type PaymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func (ps PaymentService) PayOrder(_ context.Context, _ *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	u := uuid.New()

	log.Printf("Оплата прошла успешно, transaction_uuid:%s", u)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: u.String(),
	}, nil
}

func main() {
	api := &PaymentService{}

	gRPCServer, err := servers.NewGRPCServer(grpcPort, api)
	if err != nil {
		log.Printf("failed to create gRPC server: %v\n", err)
		return
	}

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", gRPCServer.GetPort())
		if err := gRPCServer.Serve(); err != nil {
			log.Printf("gRPC server failed: %v", err)
			return
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gatewayServer, err := servers.NewHTTPServer(ctx, httpPort, grpcPort)
	if err != nil {
		log.Printf("failed to create HTTP server: %v\n", err)
		return
	}

	// Запускаем HTTP сервер с gRPC Gateway
	go func() {
		log.Printf("🌐 HTTP server with gRPC-Gateway listening on %d\n", gatewayServer.GetPort())
		if err := gatewayServer.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to serve HTTP: %s\n", err)
			return
		}
	}()

	// Graceful shutdown
	gracefulShutdown()

	log.Println("🛑 Shutting down servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	// Сначала останавливаем HTTP сервер
	log.Println("🛑 Shutting down HTTP server...")
	if err := gatewayServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}
	log.Println("✅ HTTP server stopped")

	log.Println("🛑 Shutting down gRPC server...")
	gRPCServer.Shutdown()
	log.Println("✅ gRPC server stopped")
}

func gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
