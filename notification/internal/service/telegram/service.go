package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"
	"time"

	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/notification/internal/client/http"
	"github.com/ZanDattSu/star-factory/notification/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

const chatID = 725700609

//go:embed templates/order_paid_notification.tmpl
var orderPaidTemplateFS embed.FS

type orderPaid struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
	RegisteredAt    time.Time
}

var orderPaidTemplate = template.Must(template.ParseFS(orderPaidTemplateFS, "templates/order_paid_notification.tmpl"))

//go:embed templates/ship_assembled_notification.tmpl
var shipAssembledTemplateFS embed.FS

type shipAssembled struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
	RegisteredAt time.Time
}

var shipAssembledTemplate = template.Must(template.ParseFS(shipAssembledTemplateFS, "templates/ship_assembled_notification.tmpl"))

type service struct {
	telegramClient http.TelegramClient
}

func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) SendPaidNotification(ctx context.Context, paidEvent model.OrderPaidEvent) error {
	message, err := s.buildPaidMessage(paidEvent)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) SendAssembledNotification(ctx context.Context, shipAssembledEvent model.ShipAssembledEvent) error {
	message, err := s.buildAssembledMessage(shipAssembledEvent)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) buildPaidMessage(paidEvent model.OrderPaidEvent) (string, error) {
	data := orderPaid{
		EventUUID:       paidEvent.EventUUID,
		OrderUUID:       paidEvent.OrderUUID,
		UserUUID:        paidEvent.UserUUID,
		PaymentMethod:   string(paidEvent.PaymentMethod),
		TransactionUUID: paidEvent.TransactionUUID,
		RegisteredAt:    time.Now(),
	}

	var buf bytes.Buffer
	err := orderPaidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) buildAssembledMessage(shipAssembledEvent model.ShipAssembledEvent) (string, error) {
	data := shipAssembled{
		EventUUID:    shipAssembledEvent.EventUUID,
		OrderUUID:    shipAssembledEvent.OrderUUID,
		UserUUID:     shipAssembledEvent.UserUUID,
		BuildTimeSec: int64(shipAssembledEvent.BuildTime.Seconds()),
		RegisteredAt: time.Now(),
	}

	var buf bytes.Buffer
	err := shipAssembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
