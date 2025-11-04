package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"shared/pkg/interceptor"
	"syscall"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentv1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort = 8082
	grpcPort = 50052

	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type PaymentService struct {
	paymentv1.UnimplementedPaymentServiceServer
}

func (ps PaymentService) PayOrder(_ context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	u := uuid.New()

	log.Printf("Оплата прошла успешно, transaction_uuid:%s", u)

	return &paymentv1.PayOrderResponse{
		TransactionUuid: u.String(),
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerInterceptor(),
			interceptor.ValidationInterceptor(),
		),
	)

	service := PaymentService{}

	paymentv1.RegisterPaymentServiceServer(server, service)

	reflection.Register(server)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err := server.Serve(listener)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Запускаем HTTP сервер с gRPC Gateway
	var gatewayServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mux := runtime.NewServeMux()

		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		err := paymentv1.RegisterPaymentServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
			return
		}

		gatewayServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           mux,
			ReadHeaderTimeout: readHeaderTimeout,
		}

		log.Printf("🌐 HTTP server with gRPC-Gateway listening on %d\n", httpPort)
		err = gatewayServer.ListenAndServe()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Сначала останавливаем HTTP сервер
	if gatewayServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := gatewayServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}

		log.Println("✅ HTTP server stopped")
	}

	log.Println("🛑 Shutting down gRPC server...")
	if closeErr := listener.Close(); closeErr != nil {
		log.Printf("failed to close listener: %v\n", closeErr)
	}

	server.GracefulStop()
	log.Println("✅ Server stopped")
}
