package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentv1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50051

type PaymentService struct {
	paymentv1.UnimplementedPaymentServiceServer
}

func (ps PaymentService) PayOrder(_ context.Context, _ *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
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
