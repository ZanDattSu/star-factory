package part

import (
	"context"

	"inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.repository.GetPart(ctx, uuid)
	if err != nil {
		return nil, &model.PartNotFoundError{PartUUID: uuid}
	}
	return part, nil
}
