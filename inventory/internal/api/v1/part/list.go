package part

import (
	"context"

	"inventory/internal/converter"

	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts := a.partService.ListParts(ctx, converter.PartsFilterToModel(req.Filter))

	return &inventoryV1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
