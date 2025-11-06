package part

import (
	"inventory/internal/model"
)

func (s *service) GetPart(uuid string) (*model.Part, error) {
	part, ok := s.repository.GetPart(uuid)
	if !ok {
		return nil, model.ErrPartNotFound(uuid)
	}
	return part, nil
}
