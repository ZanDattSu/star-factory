package service

import (
	"context"
	"inventory/internal/model"
)

type PartService interface {
	GetPart(uuid string) (*model.Part, error)
	ListParts(_ context.Context, req *model.PartsFilter) []*model.Part
}
