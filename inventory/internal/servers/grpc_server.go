package servers

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"shared/pkg/interceptor"

	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
	port     int
}

func NewGRPCServer(grpcPort int, api inventoryV1.InventoryServiceServer) (*GRPCServer, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerInterceptor(),
			interceptor.ValidationInterceptor(),
		),
	)

	inventoryV1.RegisterInventoryServiceServer(server, api)
	reflection.Register(server)

	return &GRPCServer{
		server:   server,
		listener: listener,
		port:     grpcPort,
	}, nil
}

func (s *GRPCServer) Serve() error {
	log.Printf("🚀 gRPC server listening on %d\n", s.port)
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Shutdown() error {
	log.Println("🛑 Shutting down gRPC server...")
	if closeErr := s.listener.Close(); closeErr != nil {
		return fmt.Errorf("failed to close listener: %w", closeErr)
	}

	s.server.GracefulStop()
	log.Println("✅ Server stopped")
	return nil
}
