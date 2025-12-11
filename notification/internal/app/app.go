package app

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/notification/internal/config"
	"github.com/ZanDattSu/star-factory/platform/pkg/closer"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
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
	errCh := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runOrderPaidConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("consumer crashed: %w", err)
		}
	}()
	go func() {
		if err := a.runShipAssembledConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("consumer crashed: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		cancel()
		<-ctx.Done()
		return err
	}
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initTelegramBot,
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

func (a *App) runOrderPaidConsumer(ctx context.Context) error {
	logger.Info(ctx, "OrderPaid Kafka consumer starting")

	err := a.diContainer.OrderPaidConsumerService().RunOrderPaidConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runShipAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "ShipAssembled Kafka consumer starting")

	err := a.diContainer.ShipAssembledConsumerService().RunShipAssembledConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initTelegramBot(ctx context.Context) error {
	telegramBot := a.diContainer.TelegramBot()

	telegramBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		logger.Info(ctx, "chat id", zap.Int64("chat_id", update.Message.Chat.ID))

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Notification Bot активирован! Теперь вы будете получать уведомления о новых заказах.",
		})
		if err != nil {
			logger.Error(ctx, "Failed to send activation message", zap.Error(err))
		}
	})

	go func() {
		logger.Info(ctx, "Telegram bot started...")
		telegramBot.Start(ctx)
	}()

	return nil
}
