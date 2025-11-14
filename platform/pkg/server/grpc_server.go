package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"platform/pkg/closer"
	"platform/pkg/grpc/health"
	"shared/pkg/interceptor"
)

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
}

type RegisterServerFunc func(server *grpc.Server)

func NewGRPCServer(grpcAddress string, registerServerFunc RegisterServerFunc) (*GRPCServer, error) {
	listener, err := newListener(grpcAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	server := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerInterceptor(),
			interceptor.ValidationInterceptor(),
		),
	)

	reflection.Register(server)

	// Регистрируем health service для проверки работоспособности
	health.RegisterService(server)

	registerServerFunc(server)

	return &GRPCServer{
		server:   server,
		listener: listener,
	}, nil
}

func newListener(grpcAddress string) (net.Listener, error) {
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return nil, err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}
		return nil
	})

	return listener, nil
}

func (s *GRPCServer) Serve() error {
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Shutdown() {
	// GracefulStop автоматически закроет listener
	s.server.GracefulStop()
}
