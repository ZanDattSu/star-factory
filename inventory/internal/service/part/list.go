package part

import (
	"context"
	"inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, req *model.PartsFilter) []*model.Part {
	parts := s.repository.ListParts(ctx, req)

	return parts
}
