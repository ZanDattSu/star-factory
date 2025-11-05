package repository

import (
	"context"

	"inventory/internal/model"
)

type PartRepository interface {
	GetPart(uuid string) (*model.Part, bool)
	ListParts(_ context.Context, req *model.PartsFilter) ([]*model.Part, error)
}
