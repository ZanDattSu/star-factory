package telegram

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/notification/internal/client/grpc/auth"
	"github.com/ZanDattSu/star-factory/notification/internal/client/http"
	"github.com/ZanDattSu/star-factory/notification/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

const defaultChatID int64 = 725700609

type service struct {
	telegramClient http.TelegramClient
	authClient     auth.AuthClient
}

func NewService(telegramClient http.TelegramClient, authClient auth.AuthClient) *service {
	return &service{telegramClient: telegramClient, authClient: authClient}
}

func (s *service) SendPaidNotification(ctx context.Context, paidEvent model.OrderPaidEvent) error {
	message, err := s.buildPaidMessage(paidEvent)
	if err != nil {
		return err
	}

	isSub, chatID, err := s.telegramSubscription(ctx, paidEvent.UserUUID)
	if err != nil {
		return err
	}

	if !isSub {
		logger.Info(
			ctx,
			"user is not subscribed to telegram notifications",
			zap.String("user_uuid", paidEvent.UserUUID),
		)
		return nil
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(
		ctx,
		"telegram message sent",
		zap.Int64("chat_id", chatID),
		zap.String("user_uuid", paidEvent.UserUUID),
	)
	return nil
}

func (s *service) SendAssembledNotification(ctx context.Context, shipAssembledEvent model.ShipAssembledEvent) error {
	message, err := s.buildAssembledMessage(shipAssembledEvent)
	if err != nil {
		return err
	}

	isSub, chatID, err := s.telegramSubscription(ctx, shipAssembledEvent.UserUUID)
	if err != nil {
		return err
	}

	if !isSub {
		logger.Info(
			ctx,
			"user is not subscribed to telegram notifications",
			zap.String("user_uuid", shipAssembledEvent.UserUUID),
		)
		return nil
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) telegramSubscription(ctx context.Context, userUUID string) (bool, int64, error) {
	user, err := s.authClient.GetUser(ctx, userUUID)
	if err != nil {
		return false, 0, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return false, 0, fmt.Errorf("user %s not found", userUUID)
	}

	for _, nm := range user.Info.NotificationMethods {
		if nm.ProviderName != "telegram" {
			continue
		}

		chatID := defaultChatID

		if nm.Target != "" {
			parsedChatID, err := parseChatID(nm.Target)
			if err != nil {
				return false, 0, err
			}
			chatID = parsedChatID
		}

		return true, chatID, nil
	}

	return false, 0, nil
}

func parseChatID(target string) (int64, error) {
	id, err := strconv.ParseInt(target, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid telegram chat_id %q: %w", target, err)
	}
	return id, nil
}
