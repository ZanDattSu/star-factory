package repository

import (
	"context"

	"github.com/ZanDattSu/star-factory/inventory/internal/model"
)

type PartRepository interface {
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	PutPart(ctx context.Context, uuid string, part *model.Part) error
	ListParts(ctx context.Context) ([]*model.Part, error)
}
