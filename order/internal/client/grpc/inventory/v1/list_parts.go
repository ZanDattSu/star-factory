package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ZanDattSu/star-factory/order/internal/client/converter"
	"github.com/ZanDattSu/star-factory/order/internal/model"
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
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, fmt.Errorf("inventory: %w", NewPartsNotFoundError(partsFilter.Uuids))
			case codes.Internal:
				return nil, fmt.Errorf("inventory internal error: %w", err)
			case codes.Unavailable:
				return nil, fmt.Errorf("inventory service unavailable: %w", err)
			}
		}

		return nil, fmt.Errorf("inventory ListParts failed: %w", err)
	}

	return converter.PartsToModel(parts.Parts), nil
}
