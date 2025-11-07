package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"order/internal/client/converter"
	"order/internal/model"

	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, partsFilter model.PartsFilter) ([]*model.Part, error) {
	parts, err := c.genClient.ListParts(
		ctx,
		&inventoryV1.ListPartsRequest{
			Filter: converter.PartsFilterToProto(partsFilter),
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, NewPartsNotFoundError(partsFilter.Uuids)
		}
		return nil, err
	}

	return converter.PartsToModel(parts.Parts), nil
}
