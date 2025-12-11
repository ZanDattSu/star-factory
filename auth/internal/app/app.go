package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/ZanDattSu/star-factory/auth/internal/config"
	"github.com/ZanDattSu/star-factory/platform/pkg/closer"
	platformServer "github.com/ZanDattSu/star-factory/platform/pkg/grpc/server"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	"github.com/ZanDattSu/star-factory/platform/pkg/migrator"
	authV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/auth/v1"
	userV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/user/v1"
)

type App struct {
	diContainer *diContainer
	gRPCServer  *platformServer.GRPCServer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initLogger,
		a.initCloser,
		a.initDI,
		a.migratorUp,
		a.initGRPCServer,
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

func (a *App) migratorUp(ctx context.Context) error {
	pool := a.diContainer.PostgreSQLPool(ctx)

	migratorDir := config.AppConfig().Postgres.MigrationsPath()
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*pool.Config().ConnConfig), migratorDir)

	err := migratorRunner.Up()
	if err != nil {
		logger.Error(ctx, "Database migration error", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	srv, err := platformServer.NewGRPCServer(
		config.AppConfig().GRPC.Address(),
		platformServer.Options{
			Register: func(s *grpc.Server) {
				authV1.RegisterAuthServiceServer(s, a.diContainer.AuthApi(ctx))
				userV1.RegisterUserServiceServer(s, a.diContainer.UserApi(ctx))
			},
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
	logger.Info(ctx, fmt.Sprintf("gRPC PaymentService server listening on %s", config.AppConfig().GRPC.Address()))

	err := a.gRPCServer.Serve()
	if err != nil {
		return err
	}

	return nil
}
