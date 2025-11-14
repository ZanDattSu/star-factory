package service

import (
	"context"

	"github.com/ZanDattSu/star-factory/inventory/internal/model"
)

type PartService interface {
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}
