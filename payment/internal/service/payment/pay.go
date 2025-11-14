package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/ZanDattSu/star-factory/payment/internal/model"
)

func (s *service) PayOrder(_ context.Context, _, _ string, _ model.PaymentMethod) string {
	return uuid.New().String()
}
