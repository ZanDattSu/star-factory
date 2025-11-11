package part

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"inventory/internal/converter"

	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.partService.ListParts(ctx, converter.PartsFilterToModel(req.Filter))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
