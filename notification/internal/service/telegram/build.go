package telegram

import (
	"bytes"
	"embed"
	"text/template"
	"time"

	"github.com/ZanDattSu/star-factory/notification/internal/model"
)

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
