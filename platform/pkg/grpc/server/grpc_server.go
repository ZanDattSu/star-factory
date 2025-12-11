package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/ZanDattSu/star-factory/platform/pkg/closer"
	"github.com/ZanDattSu/star-factory/platform/pkg/grpc/health"
	"github.com/ZanDattSu/star-factory/platform/pkg/grpc/interceptor"
)

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
}

type RegisterServerFunc func(server *grpc.Server)

type Options struct {
	Register RegisterServerFunc
	Auth     *interceptor.AuthInterceptor
}

func NewGRPCServer(grpcAddress string, opts Options) (*GRPCServer, error) {
	listener, err := newListener(grpcAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	interceptors := []grpc.UnaryServerInterceptor{
		interceptor.LoggerInterceptor(),
		interceptor.ValidationInterceptor(),
	}

	if opts.Auth != nil {
		interceptors = append(interceptors, opts.Auth.Unary())
	}

	server := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(interceptors...),
	)

	reflection.Register(server)

	// Регистрируем health service для проверки работоспособности
	health.RegisterService(server)

	if opts.Register != nil {
		opts.Register(server)
	}

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
