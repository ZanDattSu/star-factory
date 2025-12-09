package decoder

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/ZanDattSu/star-factory/notification/internal/model"
	eventsV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/events/v1"
)

type shipAssembledDecoder struct{}

func NewAssemblyDecoder() *shipAssembledDecoder {
	return &shipAssembledDecoder{}
}

func (d *shipAssembledDecoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsV1.ShipAssembledEvent
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUUID:    pb.EventUuid,
		OrderUUID:    pb.OrderUuid,
		UserUUID:     pb.UserUuid,
		BuildTimeSec: time.Duration(pb.BuildTimeSec) * time.Second,
	}, nil
}
