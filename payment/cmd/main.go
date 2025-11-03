package main

import (
	"context"
	"fmt"
	paymentv1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const grpcPort = 50051

type InventoryService struct {
	paymentv1.UnimplementedPaymentServiceServer
}

func (is InventoryService) PayOrder(_ context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
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

	server := grpc.NewServer()

	service := InventoryService{}

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

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down gRPC server...")
	if closeErr := listener.Close(); closeErr != nil {
		log.Printf("failed to close listener: %v\n", closeErr)
	}

	server.GracefulStop()
	log.Println("✅ Server stopped")
}
