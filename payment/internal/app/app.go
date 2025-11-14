package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc"

	"github.com/ZanDattSu/star-factory/payment/internal/config"
	"github.com/ZanDattSu/star-factory/payment/internal/redirect"
	"github.com/ZanDattSu/star-factory/platform/pkg/closer"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	platformServer "github.com/ZanDattSu/star-factory/platform/pkg/server"
	paymentV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/payment/v1"
)

type App struct {
	diContainer *diContainer
	gRPCServer  *platformServer.GRPCServer
	httpServer  *redirect.HTTPServer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) RunGRPC(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) RunHTTP(ctx context.Context) error {
	return a.runHTTPServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initLogger,
		a.initCloser,
		a.initDI,
		a.initGRPCServer,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	srv, err := platformServer.NewGRPCServer(
		config.AppConfig().PaymentGRPC.GRPCAddress(),
		func(s *grpc.Server) {
			paymentV1.RegisterPaymentServiceServer(s, a.diContainer.PaymentV1Api(ctx))
		},
	)
	if err != nil {
		return err
	}

	a.gRPCServer = srv

	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.gRPCServer.Shutdown()
		return nil
	})

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("gRPC PaymentService server listening on %s", config.AppConfig().PaymentGRPC.GRPCAddress()))

	err := a.gRPCServer.Serve()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	srv, err := redirect.NewHTTPServer(ctx,
		config.AppConfig().PaymentGRPC.GRPCAddress(),
		config.AppConfig().PaymentGRPC.HTTPAddress())
	if err != nil {
		return err
	}

	a.httpServer = srv

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("HTTP server with gRPC-Gateway listening on %s", config.AppConfig().PaymentGRPC.HTTPAddress()))
	if err := a.httpServer.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
