package inmemory

import (
	"context"

	"inventory/internal/model"
	"inventory/internal/repository/converter"
	repoModel "inventory/internal/repository/model"
)

func (r *repository) ListParts(_ context.Context) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]*repoModel.Part, 0, len(r.parts))
	for _, part := range r.parts {
		parts = append(parts, part)
	}
	return converter.PartsToModel(parts), nil
}
