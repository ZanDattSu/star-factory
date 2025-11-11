package mongodb

import (
	"context"
	"fmt"

	"inventory/internal/model"
	"inventory/internal/repository/converter"
)

func (r *repository) PutPart(ctx context.Context, uuid string, part *model.Part) error {
	if part == nil {
		return fmt.Errorf("part is nil")
	}

	_, err := r.collection.InsertOne(ctx, converter.PartToRepoModel(part))
	if err != nil {
		return fmt.Errorf("failed to insert part %s: %w", uuid, err)
	}

	return nil
}
