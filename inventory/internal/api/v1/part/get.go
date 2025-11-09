package part

import (
	"context"
	"errors"

	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"inventory/internal/converter"
	"inventory/internal/model"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.partService.GetPart(ctx, req.Uuid)
	if err != nil {
		var errNotFound *model.PartNotFoundError
		if errors.As(err, &errNotFound) {
			return nil, status.Error(codes.NotFound, errNotFound.Error())
		}
		return nil, err
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
