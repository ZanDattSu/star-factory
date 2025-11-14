package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"order/internal/config"
	"order/internal/server"
	"platform/pkg/closer"
	"platform/pkg/logger"
	"platform/pkg/migrator"
)

type App struct {
	diContainer *diContainer
	server      *server.HTTPServer
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
	return a.runServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initLogger,
		a.initCloser,
		a.initDI,
		a.migratorUp,
		a.initServer,
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

func (a *App) initServer(ctx context.Context) error {
	httpServer, err := server.NewHTTPServer(
		config.AppConfig().OrderHTTP.OrderAddress(),
		a.diContainer.OrderApi(ctx))
	if err != nil {
		logger.Error(ctx, "Failed to create HTTP server", zap.Error(err))
		return err
	}

	a.server = httpServer

	closer.AddNamed("Http Server", func(ctx context.Context) error {
		return a.server.Shutdown(ctx)
	})

	return nil
}

func (a *App) runServer(ctx context.Context) error {
	logger.Info(ctx, "HTTP server listening on: "+config.AppConfig().OrderHTTP.OrderPort())
	err := a.server.Serve()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
