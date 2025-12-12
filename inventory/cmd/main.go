package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/inventory/internal/app"
	"github.com/ZanDattSu/star-factory/inventory/internal/config"
	"github.com/ZanDattSu/star-factory/platform/pkg/closer"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	"github.com/ZanDattSu/star-factory/platform/pkg/path"
)

var configPath = path.GetPathRelativeToRoot("deploy/compose/inventory/.env")

func main() {
	if err := config.Load(configPath); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	// SIGTERM - "вежливая" просьба завершиться

	// SIGINT - прерывание с клавиатуры (Ctrl+C)

	osSignals := []os.Signal{syscall.SIGINT, syscall.SIGTERM}

	appCtx, appCancel := signal.NotifyContext(context.Background(), osSignals...)

	defer appCancel()

	defer gracefulShutdown()

	closer.Configure(osSignals...)

	a, err := app.New(appCtx)
	if err != nil {

		logger.Error(appCtx, "Не удалось создать приложение", zap.Error(err))

		return

	}

	go func() {
		if err = a.RunGRPC(appCtx); err != nil {

			logger.Error(appCtx, "Ошибка при работе приложения", zap.Error(err))

			return

		}
	}()

	if err = a.RunHTTP(appCtx); err != nil {

		logger.Error(appCtx, "Ошибка при работе приложения", zap.Error(err))

		return

	}
}

// gracefulShutdown мягко завершает работу программы

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig().App.ShutdownTimeout())

	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "Ошибка при завершении работы", zap.Error(err))
	}
}
