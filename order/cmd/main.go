package main

import (
	"context"
	"log"
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
	log.Println("Creating payment gRPC client...")
	paymentConn, err := newGRPCConnectWithoutSecure(paymentServicePort)
	if err != nil {
		log.Printf("Failed to connect to payment gRPC service (%s): %v", inventoryServicePort, err)
		return
	}
	defer func() {
		if closeErr := paymentConn.Close(); closeErr != nil {
			log.Printf("Failed to close payment gRPC connection: %v", closeErr)
		}
	}()

	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)
	log.Printf("Payment gRPC client created successfully (%s)", paymentServicePort)

	log.Println("======================================")

	log.Println("Creating inventory gRPC client...")
	inventoryConn, err := newGRPCConnectWithoutSecure(inventoryServicePort)
	if err != nil {
		log.Printf("Failed to connect to inventory gRPC service (%s): %v", inventoryServicePort, err)
		return
	}
	defer func() {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			log.Printf("Failed to close inventory gRPC connection: %v", closeErr)
		}
	}()

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	log.Printf("Inventory gRPC client created successfully (%s)", inventoryServicePort)

	log.Println("======================================")

	log.Println("Creating order API handler...")
	orderStorage := ordRepo.NewRepository()
	orderService := ordService.NewService(orderStorage, paymentClient, inventoryClient)
	api := ordApi.NewApi(orderService)

	log.Println("Creating HTTP server...")
	server, err := httpServer.NewHTTPServer(orderServicePort, api)
	if err != nil {
		log.Printf("Failed to create HTTP server: %v", err)
		return
	}
	log.Println("HTTP server created successfully")

	go func() {
		log.Printf("HTTP server listening on :%s\n", orderServicePort)
		if err := server.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
			return
		}
	}()

	gracefulShutdown()

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
		return
	}

	log.Println("HTTP server stopped successfully")
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
