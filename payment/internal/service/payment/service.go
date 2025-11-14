package payment

import srvc "github.com/ZanDattSu/star-factory/payment/internal/service"

// Компиляторная проверка: убеждаемся, что *service реализует интерфейс PaymentService.
var _ srvc.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
