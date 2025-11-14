package converter

import (
	"github.com/ZanDattSu/star-factory/order/internal/model"
	orderV1 "github.com/ZanDattSu/star-factory/shared/pkg/openapi/order/v1"
)

// OrderToAPI конвертирует service-модель → OpenAPI DTO.
func OrderToAPI(o *model.Order) *orderV1.OrderDto {
	if o == nil {
		return nil
	}

	dto := &orderV1.OrderDto{
		OrderUUID:  o.OrderUUID,
		UserUUID:   o.UserUUID,
		PartUuids:  o.PartUuids,
		TotalPrice: o.TotalPrice,
		Status:     OrderStatusToAPI(o.Status),
	}

	// TransactionUUID
	if o.TransactionUUID != nil {
		dto.TransactionUUID = orderV1.NewOptString(*o.TransactionUUID)
	}

	// PaymentMethod
	dto.PaymentMethod = orderV1.NewOptPaymentMethod(PaymentMethodToAPI(o.PaymentMethod))

	return dto
}

// OrderToModel конвертирует OpenAPI DTO → service-модель.
func OrderToModel(orderDto *orderV1.OrderDto) *model.Order {
	if orderDto == nil {
		return nil
	}

	o := &model.Order{
		OrderUUID:  orderDto.OrderUUID,
		UserUUID:   orderDto.UserUUID,
		PartUuids:  orderDto.PartUuids,
		TotalPrice: orderDto.TotalPrice,
		Status:     OrderStatusFromAPI(orderDto.Status),
	}

	// TransactionUUID
	if val, ok := orderDto.TransactionUUID.Get(); ok {
		o.TransactionUUID = &val
	}

	// PaymentMethod
	if val, ok := orderDto.PaymentMethod.Get(); ok {
		o.PaymentMethod = PaymentMethodToModel(val)
	}

	return o
}

// === OrderStatus converters ===

func OrderStatusToAPI(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusPendingPayment:
		return orderV1.OrderStatusPENDINGPAYMENT
	case model.OrderStatusPaid:
		return orderV1.OrderStatusPAID
	case model.OrderStatusCancelled:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusNOTSET
	}
}

func OrderStatusFromAPI(status orderV1.OrderStatus) model.OrderStatus {
	switch status {
	case orderV1.OrderStatusPENDINGPAYMENT:
		return model.OrderStatusPendingPayment
	case orderV1.OrderStatusPAID:
		return model.OrderStatusPaid
	case orderV1.OrderStatusCANCELLED:
		return model.OrderStatusCancelled
	default:
		return model.OrderStatusUnspecified
	}
}

// === PaymentMethod ===

// PaymentMethodToModel конвертирует orderV1.PaymentMethod → model.PaymentMethod
func PaymentMethodToModel(method orderV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCard
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSbp
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

// PaymentMethodToAPI конвертирует model.PaymentMethod → orderV1.PaymentMethod
func PaymentMethodToAPI(method model.PaymentMethod) orderV1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodSbp:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodInvestorMoney:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}
