package repository

import (
	"context"

	"inventory/internal/model"
)

type PartRepository interface {
	GetPart(cxt context.Context, uuid string) (*model.Part, bool)
	ListParts(ctx context.Context, req *model.PartsFilter) []*model.Part
}
