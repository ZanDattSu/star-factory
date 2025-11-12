package servers

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"shared/pkg/interceptor"

	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
	port     string
}

func (s *GRPCServer) GetPort() string {
	return s.port
}

func NewGRPCServer(grpcPort string, api inventoryV1.InventoryServiceServer) (*GRPCServer, error) {
	listener, err := net.Listen("tcp", grpcPort)
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
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Shutdown() {
	// GracefulStop автоматически закроет listener
	s.server.GracefulStop()
}
