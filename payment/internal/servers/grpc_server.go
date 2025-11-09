package servers

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"shared/pkg/interceptor"

	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
	port     int
}

func (s *GRPCServer) GetPort() int {
	return s.port
}

func NewGRPCServer(grpcPort int, api paymentV1.PaymentServiceServer) (*GRPCServer, error) {
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

	paymentV1.RegisterPaymentServiceServer(server, api)
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
