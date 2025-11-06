package service

import (
	"context"

	"inventory/internal/model"
)

type PartService interface {
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	ListParts(ctx context.Context, req *model.PartsFilter) []*model.Part
}
