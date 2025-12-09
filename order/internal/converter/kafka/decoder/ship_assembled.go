package decoder

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/ZanDattSu/star-factory/order/internal/model"
	eventsV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewAssemblyDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsV1.ShipAssembledEvent
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUuid: pb.EventUuid,
		OrderUuid: pb.OrderUuid,
		UserUuid:  pb.UserUuid,
		BuildTime: time.Duration(pb.BuildTimeSec) * time.Second,
	}, nil
}
