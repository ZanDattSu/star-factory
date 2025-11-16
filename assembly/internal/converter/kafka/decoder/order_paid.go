package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/ZanDattSu/star-factory/assembly/internal/model"
	eventsV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderPaidEvent{
		EventUuid:       pb.EventUuid,
		OrderUuid:       pb.OrderUuid,
		UserUuid:        pb.UserUuid,
		PaymentMethod:   model.PaymentMethod(pb.PaymentMethod),
		TransactionUuid: pb.TransactionUuid,
	}, nil
}
