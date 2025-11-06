package part

import (
	"context"

	"inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part, ok := s.repository.GetPart(ctx, uuid)
	if !ok {
		return nil, &model.PartNotFoundError{PartUUID: uuid}
	}
	return part, nil
}
