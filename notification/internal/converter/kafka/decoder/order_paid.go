package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/ZanDattSu/star-factory/notification/internal/model"
	eventsV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/events/v1"
)

type orderPaidDecoder struct{}

func NewOrderPaidDecoder() *orderPaidDecoder {
	return &orderPaidDecoder{}
}

func (d *orderPaidDecoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderPaidEvent{
		EventUUID:       pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		PaymentMethod:   mapPaymentMethodFromProto(pb.PaymentMethod),
		TransactionUUID: pb.TransactionUuid,
	}, nil
}

func mapPaymentMethodFromProto(pm eventsV1.PaymentMethod) model.PaymentMethod {
	switch pm {
	case eventsV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case eventsV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSbp
	case eventsV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case eventsV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}
